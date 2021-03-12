package uint128

import (
	"fmt"
	"testing"
)

// TestUint128String unit tests for Uint128.String() method
func TestUint128String(t *testing.T) {
	t.Run("manual", func(t *testing.T) {
		// Zero()
		if expected, got := "0", Zero().String(); got != expected {
			t.Errorf(`Zero() should be %q, got %q`, expected, got)
		}

		// One()
		if expected, got := "1", One().String(); got != expected {
			t.Errorf(`One() should be %q, got %q`, expected, got)
		}

		// Max()
		if expected, got := "340282366920938463463374607431768211455", Max().String(); got != expected {
			t.Errorf(`Max() should be %q, got %q`, expected, got)
		}
	})

	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint128)
		go generate128s(1000, values)
		for x := range values {
			if expected, got := x.Big().String(), x.String(); got != expected {
				t.Errorf("String() mismatch:\n\t(-) expected %q\n\t(+)   actual %q", expected, got)
			}
		}
	})
}

// BenchmarkUint128String performance tests for Uint128.String() method
func BenchmarkUint128String(b *testing.B) {
	b.ReportAllocs()

	x := rand128()
	xb := x.Big()

	b.Run("Uint128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = x.String()
		}
	})

	b.Run("big.Int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xb.String()
		}
	})
}

// TestUint128Format unit tests for Uint128.Format() method
func TestUint128Format(t *testing.T) {
	t.Run("manual", func(t *testing.T) {
		// Zero()
		if expected, got := "0o0", fmt.Sprintf("%#O", Zero()); got != expected {
			t.Errorf(`Zero() should be %q, got %q`, expected, got)
		}

		// One()
		if expected, got := "0001", fmt.Sprintf("%04b", One()); got != expected {
			t.Errorf(`One() should be %q, got %q`, expected, got)
		}

		// Max()
		if expected, got := "ffffffffffffffffffffffffffffffff", fmt.Sprintf("%x", Max()); got != expected {
			t.Errorf(`Max() should be %q, got %q`, expected, got)
		}
	})
}

// TestStoreLoad unit tests for bytes load/store functions
func TestStoreLoad(t *testing.T) {
	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint128)
		go generate128s(1000, values)
		for x := range values {
			buf := make([]byte, 16)

			// little-endian
			StoreUint128LE(buf, x)
			if got := LoadUint128LE(buf); got != x {
				t.Errorf("LoadUint128LE is not the inverse of StoreUint128LE for %#x, got %#x", x, got)
			}

			// big-endian
			StoreUint128BE(buf, x)
			if got := LoadUint128BE(buf); got != x {
				t.Errorf("LoadUint128BE is not the inverse of StoreUint128BE for %#x, got %#x", x, got)
			}

			// reverse bytes
			if got := LoadUint128LE(buf); got != x.ReverseBytes() {
				t.Errorf("LoadUint128LE is not the inverse of StoreUint128BE.ReverseBytes for %#x, got %#x", x, got)
			}
		}
	})
}
