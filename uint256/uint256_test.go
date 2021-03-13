// + build ignore

package uint256

import (
	"crypto/rand"
	"math"
	"math/big"
	"testing"
)

// rand256 generates single Uint256 random value.
func rand256() Uint256 {
	buf := make([]byte, 32+1) // one extra random byte!
	rand.Read(buf)
	u := LoadLittleEndian(buf)
	if buf[32]&0x03 == 0 {
		u.Lo.Lo = 0 // reset low half
	}
	if buf[32]&0x0C == 0 {
		u.Lo.Hi = 0 // reset low half
	}
	if buf[32]&0x30 == 0 {
		u.Hi.Lo = 0 // reset high half
	}
	if buf[32]&0xC0 == 0 {
		u.Hi.Hi = 0 // reset high half
	}
	return u
}

// generate256s generates a series of pseudo-random Uint256 values
func generate256s(count int, values chan Uint256) {
	defer close(values)

	// a few fixed values
	fixed := []uint64{0, 1, 2, math.MaxUint64 - 1, math.MaxUint64}
	for _, hihi := range fixed {
		for _, hilo := range fixed {
			for _, lohi := range fixed {
				for _, lolo := range fixed {
					values <- Uint256{
						Lo: Uint128{Lo: lolo, Hi: lohi},
						Hi: Uint128{Lo: hilo, Hi: hihi},
					}
				}
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
			if !From64(x.Lo.Lo).Equals64(x.Lo.Lo) {
				t.Fatalf("%#v does not equal64 itself", x)
			}
		}
	})
}

// TestUint128Bits unit tests for bit counting helpers.
/*func TestUint128Bits(t *testing.T) {
	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint256)
		go generate128s(1000, values)
		for x := range values {
			d := newDummy128(x.Big())
			k := int(x.Lo & 0xFF)

			if expected, got := d.LeadingZeros(), x.LeadingZeros(); got != expected {
				t.Errorf("mismatch: %#x LeadingZeros should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.TrailingZeros(), x.TrailingZeros(); got != expected {
				t.Errorf("mismatch: %#x TrailingZeros should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.OnesCount(), x.OnesCount(); got != expected {
				t.Errorf("mismatch: %#x OnesCount should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.RotateRight(k), newDummy128(x.RotateRight(k).Big()); !expected.Equals(got) {
				t.Errorf("mismatch: %#x RotateRight should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.RotateLeft(k), newDummy128(x.RotateLeft(k).Big()); !expected.Equals(got) {
				t.Errorf("mismatch: %#x RotateLeft should equal %v, got %v", x, expected, got)
			}
			if expected, got := d.Reverse(), newDummy128(x.Reverse().Big()); !expected.Equals(got) {
				t.Errorf("mismatch: %#x RotateRight should equal %v, got %v", x, expected, got)
			}
			if expected, got := x.Big().BitLen(), x.BitLen(); expected != got {
				t.Errorf("mismatch: %#x BitLen should equal %v, got %v", x, expected, got)
			}
		}
	})
}*/

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

// TestArithmetic compare Uint256 arithmetic methods to their math/big equivalents
func TestArithmetic(t *testing.T) {
	xvalues := make(chan Uint256)
	go generate256s(200, xvalues)
	for x := range xvalues {
		yvalues := make(chan Uint256)
		go generate256s(200, yvalues)
		for y := range yvalues {
			// 128 op 128
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
				t.Fatalf("mismatch: cmp(%#x,%#x) should equal %v, got %v", x, y, expected, got)
			}

			// 128 op 64
			y64 := y.Lo.Lo
			checkBinOp64(t, x, "+", y64, Uint256.Add64, (*big.Int).Add)
			checkBinOp64(t, x, "-", y64, Uint256.Sub64, (*big.Int).Sub)
			checkBinOp64(t, x, "*", y64, Uint256.Mul64, (*big.Int).Mul)
			if y64 != 0 {
				mod64 := func(x Uint256, y uint64) Uint256 {
					return From64(x.Mod64(y)) // helper to fix signature
				}
				checkBinOp64(t, x, "/", y64, Uint256.Div64, (*big.Int).Div)
				checkBinOp64(t, x, "%", y64, mod64, (*big.Int).Mod)
			}
			checkBinOp64(t, x, "&^", y64, Uint256.AndNot64, (*big.Int).AndNot)
			checkBinOp64(t, x, "&", y64, Uint256.And64, (*big.Int).And)
			checkBinOp64(t, x, "|", y64, Uint256.Or64, (*big.Int).Or)
			checkBinOp64(t, x, "^", y64, Uint256.Xor64, (*big.Int).Xor)
			if expected, got := x.Big().Cmp(From64(y64).Big()), x.Cmp64(y64); expected != got {
				t.Fatalf("mismatch: cmp64(%#x,%#x) should equal %v, got %v", x, y64, expected, got)
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
		if got := From64(x.Lo.Lo).Cmp64(x.Lo.Lo); got != 0 {
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
		out[256-i] = u[i]
	}
	return out
}
