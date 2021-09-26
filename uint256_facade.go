package bigz

import (
	u256 "github.com/Pilatuz/bigz/uint256"
)

// Uint256 is type alias for 256-bit unsigned integer.
type Uint256 = u256.Uint256

// Note, there in no New(lo, hi) just not to confuse
// which half goes first: lower or upper.
// Use structure initialization Uint256{Lo: ..., Hi: ...} instead.

// Zero256 is the lowest possible Uint256 value.
func Zero256() Uint256 {
	return u256.Zero()
}

// One256 is the lowest non-zero Uint256 value.
func One256() Uint256 {
	return u256.One()
}

// Max256 is the largest possible Uint256 value.
func Max256() Uint256 {
	return u256.Max()
}
