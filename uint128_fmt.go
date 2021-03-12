package uint128

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// String returns the base-10 representation of 128-bit value.
func (u Uint128) String() string {
	if u.Hi == 0 {
		if u.Lo == 0 {
			return "0" // zero
		}
		return strconv.FormatUint(u.Lo, 10) // lower 64-bit
	}

	buf := []byte("0000000000000000000000000000000000000000") // log10(2^128) < 40
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

// Format does custom formatting of 128-bit value.
func (u Uint128) Format(s fmt.State, ch rune) {
	// do it in the simplest
	// but unefficient way for now
	u.Big().Format(s, ch)
}

// StoreUint128LE stores 128-bit value in byte slice in little-endian order.
// It panics if byte slice length is less than 16.
func StoreUint128LE(b []byte, u Uint128) {
	binary.LittleEndian.PutUint64(b[:8], u.Lo)
	binary.LittleEndian.PutUint64(b[8:], u.Hi)
}

// StoreUint128BE stores 128-bit value in byte slice in big-endian order.
// It panics if byte slice length is less than 16.
func StoreUint128BE(b []byte, u Uint128) {
	binary.BigEndian.PutUint64(b[:8], u.Hi)
	binary.BigEndian.PutUint64(b[8:], u.Lo)
}

// LoadUint128LE loads 128-bit value from byte slice in little-endian order.
// It panics if byte slice length is less than 16.
func LoadUint128LE(b []byte) Uint128 {
	return Uint128{
		Lo: binary.LittleEndian.Uint64(b[:8]),
		Hi: binary.LittleEndian.Uint64(b[8:]),
	}
}

// LoadUint128BE loads 128-bit value from byte slice in big-endian order.
// It panics if byte slice length is less than 16.
func LoadUint128BE(b []byte) Uint128 {
	return Uint128{
		Lo: binary.BigEndian.Uint64(b[8:]),
		Hi: binary.BigEndian.Uint64(b[:8]),
	}
}
