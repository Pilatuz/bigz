package bigx

import (
	u128 "github.com/Pilatuz/bigx/v2/uint128"
)

// Uint128 is type alias for 128-bit unsigned integer.
type Uint128 = u128.Uint128

// Note, there in no New(lo, hi) just not to confuse
// which half goes first: lower or upper.
// Use structure initialization Uint128{Lo: ..., Hi: ...} instead.

// Zero128 is the lowest possible Uint128 value.
func Zero128() Uint128 {
	return u128.Zero()
}

// One128 is the lowest non-zero Uint128 value.
func One128() Uint128 {
	return u128.One()
}

// Max128 is the largest possible Uint128 value.
func Max128() Uint128 {
	return u128.Max()
}
