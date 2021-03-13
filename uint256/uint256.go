package uint256

import (
	"math"
	"math/big"
	"math/bits"

	"github.com/Pilatuz/bigx/uint128"
)

// Note, Zero and Max are functions just to make read-only values.
// We cannot define constants for structures, and global variables
// are unacceptable because it will be possible to change them.

// Zero is the lowest possible Uint256 value.
func Zero() Uint256 {
	return From64(0)
}

// One is the lowest non-zero Uint256 value.
func One() Uint256 {
	return From64(1)
}

// Max is the largest possible Uint256 value.
func Max() Uint256 {
	return Uint256{
		Lo: Uint128{
			Lo: math.MaxUint64,
			Hi: math.MaxUint64,
		},
		Hi: Uint128{
			Lo: math.MaxUint64,
			Hi: math.MaxUint64,
		},
	}
}

// Uint128 is an unsigned 128-bit number alias.
type Uint128 = uint128.Uint128

// Uint256 is an unsigned 256-bit number.
// All methods are immutable, works just like standard uint64.
type Uint256 struct {
	Lo Uint128 // low 128-bit half
	Hi Uint128 // high 128-bit half
}

// From64 converts 64-bit value v to a Uint256 value.
// High 64-bit half will be zero.
func From64(v uint64) Uint256 {
	return Uint256{Lo: Uint128{Lo: v}}
}

// FromBig converts *big.Int to 256-bit Uint256 value ignoring overflows.
// If input integer is nil or negative then return Zero.
// If input interger overflows 256-bit then return Max.
func FromBig(i *big.Int) Uint256 {
	u, _ := FromBigX(i)
	return u
}

// FromBigX converts *big.Int to 256-bit Uint256 value (eXtended version).
// Provides ok successful flag as a second return value.
// If input integer is negative or overflows 256-bit then ok=false.
// If input is nil then zero 256-bit returned.
func FromBigX(i *big.Int) (Uint256, bool) {
	switch {
	case i == nil:
		return Zero(), true // assuming nil === 0
	case i.Sign() < 0:
		return Zero(), false // value cannot be negative!
	case i.BitLen() > 256:
		return Max(), false // value overflows 256-bit!
	}

	lo := i
	hi := new(big.Int).Rsh(i, 128)
	return Uint256{
		Lo: uint128.FromBig(lo),
		Hi: uint128.FromBig(hi),
	}, true
}

// Big returns 256-bit value as a *big.Int.
func (u Uint256) Big() *big.Int {
	t := new(big.Int)
	i := new(big.Int).SetUint64(u.Hi.Hi)
	i = i.Lsh(i, 64)
	i = i.Or(i, t.SetUint64(u.Hi.Lo))
	i = i.Lsh(i, 64)
	i = i.Or(i, t.SetUint64(u.Lo.Hi))
	i = i.Lsh(i, 64)
	i = i.Or(i, t.SetUint64(u.Lo.Lo))
	return i
}

// IsZero returns true if stored 256-bit value is zero.
func (u Uint256) IsZero() bool {
	return u.Lo.IsZero() && u.Hi.IsZero()
}

// Equals returns true if two 256-bit values are equal.
// Uint256 values can be compared directly with == operator
// but use of the Equals method is preferred for consistency.
func (u Uint256) Equals(v Uint256) bool {
	return u.Lo.Equals(v.Lo) && u.Hi.Equals(v.Hi)
}

// Equals64 returns true if 256-bit value equals to a 64-bit value.
func (u Uint256) Equals64(v uint64) bool {
	return u.Lo.Equals64(v) && u.Hi.IsZero()
}

// Cmp compares two 256-bit values and returns:
//   -1 if u <  v
//    0 if u == v
//   +1 if u >  v
func (u Uint256) Cmp(v Uint256) int {
	switch {
	case u.Hi.Hi > v.Hi.Hi:
		return +1 // u > v
	case u.Hi.Hi < v.Hi.Hi:
		return -1 // u < v
	case u.Hi.Lo > v.Hi.Lo:
		return +1 // u > v
	case u.Hi.Lo < v.Hi.Lo:
		return -1 // u < v
	case u.Lo.Hi > v.Lo.Hi:
		return +1 // u > v
	case u.Lo.Hi < v.Lo.Hi:
		return -1 // u < v
	case u.Lo.Lo > v.Lo.Lo:
		return +1 // u > v
	case u.Lo.Lo < v.Lo.Lo:
		return -1 // u < v
	}
	return 0 // u == v
}

// Cmp64 compares 256-bit and 64-bit values and returns:
//   -1 if u <  v
//    0 if u == v
//   +1 if u >  v
func (u Uint256) Cmp64(v uint64) int {
	switch {
	case u.Hi.Hi != 0:
		return +1 // u > v
	case u.Hi.Lo != 0:
		return +1 // u > v
	case u.Lo.Hi != 0:
		return +1 // u > v
	case u.Lo.Lo > v:
		return +1 // u > v
	case u.Lo.Lo < v:
		return -1 // u < v
	}
	return 0 // u == v
}

///////////////////////////////////////////////////////////////////////////////
/// logical operators /////////////////////////////////////////////////////////

// Not returns logical NOT (^u) of 256-bit value.
func (u Uint256) Not() Uint256 {
	return Uint256{
		Lo: u.Lo.Not(),
		Hi: u.Hi.Not(),
	}
}

// AndNot returns logical AND NOT (u&^v) of two 256-bit values.
func (u Uint256) AndNot(v Uint256) Uint256 {
	return Uint256{
		Lo: u.Lo.AndNot(v.Lo),
		Hi: u.Hi.AndNot(v.Hi),
	}
}

// AndNot64 returns logical AND NOT (u&v) of 256-bit and 64-bit values.
func (u Uint256) AndNot64(v uint64) Uint256 {
	return Uint256{
		Lo: u.Lo.AndNot64(v),
		Hi: u.Hi, // ^0 == ff..ff
	}
}

// And returns logical AND (u&v) of two 256-bit values.
func (u Uint256) And(v Uint256) Uint256 {
	return Uint256{
		Lo: u.Lo.And(v.Lo),
		Hi: u.Hi.And(v.Hi),
	}
}

// And64 returns logical AND (u&v) of 256-bit and 64-bit values.
func (u Uint256) And64(v uint64) Uint256 {
	return Uint256{
		Lo: u.Lo.And64(v),
		Hi: Uint128{0, 0},
	}
}

// Or returns logical OR (u|v) of two 256-bit values.
func (u Uint256) Or(v Uint256) Uint256 {
	return Uint256{
		Lo: u.Lo.Or(v.Lo),
		Hi: u.Hi.Or(v.Hi),
	}
}

// Or64 returns logical OR (u|v) of 256-bit and 64-bit values.
func (u Uint256) Or64(v uint64) Uint256 {
	return Uint256{
		Lo: u.Lo.Or64(v),
		Hi: u.Hi,
	}
}

// Xor returns logical XOR (u^v) of two 256-bit values.
func (u Uint256) Xor(v Uint256) Uint256 {
	return Uint256{
		Lo: u.Lo.Xor(v.Lo),
		Hi: u.Hi.Xor(v.Hi),
	}
}

// Xor64 returns logical XOR (u^v) of 256-bit and 64-bit values.
func (u Uint256) Xor64(v uint64) Uint256 {
	return Uint256{
		Lo: u.Lo.Xor64(v),
		Hi: u.Hi,
	}
}

///////////////////////////////////////////////////////////////////////////////
/// arithmetic operators //////////////////////////////////////////////////////

// Add returns sum (u+v) of two 256-bit values.
// Wrap-around semantic is used here: Max().Add(From64(1)) == Zero()
func (u Uint256) Add(v Uint256) Uint256 {
	lolo, c0 := bits.Add64(u.Lo.Lo, v.Lo.Lo, 0)
	lohi, c1 := bits.Add64(u.Lo.Hi, v.Lo.Hi, c0)
	hilo, c2 := bits.Add64(u.Hi.Lo, v.Hi.Lo, c1)
	hihi, _ := bits.Add64(u.Hi.Hi, v.Hi.Hi, c2)
	return Uint256{
		Lo: Uint128{Lo: lolo, Hi: lohi},
		Hi: Uint128{Lo: hilo, Hi: hihi},
	}
}

// Add64 returns sum u+v of 256-bit and 64-bit values.
// Wrap-around semantic is used here: Max().Add64(1) == Zero()
func (u Uint256) Add64(v uint64) Uint256 {
	lolo, c0 := bits.Add64(u.Lo.Lo, v, 0)
	lohi, c1 := bits.Add64(u.Lo.Hi, 0, c0)
	hilo, c2 := bits.Add64(u.Hi.Lo, 0, c1)
	hihi, _ := bits.Add64(u.Hi.Hi, 0, c2)
	return Uint256{
		Lo: Uint128{Lo: lolo, Hi: lohi},
		Hi: Uint128{Lo: hilo, Hi: hihi},
	}
}

// Sub returns difference (u-v) of two 256-bit values.
// Wrap-around semantic is used here: Zero().Sub(From64(1)) == Max().
func (u Uint256) Sub(v Uint256) Uint256 {
	lolo, b0 := bits.Sub64(u.Lo.Lo, v.Lo.Lo, 0)
	lohi, b1 := bits.Sub64(u.Lo.Hi, v.Lo.Hi, b0)
	hilo, b2 := bits.Sub64(u.Hi.Lo, v.Hi.Lo, b1)
	hihi, _ := bits.Sub64(u.Hi.Hi, v.Hi.Hi, b2)
	return Uint256{
		Lo: Uint128{Lo: lolo, Hi: lohi},
		Hi: Uint128{Lo: hilo, Hi: hihi},
	}
}

// Sub64 returns difference (u-v) of 256-bit and 64-bit values.
// Wrap-around semantic is used here: Zero().Sub64(1) == Max().
func (u Uint256) Sub64(v uint64) Uint256 {
	lolo, b0 := bits.Sub64(u.Lo.Lo, v, 0)
	lohi, b1 := bits.Sub64(u.Lo.Hi, 0, b0)
	hilo, b2 := bits.Sub64(u.Hi.Lo, 0, b1)
	hihi, _ := bits.Sub64(u.Hi.Hi, 0, b2)
	return Uint256{
		Lo: Uint128{Lo: lolo, Hi: lohi},
		Hi: Uint128{Lo: hilo, Hi: hihi},
	}
}

// Mul returns multiplication (u*v) of two 256-bit values.
// Wrap-around semantic is used here: Max().Mul(Max()) == From64(1).
func (u Uint256) Mul(v Uint256) Uint256 {
	/*hi, lo := bits.Mul64(u.Lo, v.Lo)
	hi += u.Hi*v.Lo + u.Lo*v.Hi
	return Uint256{Lo: lo, Hi: hi}*/
	return Zero()
}

// Mul64 returns multiplication (u*v) of 256-bit and 64-bit values.
// Wrap-around semantic is used here: Max().Mul64(2) == Max().Sub64(1).
func (u Uint256) Mul64(v uint64) Uint256 {
	/*hi, lo := bits.Mul64(u.Lo, v)
	return Uint256{
		Lo: lo,
		Hi: hi + u.Hi*v,
	}*/
	return Zero()
}

// Div returns division (u/v) of two 256-bit values.
func (u Uint256) Div(v Uint256) Uint256 {
	q, _ := u.QuoRem(v)
	return q
}

// Div64 returns division (u/v) of 256-bit and 64-bit values.
func (u Uint256) Div64(v uint64) Uint256 {
	q, _ := u.QuoRem64(v)
	return q
}

// Mod returns modulo (u%v) of two 256-bit values.
func (u Uint256) Mod(v Uint256) Uint256 {
	_, r := u.QuoRem(v)
	return r
}

// Mod64 returns modulo (u%v) of 256-bit and 64-bit values.
func (u Uint256) Mod64(v uint64) uint64 {
	_, r := u.QuoRem64(v)
	return r
}

// QuoRem returns quotient (u/v) and remainder (u%v) of two 256-bit values.
func (u Uint256) QuoRem(v Uint256) (Uint256, Uint256) {
	/*if v.Hi.IsZero() {
		q, r := u.QuoRem128(v.Lo)
		return q, From128(r)
	}

	// generate a "trial quotient," guaranteed to be
	// within 1 of the actual quotient, then adjust.
	n := uint(bits.LeadingZeros64(v.Hi))
	u1, v1 := u.Rsh(1), v.Lsh(n)
	tq, _ := bits.Div64(u1.Hi, u1.Lo, v1.Hi)
	tq >>= 63 - n
	if tq != 0 {
		tq--
	}

	// calculate remainder using trial quotient, then
	// adjust if remainder is greater than divisor
	q, r := From64(tq), u.Sub(v.Mul64(tq))
	if r.Cmp(v) >= 0 {
		q = q.Add64(1)
		r = r.Sub(v)
	}

	return q, r*/
	return Zero(), Zero()
}

// QuoRem64 returns quotient (u/v) and remainder (u%v) of 256-bit and 64-bit values.
func (u Uint256) QuoRem64(v uint64) (Uint256, uint64) {
	/*if u.Hi < v {
		lo, r := bits.Div64(u.Hi, u.Lo, v)
		return Uint256{Lo: lo}, r
	}

	hi, r := bits.Div64(0, u.Hi, v)
	lo, r := bits.Div64(r, u.Lo, v)
	return Uint256{Lo: lo, Hi: hi}, r*/
	return Zero(), 0
}

///////////////////////////////////////////////////////////////////////////////
/// shift operators ///////////////////////////////////////////////////////////

// Lsh returns left shift (u<<n).
func (u Uint256) Lsh(n uint) Uint256 {
	if n > 128 {
		return Uint256{
			// Lo: 0,
			Hi: u.Lo.Lsh(n - 128),
		}
	}

	return Zero() /*Uint256{
		Lo: u.Lo << n,
		Hi: u.Hi<<n | u.Lo>>(64-n),
	}*/
}

// Rsh returns right shift (u>>n).
func (u Uint256) Rsh(n uint) Uint256 {
	if n > 128 {
		return Uint256{
			Lo: u.Hi.Rsh(n - 128),
			// Hi: 0,
		}
	}

	return Zero() /*Uint256{
		Lo: u.Lo>>n | u.Hi<<(64-n),
		Hi: u.Hi >> n,
	}*/
}

// RotateLeft returns the value of u rotated left by (k mod 256) bits.
func (u Uint256) RotateLeft(k int) Uint256 {
	n := uint(k) & 255

	if n < 64 {
		if n == 0 {
			return u
		}

		return Uint256{
			Lo: Uint128{
				Lo: u.Lo.Lo<<n | u.Hi.Hi>>(64-n),
				Hi: u.Lo.Hi<<n | u.Lo.Lo>>(64-n),
			},
			Hi: Uint128{
				Lo: u.Lo.Lo<<n | u.Hi.Hi>>(64-n),
				Hi: u.Hi.Lo<<n | u.Lo.Hi>>(64-n),
			},
		}
	}

	n -= 64
	if n == 0 {
		return Uint256{
			Lo: u.Hi,
			Hi: u.Lo,
		}
	}

	return Zero() /*Uint256{
		Lo: u.Lo>>(64-n) | u.Hi<<n,
		Hi: u.Hi>>(64-n) | u.Lo<<n,
	}*/
}

// RotateRight returns the value of u rotated left by (k mod 256) bits.
func (u Uint256) RotateRight(k int) Uint256 {
	return u.RotateLeft(-k)
}

///////////////////////////////////////////////////////////////////////////////
/// bit counting //////////////////////////////////////////////////////////////

// BitLen returns the minimum number of bits required to represent 256-bit value.
// The result is 0 for u == 0.
func (u Uint256) BitLen() int {
	if !u.Hi.IsZero() {
		return 128 + u.Hi.BitLen()
	}
	return u.Lo.BitLen()
}

// LeadingZeros returns the number of leading zero bits.
// The result is 256 for u == 0.
func (u Uint256) LeadingZeros() int {
	if !u.Hi.IsZero() {
		return u.Hi.LeadingZeros()
	}
	return 128 + u.Lo.LeadingZeros()
}

// TrailingZeros returns the number of trailing zero bits.
// The result is 256 for u == 0.
func (u Uint256) TrailingZeros() int {
	if !u.Lo.IsZero() {
		return u.Lo.TrailingZeros()
	}
	return 128 + u.Hi.TrailingZeros()
}

// OnesCount returns the number of one bits ("population count").
func (u Uint256) OnesCount() int {
	return u.Lo.OnesCount() +
		u.Hi.OnesCount()
}

// Reverse returns the value with bits in reversed order.
func (u Uint256) Reverse() Uint256 {
	return Uint256{
		Lo: u.Hi.Reverse(),
		Hi: u.Lo.Reverse(),
	}
}

// ReverseBytes returns the value with bytes in reversed order.
func (u Uint256) ReverseBytes() Uint256 {
	return Uint256{
		Lo: u.Hi.ReverseBytes(),
		Hi: u.Lo.ReverseBytes(),
	}
}
