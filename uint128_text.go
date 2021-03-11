package uint128

import (
	"encoding/binary"
)

// TODO: format

// String returns the base-10 representation of u as a string.
func (u Uint128) String() string {
	if u.IsZero() {
		return "0"
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

// PutBytes stores u in b in little-endian order. It panics if len(b) < 16.
func (u Uint128) PutBytes(b []byte) {
	// little-endian order is used by default
	u.PutBytesLE(b)
}

// PutBytesLE stores u in b in little-endian order. It panics if len(b) < 16.
func (u Uint128) PutBytesLE(b []byte) {
	binary.LittleEndian.PutUint64(b[:8], u.Lo)
	binary.LittleEndian.PutUint64(b[8:], u.Hi)
}

// PutBytesBE stores u in b in big-endian order. It panics if len(b) < 16.
func (u Uint128) PutBytesBE(b []byte) {
	binary.BigEndian.PutUint64(b[:8], u.Hi)
	binary.BigEndian.PutUint64(b[8:], u.Lo)
}

// FromBytes converts b to a Uint128 value (little-endian order).
func FromBytes(b []byte) Uint128 {
	// little-endian order is used by default
	return FromBytesLE(b)
}

// FromBytesLE converts b to a Uint128 value (little-endian order).
func FromBytesLE(b []byte) Uint128 {
	return Uint128{
		Lo: binary.LittleEndian.Uint64(b[:8]),
		Hi: binary.LittleEndian.Uint64(b[8:]),
	}
}

// FromBytesBE converts b to a Uint128 value (big-endian order).
func FromBytesBE(b []byte) Uint128 {
	return Uint128{
		Lo: binary.BigEndian.Uint64(b[8:]),
		Hi: binary.BigEndian.Uint64(b[:8]),
	}
}
