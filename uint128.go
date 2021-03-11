package uint128

import (
	"math"
	"math/big"
	"math/bits"
)

// Note, Zero and Max are functions just to make read-only values.
// We cannot define constants for structures, and global variables
// are unacceptable because it will be possible to change them.

// Zero is the lowest possible Uint128 value.
func Zero() Uint128 {
	return Uint128{0, 0}
}

// Max is the largest possible Uint128 value.
func Max() Uint128 {
	return Uint128{
		Lo: math.MaxUint64,
		Hi: math.MaxUint64,
	}
}

// Uint128 is an unsigned 128-bit number.
// All methods are immutable, works just like standard uint64.
type Uint128 struct {
	Lo uint64 // lower 64-bit half
	Hi uint64 // higher 64-bit half
}

// Note, there in no New(lo, hi) just not to confuse
// which half goes first: lower or hight.
// Use structure initialization Uint128{Lo: ..., Hi: ...} instead.

// From64 converts 64-bit value v to a Uint128 value.
// Higher 64-bit half will be zero.
func From64(v uint64) Uint128 {
	return Uint128{Lo: v}
}

// FromBig converts *big.Int to 128-bit Uint128 value ignoring overflows.
// If input integer is nil or negative then return Zero.
// If input interger overflows 128-bit then return Max.
func FromBig(i *big.Int) Uint128 {
	u, _ := FromBigX(i)
	return u
}

// FromBigX converts *big.Int to 128-bit Uint128 value (eXtended version).
// Provides ok successful flag as a second return value.
// If input integer is negative or overflows 128-bit then ok=false.
// If input if nil then zero 128-bit returned.
func FromBigX(i *big.Int) (Uint128, bool) {
	switch {
	case i == nil:
		return Zero(), true // assuming nil === 0
	case i.Sign() < 0:
		return Zero(), false // value cannot be negative!
	case i.BitLen() > 128:
		return Max(), false // value overflows 128-bit!
	}

	// Note, actually result of big.Int.Uint64 is undefined
	// if stored value is greater than 2^64
	// but we assume that it just gets lower 64 bits.
	lo := i.Uint64()
	hi := new(big.Int).Rsh(i, 64).Uint64()
	return Uint128{
		Lo: lo,
		Hi: hi,
	}, true
}

// Big returns 128-bit value as a *big.Int.
func (u Uint128) Big() *big.Int {
	i := new(big.Int).SetUint64(u.Hi)
	i = i.Lsh(i, 64)
	i = i.Or(i, new(big.Int).SetUint64(u.Lo))
	return i
}

// IsZero returns true if stored 128-bit value is zero.
func (u Uint128) IsZero() bool {
	return (u.Lo == 0) && (u.Hi == 0)
}

// Equals returns true if two 128-bit values are equal.
// Uint128 values can be compared directly with == operator
// but use of the Equals method is preferred for consistency.
func (u Uint128) Equals(v Uint128) bool {
	return (u.Lo == v.Lo) && (u.Hi == v.Hi)
}

// Equals64 returns true if 128-bit value equals to a 64-bit value.
func (u Uint128) Equals64(v uint64) bool {
	return (u.Lo == v) && (u.Hi == 0)
}

// Cmp compares two 128-bit values and returns:
//   -1 if u <  v
//    0 if u == v
//   +1 if u >  v
func (u Uint128) Cmp(v Uint128) int {
	switch {
	case u.Hi > v.Hi:
		return +1 // u > v
	case u.Hi < v.Hi:
		return -1 // u < v
	case u.Lo > v.Lo:
		return +1 // u > v
	case u.Lo < v.Lo:
		return -1 // u < v
	}
	return 0 // u == v
}

// Cmp64 compares 128-bit and 64-bit values and returns:
//   -1 if u <  v
//    0 if u == v
//   +1 if u >  v
func (u Uint128) Cmp64(v uint64) int {
	switch {
	case u.Hi != 0:
		return +1 // u > v
	case u.Lo > v:
		return +1 // u > v
	case u.Lo < v:
		return -1 // u < v
	}
	return 0 // u == v
}

///////////////////////////////////////////////////////////////////////////////
/// logical operators /////////////////////////////////////////////////////////

// Not returns logical NOT (^u) of 128-bit value.
func (u Uint128) Not() Uint128 {
	return Uint128{
		Lo: ^u.Lo,
		Hi: ^u.Hi,
	}
}

// AndNot returns logical AND NOT (u&^v) of two 128-bit values.
func (u Uint128) AndNot(v Uint128) Uint128 {
	return Uint128{
		Lo: u.Lo & ^v.Lo,
		Hi: u.Hi & ^v.Hi,
	}
}

// AndNot64 returns logical AND NOT (u&v) of 128-bit and 64-bit values.
func (u Uint128) AndNot64(v uint64) Uint128 {
	return Uint128{
		Lo: u.Lo & ^v,
		Hi: u.Hi, // ^0 == ff..ff
	}
}

// And returns logical AND (u&v) of two 128-bit values.
func (u Uint128) And(v Uint128) Uint128 {
	return Uint128{
		Lo: u.Lo & v.Lo,
		Hi: u.Hi & v.Hi,
	}
}

// And64 returns logical AND (u&v) of 128-bit and 64-bit values.
func (u Uint128) And64(v uint64) Uint128 {
	return Uint128{
		Lo: u.Lo & v,
		Hi: u.Hi & 0, // hope compiler is able to always put just zero instead
	}
}

// Or returns logical OR (u|v) of two 128-bit values.
func (u Uint128) Or(v Uint128) Uint128 {
	return Uint128{
		Lo: u.Lo | v.Lo,
		Hi: u.Hi | v.Hi,
	}
}

// Or64 returns logical OR (u|v) of 128-bit and 64-bit values.
func (u Uint128) Or64(v uint64) Uint128 {
	return Uint128{
		Lo: u.Lo | v,
		Hi: u.Hi | 0, // hope compiler is able to always put just u.Hi instead
	}
}

// Xor returns logical XOR (u^v) of two 128-bit values.
func (u Uint128) Xor(v Uint128) Uint128 {
	return Uint128{
		Lo: u.Lo ^ v.Lo,
		Hi: u.Hi ^ v.Hi,
	}
}

// Xor64 returns logical XOR (u^v) of 128-bit and 64-bit values.
func (u Uint128) Xor64(v uint64) Uint128 {
	return Uint128{
		Lo: u.Lo ^ v,
		Hi: u.Hi ^ 0, // hope compiler is able to always put just u.Hi instead
	}
}

///////////////////////////////////////////////////////////////////////////////
/// arithmetic operators //////////////////////////////////////////////////////

// Add returns sum (u+v) of two 128-bit values.
// Wrap-around semantic is used here: Max().Add(From64(1)) == Zero()
func (u Uint128) Add(v Uint128) Uint128 {
	lo, c0 := bits.Add64(u.Lo, v.Lo, 0)
	hi, _ := bits.Add64(u.Hi, v.Hi, c0)
	return Uint128{Lo: lo, Hi: hi}
}

// Add64 returns sum u+v of 128-bit and 64-bit values.
// Wrap-around semantic is used here: Max().Add64(1) == Zero()
func (u Uint128) Add64(v uint64) Uint128 {
	lo, c0 := bits.Add64(u.Lo, v, 0)
	return Uint128{Lo: lo, Hi: u.Hi + c0}
}

// Sub returns difference (u-v) of two 128-bit values.
// Wrap-around semantic is used here: Zero().Sub(From64(1)) == Max().
func (u Uint128) Sub(v Uint128) Uint128 {
	lo, b0 := bits.Sub64(u.Lo, v.Lo, 0)
	hi, _ := bits.Sub64(u.Hi, v.Hi, b0)
	return Uint128{Lo: lo, Hi: hi}
}

// Sub64 returns difference (u-v) of 128-bit and 64-bit values.
// Wrap-around semantic is used here: Zero().Sub64(1) == Max().
func (u Uint128) Sub64(v uint64) Uint128 {
	lo, b0 := bits.Sub64(u.Lo, v, 0)
	return Uint128{Lo: lo, Hi: u.Hi - b0}
}

// Mul returns multiplication (u*v) of two 128-bit values.
// Wrap-around semantic is used here: Max().Mul(Max()) == From64(1).
func (u Uint128) Mul(v Uint128) Uint128 {
	hi, lo := bits.Mul64(u.Lo, v.Lo)
	hi += u.Hi*v.Lo + u.Lo*v.Hi
	return Uint128{Lo: lo, Hi: hi}
}

// Mul64 returns multiplication (u*v) of 128-bit and 64-bit values.
// Wrap-around semantic is used here: Max().Mul64(2) == Max().Sub64(1).
func (u Uint128) Mul64(v uint64) Uint128 {
	hi, lo := bits.Mul64(u.Lo, v)
	return Uint128{
		Lo: lo,
		Hi: hi + u.Hi*v,
	}
}

// Div returns division (u/v) of two 128-bit values.
func (u Uint128) Div(v Uint128) Uint128 {
	q, _ := u.QuoRem(v)
	return q
}

// Div64 returns division (u/v) of 128-bit and 64-bit values.
func (u Uint128) Div64(v uint64) Uint128 {
	q, _ := u.QuoRem64(v)
	return q
}

// Mod returns modulo (u%v) of two 128-bit values.
func (u Uint128) Mod(v Uint128) Uint128 {
	_, r := u.QuoRem(v)
	return r
}

// Mod64 returns modulo (u%v) of 128-bit and 64-bit values.
func (u Uint128) Mod64(v uint64) uint64 {
	_, r := u.QuoRem64(v)
	return r
}

// QuoRem returns quotient (u/v) and remainder (u%v) of two 128-bit values.
func (u Uint128) QuoRem(v Uint128) (Uint128, Uint128) {
	if v.Hi == 0 {
		q, r := u.QuoRem64(v.Lo)
		return q, From64(r)
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

	return q, r
}

// QuoRem64 returns quotient (u/v) and remainder (u%v) of 128-bit and 64-bit values.
func (u Uint128) QuoRem64(v uint64) (Uint128, uint64) {
	if u.Hi < v {
		lo, r := bits.Div64(u.Hi, u.Lo, v)
		return Uint128{Lo: lo}, r
	}

	hi, r := bits.Div64(0, u.Hi, v)
	lo, r := bits.Div64(r, u.Lo, v)
	return Uint128{Lo: lo, Hi: hi}, r
}

///////////////////////////////////////////////////////////////////////////////
/// shift operators ///////////////////////////////////////////////////////////

// Lsh returns left shift (u<<n).
func (u Uint128) Lsh(n uint) Uint128 {
	if n > 64 {
		return Uint128{
			// Lo: 0,
			Hi: u.Lo << (n - 64),
		}
	}

	return Uint128{
		Lo: u.Lo << n,
		Hi: u.Hi<<n | u.Lo>>(64-n),
	}
}

// Rsh returns right shift (u>>n).
func (u Uint128) Rsh(n uint) Uint128 {
	if n > 64 {
		return Uint128{
			Lo: u.Hi >> (n - 64),
			// Hi: 0,
		}
	}

	return Uint128{
		Lo: u.Lo>>n | u.Hi<<(64-n),
		Hi: u.Hi >> n,
	}
}

// RotateLeft returns the value of u rotated left by (k mod 128) bits.
func (u Uint128) RotateLeft(k int) Uint128 {
	const n = 128
	s := uint(k) & (n - 1)
	return u.Lsh(s).Or(u.Rsh(n - s)) // TODO: consider to optimize this
}

// RotateRight returns the value of u rotated left by (k mod 128) bits.
func (u Uint128) RotateRight(k int) Uint128 {
	return u.RotateLeft(-k)
}

///////////////////////////////////////////////////////////////////////////////
/// bit counting //////////////////////////////////////////////////////////////

// BitLen returns the minimum number of bits required to represent 128-bit value.
// The result is 0 for u == 0.
func (u Uint128) BitLen() int {
	if u.Hi != 0 {
		return 64 + bits.Len64(u.Hi)
	}
	return bits.Len64(u.Lo)
}

// LeadingZeros returns the number of leading zero bits.
// The result is 128 for u == 0.
func (u Uint128) LeadingZeros() int {
	if u.Hi != 0 {
		return bits.LeadingZeros64(u.Hi)
	}
	return 64 + bits.LeadingZeros64(u.Lo)
}

// TrailingZeros returns the number of trailing zero bits.
// The result is 128 for u == 0.
func (u Uint128) TrailingZeros() int {
	if u.Lo != 0 {
		return bits.TrailingZeros64(u.Lo)
	}
	return 64 + bits.TrailingZeros64(u.Hi)
}

// OnesCount returns the number of one bits ("population count").
func (u Uint128) OnesCount() int {
	return bits.OnesCount64(u.Lo) +
		bits.OnesCount64(u.Hi)
}

// Reverse returns the value with bits in reversed order.
func (u Uint128) Reverse() Uint128 {
	return Uint128{
		Lo: bits.Reverse64(u.Hi),
		Hi: bits.Reverse64(u.Lo),
	}
}

// ReverseBytes returns the value with bytes in reversed order.
func (u Uint128) ReverseBytes() Uint128 {
	return Uint128{
		Lo: bits.ReverseBytes64(u.Hi),
		Hi: bits.ReverseBytes64(u.Lo),
	}
}
