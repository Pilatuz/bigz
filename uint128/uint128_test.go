package uint128

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"testing"
)

// rand128 generates single Uint128 random value.
func rand128() Uint128 {
	buf := make([]byte, 16+1) // one extra random byte!
	rand.Read(buf)
	u := LoadLittleEndian(buf)
	if buf[16]&0x07 == 0 {
		u.Lo = 0 // reset lower half
	}
	if buf[16]&0x70 == 0 {
		u.Hi = 0 // reset upper half
	}
	return u
}

// rand128slice generates slice of Uint128 pure random values.
func rand128slice(count int) []Uint128 {
	buf := make([]byte, 16)
	out := make([]Uint128, count)
	for i := range out {
		rand.Read(buf)
		out[i] = LoadLittleEndian(buf)
	}
	return out
}

// generate128s generates a series of pseudo-random Uint128 values
func generate128s(count int, values chan Uint128) {
	defer close(values)

	// a few fixed values
	fixed := []uint64{0, 1, 2, math.MaxUint64 - 1, math.MaxUint64}
	for _, hi := range fixed {
		for _, lo := range fixed {
			values <- Uint128{Lo: lo, Hi: hi}
		}
	}

	// a few random values
	for i := 0; i < count; i++ {
		values <- rand128()
	}
}

// TestUint128Helpers unit tests for various Uint128 helpers.
func TestUint128Helpers(t *testing.T) {
	t.Run("FromBig", func(t *testing.T) {
		if got := FromBig(nil); !got.Equals(Zero()) {
			t.Fatalf("FromBig(nil) does not equal to 0, got %#x", got)
		}

		if got := FromBig(big.NewInt(-1)); !got.Equals(Zero()) {
			t.Fatalf("FromBig(-1) does not equal to 0, got %#x", got)
		}

		if got := FromBig(new(big.Int).Lsh(big.NewInt(1), 129)); !got.Equals(Max()) {
			t.Fatalf("FromBig(2^129) does not equal to Max(), got %#x", got)
		}
	})

	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint128)
		go generate128s(1000, values)
		for x := range values {
			if got := FromBig(x.Big()); got != x {
				t.Fatalf("FromBig is not the inverse of Big for #%x, got %#x", x, got)
			}

			if !x.Equals(x) {
				t.Fatalf("%#x does not equal itself", x)
			}
			if !From64(x.Lo).Equals64(x.Lo) {
				t.Fatalf("%#v does not equal64 itself", x)
			}
		}
	})
}

// TestUint128Bits unit tests for bit counting helpers.
func TestUint128Bits(t *testing.T) {
	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint128)
		go generate128s(1000, values)
		for x := range values {
			d := newDummy128(x.Big())
			k := int(x.Lo & 0xFF)

			if expected, got := d.LeadingZeros(), x.LeadingZeros(); got != expected {
				t.Fatalf("mismatch: %#x LeadingZeros should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.TrailingZeros(), x.TrailingZeros(); got != expected {
				t.Fatalf("mismatch: %#x TrailingZeros should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.OnesCount(), x.OnesCount(); got != expected {
				t.Fatalf("mismatch: %#x OnesCount should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.RotateRight(k), newDummy128(x.RotateRight(k).Big()); !expected.Equals(got) {
				t.Fatalf("mismatch: %#x RotateRight should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.RotateLeft(k), newDummy128(x.RotateLeft(k).Big()); !expected.Equals(got) {
				t.Fatalf("mismatch: %#x RotateLeft should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.Reverse(), newDummy128(x.Reverse().Big()); !expected.Equals(got) {
				t.Fatalf("mismatch: %#x RotateRight should equal %v, got %v", x, expected, got)
			}
			if expected, got := x.Big().BitLen(), x.BitLen(); expected != got {
				t.Fatalf("mismatch: %#x BitLen should equal %v, got %v", x, expected, got)
			}
		}
	})
}

// big.Int 2^128 wraparound semantics
var (
	bigOne  = big.NewInt(1)                    // = 1
	bigMod  = new(big.Int).Lsh(bigOne, 128)    // = 2^128
	bigMask = new(big.Int).Sub(bigMod, bigOne) // = 2^128 - 1
)

func mod128(i *big.Int) *big.Int {
	if i.Sign() < 0 {
		i = i.Add(i, bigMod) // just add 2^128 to make it positive
	}
	return i.And(i, bigMask)
}

type (
	BinOp    func(x, y Uint128) Uint128
	BinOp64  func(x Uint128, y uint64) Uint128
	BigBinOp func(z, x, y *big.Int) *big.Int

	ShiftOp    func(x Uint128, n uint) Uint128
	BigShiftOp func(z, x *big.Int, n uint) *big.Int
)

// z = op(x, y)
func checkBinOp(t *testing.T, x Uint128, op string, y Uint128, fn BinOp, fnb BigBinOp) {
	t.Helper()
	expected := mod128(fnb(new(big.Int), x.Big(), y.Big()))
	if got := fn(x, y); expected.Cmp(got.Big()) != 0 {
		t.Fatalf("mismatch: (%#x %v %#x) should equal %#x, got %#x", x, op, y, expected, got)
	}
}
func checkBinOp64(t *testing.T, x Uint128, op string, y uint64, fn BinOp64, fnb BigBinOp) {
	t.Helper()
	expected := mod128(fnb(new(big.Int), x.Big(), From64(y).Big()))
	if got := fn(x, y); expected.Cmp(got.Big()) != 0 {
		t.Fatalf("mismatch: (%#x %v %#x) should equal %#x, got %#x", x, op, y, expected, got)
	}
}

// z = op(x, n)
func checkShiftOp(t *testing.T, x Uint128, op string, n uint, fn ShiftOp, fnb BigShiftOp) {
	t.Helper()
	expected := mod128(fnb(new(big.Int), x.Big(), n))
	if got := fn(x, n); expected.Cmp(got.Big()) != 0 {
		t.Fatalf("mismatch: (%#x %v %v) should equal %#x, got %#x", x, op, n, expected, got)
	}
}

// TestMul unit tests for full 128-bit multiplication.
func TestMul(t *testing.T) {
	xvalues := make(chan Uint128)
	go generate128s(200, xvalues)
	for x := range xvalues {
		yvalues := make(chan Uint128)
		go generate128s(200, yvalues)
		for y := range yvalues {
			hi, lo := Mul(x, y)
			expected := new(big.Int).Mul(x.Big(), y.Big())
			got := new(big.Int).Lsh(hi.Big(), 128)
			got.Or(got, lo.Big())
			if expected.Cmp(got) != 0 {
				t.Fatalf("%x * %x != %x, got %x", x, y, expected, got)
			}
		}
	}
}

// TestDiv unit tests for full 256-bit division.
func TestDiv(t *testing.T) {
	t.Run("div_by_zero", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				expected := "integer divide by zero"
				if fmt.Sprintf("%v", r) != expected {
					t.Fatalf("unexpected panic: %v", r)
				}
			} else {
				t.Fatalf("expected panic, got nothing")
			}
		}()
		Div(One(), One(), Zero())
	})

	t.Run("overflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				expected := "integer overflow"
				if fmt.Sprintf("%v", r) != expected {
					t.Fatalf("unexpected panic: %v", r)
				}
			} else {
				t.Fatalf("expected panic, got nothing")
			}
		}()
		Div(Max(), One(), One())
	})

	xvalues := make(chan Uint128)
	go generate128s(100, xvalues)
	for x := range xvalues {
		yvalues := make(chan Uint128)
		go generate128s(100, yvalues)
		for y := range yvalues {
			zvalues := make(chan Uint128)
			go generate128s(100, zvalues)
			for z := range zvalues {
				if z.IsZero() {
					continue
				}
				if z.Cmp(x) <= 0 {
					continue
				}
				q, r := Div(x, y, z)
				xy := new(big.Int).Lsh(x.Big(), 128)
				xy.Or(xy, y.Big())
				expectedq, expectedr := new(big.Int).QuoRem(xy, z.Big(), new(big.Int))
				if expectedq.Cmp(q.Big()) != 0 {
					t.Fatalf("%x / %x != %x, got %x", xy, z, expectedq, q)
				}
				if expectedr.Cmp(r.Big()) != 0 {
					t.Fatalf("%x %% %x != %x, got %x", xy, z, expectedr, r)
				}
			}
		}
	}
}

// TestArithmetic compare Uint128 arithmetic methods to their math/big equivalents
func TestArithmetic(t *testing.T) {
	xvalues := make(chan Uint128)
	go generate128s(200, xvalues)
	for x := range xvalues {
		yvalues := make(chan Uint128)
		go generate128s(200, yvalues)
		for y := range yvalues {
			// 128 op 128
			checkBinOp(t, x, "+", y, Uint128.Add, (*big.Int).Add)
			checkBinOp(t, x, "-", y, Uint128.Sub, (*big.Int).Sub)
			checkBinOp(t, x, "*", y, Uint128.Mul, (*big.Int).Mul)
			if !y.IsZero() {
				checkBinOp(t, x, "/", y, Uint128.Div, (*big.Int).Div)
				checkBinOp(t, x, "%", y, Uint128.Mod, (*big.Int).Mod)
			}
			checkBinOp(t, x, "&^", y, Uint128.AndNot, (*big.Int).AndNot)
			checkBinOp(t, x, "&", y, Uint128.And, (*big.Int).And)
			checkBinOp(t, x, "|", y, Uint128.Or, (*big.Int).Or)
			checkBinOp(t, x, "^", y, Uint128.Xor, (*big.Int).Xor)
			if expected, got := x.Big().Cmp(y.Big()), x.Cmp(y); expected != got {
				t.Fatalf("mismatch: Cmp(%#x,%#x) should equal %v, got %v", x, y, expected, got)
			}

			// 128 op 64
			y64 := y.Lo
			checkBinOp64(t, x, "+", y64, Uint128.Add64, (*big.Int).Add)
			checkBinOp64(t, x, "-", y64, Uint128.Sub64, (*big.Int).Sub)
			checkBinOp64(t, x, "*", y64, Uint128.Mul64, (*big.Int).Mul)
			if y64 != 0 {
				mod64 := func(x Uint128, y uint64) Uint128 {
					return From64(x.Mod64(y)) // helper to fix signature
				}
				checkBinOp64(t, x, "/", y64, Uint128.Div64, (*big.Int).Div)
				checkBinOp64(t, x, "%", y64, mod64, (*big.Int).Mod)
			}
			checkBinOp64(t, x, "&^", y64, Uint128.AndNot64, (*big.Int).AndNot)
			checkBinOp64(t, x, "&", y64, Uint128.And64, (*big.Int).And)
			checkBinOp64(t, x, "|", y64, Uint128.Or64, (*big.Int).Or)
			checkBinOp64(t, x, "^", y64, Uint128.Xor64, (*big.Int).Xor)
			if expected, got := x.Big().Cmp(From64(y64).Big()), x.Cmp64(y64); expected != got {
				t.Fatalf("mismatch: Cmp64(%#x,%#x) should equal %v, got %v", x, y64, expected, got)
			}

			// shift op
			z := uint(y.Lo & 0xFF)
			checkShiftOp(t, x, "<<", z, Uint128.Lsh, (*big.Int).Lsh)
			checkShiftOp(t, x, ">>", z, Uint128.Rsh, (*big.Int).Rsh)
		}

		// unary Cmp
		if got := x.Cmp(x); got != 0 {
			t.Fatalf("%#x does not equal itself, got %v", x, got)
		}
		if got := From64(x.Lo).Cmp64(x.Lo); got != 0 {
			t.Fatalf("%#x does not equal itself, got %v", x.Lo, got)
		}

		// unary Not
		if expected, got := mod128(new(big.Int).Not(x.Big())), x.Not(); expected.Cmp(got.Big()) != 0 {
			t.Fatalf("mismatch: (%v %#x) should equal %#x, got %#x", "~", x, expected, got)
		}
	}
}

// dummy raw 128 bits
type dummy128 [128]uint

func newDummy128(b *big.Int) dummy128 {
	n := b.BitLen()
	if n > 128 {
		n = 128 // truncate
	}

	var out dummy128
	for i := 0; i < n; i++ {
		out[i] = b.Bit(i)
	}
	return out
}

func (u dummy128) Equals(v dummy128) bool {
	for i := range u {
		if u[i] != v[i] {
			return false
		}
	}
	return true
}

func (u dummy128) LeadingZeros() int {
	return u.Reverse().TrailingZeros()
}

func (u dummy128) TrailingZeros() int {
	var out int
	for i := range u {
		if u[i] != 0 {
			break
		}
		out++
	}
	return out
}

func (u dummy128) OnesCount() int {
	var out int
	for i := range u {
		if u[i] != 0 {
			out++
		}
	}
	return out
}

func (u dummy128) RotateLeft(k int) dummy128 {
	var out dummy128
	for i := range u {
		out[uint(i+k)%128] = u[i]
	}
	return out
}

func (u dummy128) RotateRight(k int) dummy128 {
	var out dummy128
	for i := range u {
		out[i] = u[uint(i+k)%128]
	}
	return out
}

func (u dummy128) Reverse() dummy128 {
	var out dummy128
	for i := range u {
		out[127-i] = u[i]
	}
	return out
}
