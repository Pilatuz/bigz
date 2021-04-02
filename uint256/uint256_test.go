// + build ignore

package uint256

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/Pilatuz/bigx/v2/uint128"
)

// rand256 generates single Uint256 random value.
func rand256() Uint256 {
	buf := make([]byte, 32+1) // one extra random byte!
	rand.Read(buf)
	u := LoadLittleEndian(buf)
	if buf[32]&0x03 == 0 {
		u.Lo.Lo = 0 // reset lower half
	}
	if buf[32]&0x0C == 0 {
		u.Lo.Hi = 0 // reset lower half
	}
	if buf[32]&0x30 == 0 {
		u.Hi.Lo = 0 // reset upper half
	}
	if buf[32]&0xC0 == 0 {
		u.Hi.Hi = 0 // reset upper half
	}
	return u
}

// rand256slice generates slice of Uint256 pure random values.
func rand256slice(count int) []Uint256 {
	buf := make([]byte, 32)
	out := make([]Uint256, count)
	for i := range out {
		rand.Read(buf)
		out[i] = LoadLittleEndian(buf)
	}
	return out
}

// generate256s generates a series of pseudo-random Uint256 values
func generate256s(count int, values chan Uint256) {
	defer close(values)

	// a few fixed values
	fixed := []Uint128{uint128.Zero(), uint128.One(), uint128.Max().Sub64(1), uint128.Max()}
	for _, hi := range fixed {
		for _, lo := range fixed {
			values <- Uint256{
				Lo: lo,
				Hi: hi,
			}
		}
	}

	// a few random values
	for i := 0; i < count; i++ {
		values <- rand256()
	}
}

// TestUint128Helpers unit tests for various Uint256 helpers.
func TestUint128Helpers(t *testing.T) {
	t.Run("FromBig", func(t *testing.T) {
		if got := FromBig(nil); !got.Equals(Zero()) {
			t.Fatalf("FromBig(nil) does not equal to 0, got %#x", got)
		}

		if got := FromBig(big.NewInt(-1)); !got.Equals(Zero()) {
			t.Fatalf("FromBig(-1) does not equal to 0, got %#x", got)
		}

		if got := FromBig(new(big.Int).Lsh(big.NewInt(1), 257)); !got.Equals(Max()) {
			t.Fatalf("FromBig(2^257) does not equal to Max(), got %#x", got)
		}
	})

	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint256)
		go generate256s(1000, values)
		for x := range values {
			if got := FromBig(x.Big()); got != x {
				t.Fatalf("FromBig is not the inverse of Big for #%x, got %#x", x, got)
			}

			if !x.Equals(x) {
				t.Fatalf("%#x does not equal itself", x)
			}
			if !From128(x.Lo).Equals128(x.Lo) {
				t.Fatalf("%#v does not equal128 itself", x)
			}
		}
	})
}

// TestUint128Bits unit tests for bit counting helpers.
func TestUint128Bits(t *testing.T) {
	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint256)
		go generate256s(1000, values)
		for x := range values {
			d := newDummy256(x.Big())
			k := int(x.Lo.Lo & 0xFF)

			if expected, got := d.LeadingZeros(), x.LeadingZeros(); got != expected {
				t.Fatalf("mismatch: %#x LeadingZeros should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.TrailingZeros(), x.TrailingZeros(); got != expected {
				t.Fatalf("mismatch: %#x TrailingZeros should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.OnesCount(), x.OnesCount(); got != expected {
				t.Fatalf("mismatch: %#x OnesCount should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.RotateRight(k), newDummy256(x.RotateRight(k).Big()); !expected.Equals(got) {
				t.Fatalf("mismatch: %#x RotateRight should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.RotateLeft(k), newDummy256(x.RotateLeft(k).Big()); !expected.Equals(got) {
				t.Fatalf("mismatch: %#x RotateLeft should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.Reverse(), newDummy256(x.Reverse().Big()); !expected.Equals(got) {
				t.Fatalf("mismatch: %#x RotateRight should equal %v, got %v", x, expected, got)
			}
			if expected, got := x.Big().BitLen(), x.BitLen(); expected != got {
				t.Fatalf("mismatch: %#x BitLen should equal %v, got %v", x, expected, got)
			}
		}
	})
}

// big.Int 2^256 wraparound semantics
var (
	bigOne  = big.NewInt(1)                    // = 1
	bigMod  = new(big.Int).Lsh(bigOne, 256)    // = 2^256
	bigMask = new(big.Int).Sub(bigMod, bigOne) // = 2^256 - 1
)

func mod256(i *big.Int) *big.Int {
	if i.Sign() < 0 {
		i = i.Add(i, bigMod) // just add 2^128 to make it positive
	}
	return i.And(i, bigMask)
}

type (
	BinOp    func(x, y Uint256) Uint256
	BinOp128 func(x Uint256, y Uint128) Uint256
	BinOp64  func(x Uint256, y uint64) Uint256
	BigBinOp func(z, x, y *big.Int) *big.Int

	ShiftOp    func(x Uint256, n uint) Uint256
	BigShiftOp func(z, x *big.Int, n uint) *big.Int
)

// z = op(x, y)
func checkBinOp(t *testing.T, x Uint256, op string, y Uint256, fn BinOp, fnb BigBinOp) {
	t.Helper()
	expected := mod256(fnb(new(big.Int), x.Big(), y.Big()))
	if got := fn(x, y); expected.Cmp(got.Big()) != 0 {
		t.Fatalf("mismatch: (%#x %v %#x) should equal %#x, got %#x", x, op, y, expected, got)
	}
}
func checkBinOp128(t *testing.T, x Uint256, op string, y Uint128, fn BinOp128, fnb BigBinOp) {
	t.Helper()
	expected := mod256(fnb(new(big.Int), x.Big(), From128(y).Big()))
	if got := fn(x, y); expected.Cmp(got.Big()) != 0 {
		t.Fatalf("mismatch: (%#x %v %#x) should equal %#x, got %#x", x, op, y, expected, got)
	}
}

func checkBinOp64(t *testing.T, x Uint256, op string, y uint64, fn BinOp64, fnb BigBinOp) {
	t.Helper()
	expected := mod256(fnb(new(big.Int), x.Big(), From64(y).Big()))
	if got := fn(x, y); expected.Cmp(got.Big()) != 0 {
		t.Fatalf("mismatch: (%#x %v %#x) should equal %#x, got %#x", x, op, y, expected, got)
	}
}

// z = op(x, n)
func checkShiftOp(t *testing.T, x Uint256, op string, n uint, fn ShiftOp, fnb BigShiftOp) {
	t.Helper()
	expected := mod256(fnb(new(big.Int), x.Big(), n))
	if got := fn(x, n); expected.Cmp(got.Big()) != 0 {
		t.Fatalf("mismatch: (%#x %v %v) should equal %#x, got %#x", x, op, n, expected, got)
	}
}

// TestMul unit tests for full 256-bit multiplication.
func TestMul(t *testing.T) {
	xvalues := make(chan Uint256)
	go generate256s(200, xvalues)
	for x := range xvalues {
		yvalues := make(chan Uint256)
		go generate256s(200, yvalues)
		for y := range yvalues {
			hi, lo := Mul(x, y)
			expected := new(big.Int).Mul(x.Big(), y.Big())
			got := new(big.Int).Lsh(hi.Big(), 256)
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

	xvalues := make(chan Uint256)
	go generate256s(10, xvalues)
	for x := range xvalues {
		yvalues := make(chan Uint256)
		go generate256s(10, yvalues)
		for y := range yvalues {
			zvalues := make(chan Uint256)
			go generate256s(10, zvalues)
			for z := range zvalues {
				if z.IsZero() {
					continue
				}
				if z.Cmp(x) <= 0 {
					continue
				}
				q, r := Div(x, y, z)
				xy := new(big.Int).Lsh(x.Big(), 256)
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

// TestArithmetic compare Uint256 arithmetic methods to their math/big equivalents
func TestArithmetic(t *testing.T) {
	xvalues := make(chan Uint256)
	go generate256s(200, xvalues)
	for x := range xvalues {
		yvalues := make(chan Uint256)
		go generate256s(200, yvalues)
		for y := range yvalues {
			// 256 op 256
			checkBinOp(t, x, "+", y, Uint256.Add, (*big.Int).Add)
			checkBinOp(t, x, "-", y, Uint256.Sub, (*big.Int).Sub)
			checkBinOp(t, x, "*", y, Uint256.Mul, (*big.Int).Mul)
			if !y.IsZero() {
				checkBinOp(t, x, "/", y, Uint256.Div, (*big.Int).Div)
				checkBinOp(t, x, "%", y, Uint256.Mod, (*big.Int).Mod)
			}
			checkBinOp(t, x, "&^", y, Uint256.AndNot, (*big.Int).AndNot)
			checkBinOp(t, x, "&", y, Uint256.And, (*big.Int).And)
			checkBinOp(t, x, "|", y, Uint256.Or, (*big.Int).Or)
			checkBinOp(t, x, "^", y, Uint256.Xor, (*big.Int).Xor)
			if expected, got := x.Big().Cmp(y.Big()), x.Cmp(y); expected != got {
				t.Fatalf("mismatch: Cmp(%#x,%#x) should equal %v, got %v", x, y, expected, got)
			}

			// 256 op 128
			y128 := y.Lo
			checkBinOp128(t, x, "+", y128, Uint256.Add128, (*big.Int).Add)
			checkBinOp128(t, x, "-", y128, Uint256.Sub128, (*big.Int).Sub)
			checkBinOp128(t, x, "*", y128, Uint256.Mul128, (*big.Int).Mul)
			if !y128.IsZero() {
				mod128 := func(x Uint256, y Uint128) Uint256 {
					return From128(x.Mod128(y)) // helper to fix signature
				}
				checkBinOp128(t, x, "/", y128, Uint256.Div128, (*big.Int).Div)
				checkBinOp128(t, x, "%", y128, mod128, (*big.Int).Mod)
			}
			if expected, got := x.Big().Cmp(From128(y128).Big()), x.Cmp128(y128); expected != got {
				t.Fatalf("mismatch: Cmp128(%#x,%#x) should equal %v, got %v", x, y128, expected, got)
			}
			checkBinOp128(t, x, "&^", y128, Uint256.AndNot128, (*big.Int).AndNot)
			checkBinOp128(t, x, "&", y128, Uint256.And128, (*big.Int).And)
			checkBinOp128(t, x, "|", y128, Uint256.Or128, (*big.Int).Or)
			checkBinOp128(t, x, "^", y128, Uint256.Xor128, (*big.Int).Xor)

			// 256 op 64
			y64 := y128.Lo
			if y64 != 0 {
				mod64 := func(x Uint256, y uint64) Uint256 {
					return From64(x.Mod64(y)) // helper to fix signature
				}
				checkBinOp64(t, x, "/", y64, Uint256.Div64, (*big.Int).Div)
				checkBinOp64(t, x, "%", y64, mod64, (*big.Int).Mod)
			}

			// shift op
			z := uint(y.Lo.Lo & 0xFF)
			checkShiftOp(t, x, "<<", z, Uint256.Lsh, (*big.Int).Lsh)
			checkShiftOp(t, x, ">>", z, Uint256.Rsh, (*big.Int).Rsh)
		}

		// unary Cmp
		if got := x.Cmp(x); got != 0 {
			t.Fatalf("%#x does not equal itself, got %v", x, got)
		}
		if got := From128(x.Lo).Cmp128(x.Lo); got != 0 {
			t.Fatalf("%#x does not equal itself, got %v", x.Lo, got)
		}

		// unary Not
		if expected, got := mod256(new(big.Int).Not(x.Big())), x.Not(); expected.Cmp(got.Big()) != 0 {
			t.Fatalf("mismatch: (%v %#x) should equal %#x, got %#x", "~", x, expected, got)
		}
	}
}

// dummy raw 256 bits
type dummy256 [256]uint

func newDummy256(b *big.Int) dummy256 {
	n := b.BitLen()
	if n > 256 {
		n = 256 // truncate
	}

	var out dummy256
	for i := 0; i < n; i++ {
		out[i] = b.Bit(i)
	}
	return out
}

func (u dummy256) Equals(v dummy256) bool {
	for i := range u {
		if u[i] != v[i] {
			return false
		}
	}
	return true
}

func (u dummy256) LeadingZeros() int {
	return u.Reverse().TrailingZeros()
}

func (u dummy256) TrailingZeros() int {
	var out int
	for i := range u {
		if u[i] != 0 {
			break
		}
		out++
	}
	return out
}

func (u dummy256) OnesCount() int {
	var out int
	for i := range u {
		if u[i] != 0 {
			out++
		}
	}
	return out
}

func (u dummy256) RotateLeft(k int) dummy256 {
	var out dummy256
	for i := range u {
		out[uint(i+k)%256] = u[i]
	}
	return out
}

func (u dummy256) RotateRight(k int) dummy256 {
	var out dummy256
	for i := range u {
		out[i] = u[uint(i+k)%256]
	}
	return out
}

func (u dummy256) Reverse() dummy256 {
	var out dummy256
	for i := range u {
		out[255-i] = u[i]
	}
	return out
}
