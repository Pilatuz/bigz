package uint256

import (
	"errors"
	"math/big"
	"math/bits"

	"github.com/Pilatuz/bigx/v2/uint128"
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
		Lo: uint128.Max(),
		Hi: uint128.Max(),
	}
}

// Uint128 is an unsigned 128-bit number alias.
type Uint128 = uint128.Uint128

// Uint256 is an unsigned 256-bit number.
// All methods are immutable, works just like standard uint64.
type Uint256 struct {
	Lo Uint128 // lower 128-bit half
	Hi Uint128 // upper 128-bit half
}

// From128 converts 128-bit value v to a Uint256 value.
// Upper 128-bit half will be zero.
func From128(v Uint128) Uint256 {
	return Uint256{Lo: v}
}

// From64 converts 64-bit value v to a Uint256 value.
// Upper 128-bit half will be zero.
func From64(v uint64) Uint256 {
	return From128(uint128.From64(v))
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

	// Note, actually result of big.Int.Uint64 is undefined
	// if stored value is greater than 2^64
	// but we assume that it just gets lower 64 bits.
	t := new(big.Int)
	lolo := i.Uint64()
	lohi := t.Rsh(i, 64).Uint64()
	hilo := t.Rsh(i, 128).Uint64()
	hihi := t.Rsh(i, 192).Uint64()
	return Uint256{
		Lo: Uint128{Lo: lolo, Hi: lohi},
		Hi: Uint128{Lo: hilo, Hi: hihi},
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

// Equals128 returns true if 256-bit value equals to a 128-bit value.
func (u Uint256) Equals128(v Uint128) bool {
	return u.Lo.Equals(v) && u.Hi.IsZero()
}

// Cmp compares two 256-bit values and returns:
//   -1 if u <  v
//    0 if u == v
//   +1 if u >  v
func (u Uint256) Cmp(v Uint256) int {
	if h := u.Hi.Cmp(v.Hi); h != 0 {
		return h
	}
	return u.Lo.Cmp(v.Lo)
}

// Cmp128 compares 256-bit and 128-bit values and returns:
//   -1 if u <  v
//    0 if u == v
//   +1 if u >  v
func (u Uint256) Cmp128(v Uint128) int {
	switch {
	case !u.Hi.IsZero():
		return +1 // u > v
	}
	return u.Lo.Cmp(v)
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

// AndNot128 returns logical AND NOT (u&v) of 256-bit and 128-bit values.
func (u Uint256) AndNot128(v Uint128) Uint256 {
	return Uint256{
		Lo: u.Lo.AndNot(v),
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

// And128 returns logical AND (u&v) of 256-bit and 128-bit values.
func (u Uint256) And128(v Uint128) Uint256 {
	return Uint256{
		Lo: u.Lo.And(v),
		// Hi: Uint128{0, 0},
	}
}

// Or returns logical OR (u|v) of two 256-bit values.
func (u Uint256) Or(v Uint256) Uint256 {
	return Uint256{
		Lo: u.Lo.Or(v.Lo),
		Hi: u.Hi.Or(v.Hi),
	}
}

// Or128 returns logical OR (u|v) of 256-bit and 128-bit values.
func (u Uint256) Or128(v Uint128) Uint256 {
	return Uint256{
		Lo: u.Lo.Or(v),
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

// Xor128 returns logical XOR (u^v) of 256-bit and 128-bit values.
func (u Uint256) Xor128(v Uint128) Uint256 {
	return Uint256{
		Lo: u.Lo.Xor(v),
		Hi: u.Hi,
	}
}

///////////////////////////////////////////////////////////////////////////////
/// arithmetic operators //////////////////////////////////////////////////////

// Add returns the sum with carry of x, y and carry: sum = x + y + carry.
// The carry input must be 0 or 1; otherwise the behavior is undefined.
// The carryOut output is guaranteed to be 0 or 1.
func Add(x, y Uint256, carry uint64) (sum Uint256, carryOut uint64) {
	sum.Lo, carryOut = uint128.Add(x.Lo, y.Lo, carry)
	sum.Hi, carryOut = uint128.Add(x.Hi, y.Hi, carryOut)
	return
}

// Add returns sum (u+v) of two 256-bit values.
// Wrap-around semantic is used here: Max().Add(From64(1)) == Zero()
func (u Uint256) Add(v Uint256) Uint256 {
	sum, _ := Add(u, v, 0)
	return sum
}

// Add128 returns sum u+v of 256-bit and 128-bit values.
// Wrap-around semantic is used here: Max().Add128(uint128.One()) == Zero()
func (u Uint256) Add128(v Uint128) Uint256 {
	lo, c0 := uint128.Add(u.Lo, v, 0)
	return Uint256{Lo: lo, Hi: u.Hi.Add64(c0)}
}

// Sub returns the difference of x, y and borrow: diff = x - y - borrow.
// The borrow input must be 0 or 1; otherwise the behavior is undefined.
// The borrowOut output is guaranteed to be 0 or 1.
func Sub(x, y Uint256, borrow uint64) (diff Uint256, borrowOut uint64) {
	diff.Lo, borrowOut = uint128.Sub(x.Lo, y.Lo, borrow)
	diff.Hi, borrowOut = uint128.Sub(x.Hi, y.Hi, borrowOut)
	return
}

// Sub returns difference (u-v) of two 256-bit values.
// Wrap-around semantic is used here: Zero().Sub(From64(1)) == Max().
func (u Uint256) Sub(v Uint256) Uint256 {
	diff, _ := Sub(u, v, 0)
	return diff
}

// Sub128 returns difference (u-v) of 256-bit and 128-bit values.
// Wrap-around semantic is used here: Zero().Sub128(uint128.One()) == Max().
func (u Uint256) Sub128(v Uint128) Uint256 {
	lo, b0 := uint128.Sub(u.Lo, v, 0)
	return Uint256{Lo: lo, Hi: u.Hi.Sub64(b0)}
}

// Mul returns the 512-bit product of x and y: (hi, lo) = x * y
// with the product bits' upper half returned in hi and the lower
// half returned in lo.
func Mul(x, y Uint256) (hi, lo Uint256) {
	lo.Hi, lo.Lo = uint128.Mul(x.Lo, y.Lo)
	hi.Hi, hi.Lo = uint128.Mul(x.Hi, y.Hi)
	t0, t1 := uint128.Mul(x.Lo, y.Hi)
	t2, t3 := uint128.Mul(x.Hi, y.Lo)

	var c0, c1 uint64
	lo.Hi, c0 = uint128.Add(lo.Hi, t1, 0)
	lo.Hi, c1 = uint128.Add(lo.Hi, t3, 0)
	hi.Lo, c0 = uint128.Add(hi.Lo, t0, c0)
	hi.Lo, c1 = uint128.Add(hi.Lo, t2, c1)
	hi.Hi = hi.Hi.Add64(c0 + c1)

	return
}

// Mul returns multiplication (u*v) of two 256-bit values.
// Wrap-around semantic is used here: Max().Mul(Max()) == From64(1).
func (u Uint256) Mul(v Uint256) Uint256 {
	hi, lo := uint128.Mul(u.Lo, v.Lo)
	hi = hi.Add(u.Hi.Mul(v.Lo))
	hi = hi.Add(u.Lo.Mul(v.Hi))
	return Uint256{Lo: lo, Hi: hi}
}

// Mul128 returns multiplication (u*v) of 256-bit and 128-bit values.
// Wrap-around semantic is used here: Max().Mul128(2) == Max().Sub128(1).
func (u Uint256) Mul128(v Uint128) Uint256 {
	hi, lo := uint128.Mul(u.Lo, v)
	return Uint256{
		Lo: lo,
		Hi: hi.Add(u.Hi.Mul(v)),
	}
}

// Div returns division (u/v) of two 256-bit values.
func (u Uint256) Div(v Uint256) Uint256 {
	q, _ := u.QuoRem(v)
	return q
}

// Div128 returns division (u/v) of 256-bit and 128-bit values.
func (u Uint256) Div128(v Uint128) Uint256 {
	q, _ := u.QuoRem128(v)
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

// Mod128 returns modulo (u%v) of 256-bit and 128-bit values.
func (u Uint256) Mod128(v Uint128) Uint128 {
	_, r := u.QuoRem128(v)
	return r
}

// Mod64 returns modulo (u%v) of 256-bit and 64-bit values.
func (u Uint256) Mod64(v uint64) uint64 {
	_, r := u.QuoRem64(v)
	return r
}

// QuoRem returns quotient (u/v) and remainder (u%v) of two 256-bit values.
func (u Uint256) QuoRem(v Uint256) (Uint256, Uint256) {
	if v.Hi.IsZero() {
		q, r := u.QuoRem128(v.Lo)
		return q, From128(r)
	}

	// generate a "trial quotient," guaranteed to be
	// within 1 of the actual quotient, then adjust.
	n := uint(v.Hi.LeadingZeros())
	u1, v1 := u.Rsh(1), v.Lsh(n)
	tq, _ := uint128.Div(u1.Hi, u1.Lo, v1.Hi)
	tq = tq.Rsh(127 - n)
	if !tq.IsZero() {
		tq = tq.Sub64(1)
	}

	// calculate remainder using trial quotient, then
	// adjust if remainder is greater than divisor
	q, r := From128(tq), u.Sub(v.Mul128(tq))
	if r.Cmp(v) >= 0 {
		q = q.Add128(uint128.One())
		r = r.Sub(v)
	}

	return q, r
}

// QuoRem128 returns quotient (u/v) and remainder (u%v) of 256-bit and 128-bit values.
func (u Uint256) QuoRem128(v Uint128) (Uint256, Uint128) {
	if u.Hi.Cmp(v) < 0 {
		lo, r := uint128.Div(u.Hi, u.Lo, v)
		return Uint256{Lo: lo}, r
	}

	hi, r := uint128.Div(uint128.Zero(), u.Hi, v)
	lo, r := uint128.Div(r, u.Lo, v)
	return Uint256{Lo: lo, Hi: hi}, r
}

// QuoRem64 returns quotient (u/v) and remainder (u%v) of 256-bit and 64-bit values.
func (u Uint256) QuoRem64(v uint64) (q Uint256, r uint64) {
	q.Hi, r = u.Hi.QuoRem64(v)
	q.Lo.Hi, r = bits.Div64(r, u.Lo.Hi, v)
	q.Lo.Lo, r = bits.Div64(r, u.Lo.Lo, v)
	return
}

// Div returns the quotient and remainder of (hi, lo) divided by y:
// quo = (hi, lo)/y, rem = (hi, lo)%y with the dividend bits' upper
// half in parameter hi and the lower half in parameter lo.
// Panics if y is less or equal to hi!
func Div(hi, lo, y Uint256) (quo, rem Uint256) {
	if y.IsZero() {
		panic(errors.New("integer divide by zero"))
	}
	if y.Cmp(hi) <= 0 {
		panic(errors.New("integer overflow"))
	}

	s := uint(y.LeadingZeros())
	y = y.Lsh(s)

	un32 := hi.Lsh(s).Or(lo.Rsh(256 - s))
	un10 := lo.Lsh(s)
	q1, rhat := un32.QuoRem128(y.Hi)
	r1 := From128(rhat)

	for !q1.Hi.IsZero() || q1.Mul128(y.Lo).Cmp(Uint256{Hi: r1.Lo, Lo: un10.Hi}) > 0 {
		q1 = q1.Sub128(uint128.One())
		r1 = r1.Add128(y.Hi)
		if !r1.Hi.IsZero() {
			break
		}
	}

	un21 := Uint256{Hi: un32.Lo, Lo: un10.Hi}.Sub(q1.Mul(y))
	q0, rhat := un21.QuoRem128(y.Hi)
	r0 := From128(rhat)

	for !q0.Hi.IsZero() || q0.Mul128(y.Lo).Cmp(Uint256{Hi: r0.Lo, Lo: un10.Lo}) > 0 {
		q0 = q0.Sub128(uint128.One())
		r0 = r0.Add128(y.Hi)
		if !r0.Hi.IsZero() {
			break
		}
	}

	return Uint256{Hi: q1.Lo, Lo: q0.Lo},
		Uint256{Hi: un21.Lo, Lo: un10.Lo}.
			Sub(q0.Mul(y)).Rsh(s)
}

///////////////////////////////////////////////////////////////////////////////
/// shift operators ///////////////////////////////////////////////////////////

// Lsh returns left shift (u<<n).
func (u Uint256) Lsh(n uint) Uint256 {
	if n > 128 {
		return Uint256{
			// Lo: Uint128{Lo: 0, Hi: 0},
			Hi: u.Lo.Lsh(n - 128),
		}
	}

	if n > 64 {
		n -= 64
		return Uint256{
			Lo: Uint128{
				// Lo: 0,
				Hi: u.Lo.Lo << n,
			},
			Hi: Uint128{
				Lo: u.Lo.Hi<<n | u.Lo.Lo>>(64-n),
				Hi: u.Hi.Lo<<n | u.Lo.Hi>>(64-n),
			},
		}
	}

	return Uint256{
		Lo: Uint128{
			Lo: u.Lo.Lo << n,
			Hi: u.Lo.Hi<<n | u.Lo.Lo>>(64-n),
		},
		Hi: Uint128{
			Lo: u.Hi.Lo<<n | u.Lo.Hi>>(64-n),
			Hi: u.Hi.Hi<<n | u.Hi.Lo>>(64-n),
		},
	}
}

// Rsh returns right shift (u>>n).
func (u Uint256) Rsh(n uint) Uint256 {
	if n > 128 {
		return Uint256{
			Lo: u.Hi.Rsh(n - 128),
			// Hi: Uint128{Lo: 0, Hi: 0},
		}
	}

	if n > 64 {
		n -= 64
		return Uint256{
			Lo: Uint128{
				Lo: u.Lo.Hi>>n | u.Hi.Lo<<(64-n),
				Hi: u.Hi.Lo>>n | u.Hi.Hi<<(64-n),
			},
			Hi: Uint128{
				Lo: u.Hi.Hi >> n,
				// Hi: 0,
			},
		}
	}

	return Uint256{
		Lo: Uint128{
			Lo: u.Lo.Lo>>n | u.Lo.Hi<<(64-n),
			Hi: u.Lo.Hi>>n | u.Hi.Lo<<(64-n),
		},
		Hi: Uint128{
			Lo: u.Hi.Lo>>n | u.Hi.Hi<<(64-n),
			Hi: u.Hi.Hi >> n,
		},
	}
}

// RotateLeft returns the value of u rotated left by (k mod 256) bits.
func (u Uint256) RotateLeft(k int) Uint256 {
	n := uint(k) & 255

	if n < 64 {
		if n == 0 {
			// no shift
			return u
		}

		// shift by [1..63]
		return Uint256{
			Lo: Uint128{
				Lo: u.Lo.Lo<<n | u.Hi.Hi>>(64-n),
				Hi: u.Lo.Hi<<n | u.Lo.Lo>>(64-n),
			},
			Hi: Uint128{
				Lo: u.Hi.Lo<<n | u.Lo.Hi>>(64-n),
				Hi: u.Hi.Hi<<n | u.Hi.Lo>>(64-n),
			},
		}
	}

	n -= 64
	if n < 64 {
		if n == 0 {
			// shift by 64
			return Uint256{
				Lo: Uint128{
					Lo: u.Hi.Hi,
					Hi: u.Lo.Lo,
				},
				Hi: Uint128{
					Lo: u.Lo.Hi,
					Hi: u.Hi.Lo,
				},
			}
		}

		// shift by [65..127]
		return Uint256{
			Lo: Uint128{
				Lo: u.Hi.Hi<<n | u.Hi.Lo>>(64-n),
				Hi: u.Lo.Lo<<n | u.Hi.Hi>>(64-n),
			},
			Hi: Uint128{
				Lo: u.Lo.Hi<<n | u.Lo.Lo>>(64-n),
				Hi: u.Hi.Lo<<n | u.Lo.Hi>>(64-n),
			},
		}
	}

	n -= 64
	if n < 64 {
		if n == 0 {
			// shift by 128
			return Uint256{
				Lo: u.Hi,
				Hi: u.Lo,
			}
		}

		// shift by [129..191]
		return Uint256{
			Lo: Uint128{
				Lo: u.Hi.Lo<<n | u.Lo.Hi>>(64-n),
				Hi: u.Hi.Hi<<n | u.Hi.Lo>>(64-n),
			},
			Hi: Uint128{
				Lo: u.Lo.Lo<<n | u.Hi.Hi>>(64-n),
				Hi: u.Lo.Hi<<n | u.Lo.Lo>>(64-n),
			},
		}
	}

	n -= 64
	if n == 0 {
		// shift by 192
		return Uint256{
			Lo: Uint128{
				Lo: u.Lo.Hi,
				Hi: u.Hi.Lo,
			},
			Hi: Uint128{
				Lo: u.Hi.Hi,
				Hi: u.Lo.Lo,
			},
		}
	}

	// shift by [193..255]
	return Uint256{
		Lo: Uint128{
			Lo: u.Lo.Hi<<n | u.Lo.Lo>>(64-n),
			Hi: u.Hi.Lo<<n | u.Lo.Hi>>(64-n),
		},
		Hi: Uint128{
			Lo: u.Hi.Hi<<n | u.Hi.Lo>>(64-n),
			Hi: u.Lo.Lo<<n | u.Hi.Hi>>(64-n),
		},
	}
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
