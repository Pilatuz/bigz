package uint256

import (
	"fmt"
	"math/big"

	"github.com/Pilatuz/bigx/v2/uint128"
)

// String returns the base-10 representation of 256-bit value.
func (u Uint256) String() string {
	if u.Hi.IsZero() {
		if u.Lo.IsZero() {
			return "0" // zero
		}
		return u.Lo.String() // lower 128-bit
	}

	buf := []byte("000000000000000000000000000000000000000000000000000000000000000000000000000000") // log10(2^256) < 78
	for i := len(buf); ; i -= 19 {
		q, r := u.QuoRem64(1e19) // largest power of 10 that fits in a uint64
		var n int
		for ; r != 0; r /= 10 {
			n++
			buf[i-n] += byte(r % 10)
		}
		if q.IsZero() {
			return string(buf[i-n:])
		}
		u = q
	}
}

// Format does custom formatting of 256-bit value.
func (u Uint256) Format(s fmt.State, ch rune) {
	u.Big().Format(s, ch) // via big.Int, unefficient! consider to optimize
}

// MarshalText implements the encoding.TextMarshaler interface.
func (u Uint256) MarshalText() (text []byte, err error) {
	return u.Big().MarshalText() // via big.Int, unefficient! consider to optimize
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (u *Uint256) UnmarshalText(text []byte) error {
	// via big.Int, unefficient! consider to optimize
	i := new(big.Int)
	if err := i.UnmarshalText(text); err != nil {
		return err
	}
	v, ok := FromBigX(i)
	if !ok {
		return fmt.Errorf("%q overflows 256-bit integer", text)
	}
	*u = v
	return nil
}

// StoreLittleEndian stores 256-bit value in byte slice in little-endian byte order.
// It panics if byte slice length is less than 32.
func StoreLittleEndian(b []byte, u Uint256) {
	uint128.StoreLittleEndian(b[:16], u.Lo)
	uint128.StoreLittleEndian(b[16:], u.Hi)
}

// StoreBigEndian stores 256-bit value in byte slice in big-endian byte order.
// It panics if byte slice length is less than 32.
func StoreBigEndian(b []byte, u Uint256) {
	uint128.StoreBigEndian(b[:16], u.Hi)
	uint128.StoreBigEndian(b[16:], u.Lo)
}

// LoadLittleEndian loads 256-bit value from byte slice in little-endian byte order.
// It panics if byte slice length is less than 32.
func LoadLittleEndian(b []byte) Uint256 {
	return Uint256{
		Lo: uint128.LoadLittleEndian(b[:16]),
		Hi: uint128.LoadLittleEndian(b[16:]),
	}
}

// LoadBigEndian loads 256-bit value from byte slice in big-endian byte order.
// It panics if byte slice length is less than 32.
func LoadBigEndian(b []byte) Uint256 {
	return Uint256{
		Lo: uint128.LoadBigEndian(b[16:]),
		Hi: uint128.LoadBigEndian(b[:16]),
	}
}
