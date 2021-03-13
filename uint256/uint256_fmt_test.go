package uint256

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestUint128String unit tests for Uint128.String() method
func TestUint128String(t *testing.T) {
	t.Run("manual", func(t *testing.T) {
		// Zero()
		if expected, got := "0", Zero().String(); got != expected {
			t.Errorf("Zero() should be %q, got %q", expected, got)
		}

		// One()
		if expected, got := "1", One().String(); got != expected {
			t.Errorf("One() should be %q, got %q", expected, got)
		}

		// Max()
		if expected, got := "115792089237316195423570985008687907853269984665640564039457584007913129639935", Max().String(); got != expected {
			t.Errorf("Max() should be %q, got %q", expected, got)
		}
	})

	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint256)
		go generate256s(1000, values)
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

	x := rand256()
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
			t.Errorf("Zero() should be %q, got %q", expected, got)
		}

		// One()
		if expected, got := "0001", fmt.Sprintf("%04b", One()); got != expected {
			t.Errorf("One() should be %q, got %q", expected, got)
		}

		// Max()
		if expected, got := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", fmt.Sprintf("%x", Max()); got != expected {
			t.Errorf("Max() should be %q, got %q", expected, got)
		}
	})
}

// TestStoreLoad unit tests for bytes load/store functions
func TestStoreLoad(t *testing.T) {
	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint256)
		go generate256s(1000, values)
		for x := range values {
			buf := make([]byte, 32)

			// little-endian
			StoreLittleEndian(buf, x)
			if got := LoadLittleEndian(buf); got != x {
				t.Errorf("LoadLittleEndian is not the inverse of StoreLittleEndian for %#x, got %#x", x, got)
			}

			// big-endian
			StoreBigEndian(buf, x)
			if got := LoadBigEndian(buf); got != x {
				t.Errorf("LoadBigEndian is not the inverse of StoreBigEndian for %#x, got %#x", x, got)
			}

			// reverse bytes
			if got := LoadLittleEndian(buf); got != x.ReverseBytes() {
				t.Errorf("LoadLittleEndian is not the inverse of StoreBigEndian.ReverseBytes for %#x, got %#x", x, got)
			}
		}
	})
}

// TestJSON unit tests for marshaling functions
func TestJSON(t *testing.T) {
	type Foo struct {
		Bar Uint256 `json:"bar"`
	}

	t.Run("bad", func(t *testing.T) {
		var tmp Foo

		// expected non-empty string
		err := json.Unmarshal([]byte(`{"bar":""}`), &tmp)
		if err == nil {
			t.Fatalf("should fail on BAD JSON")
		}

		// expected positive integer in range [0, 2^128)
		err = json.Unmarshal([]byte(`{"bar":"-1"}`), &tmp)
		if err == nil {
			t.Fatalf("should fail on BAD JSON")
		}
	})

	t.Run("rand", func(t *testing.T) {
		values := make(chan Uint256)
		go generate256s(1000, values)
		for x := range values {
			buf, err := json.Marshal(Foo{Bar: x})
			if err != nil {
				t.Fatalf("failed to marshal to JSON: %v", err)
			}

			var tmp Foo
			err = json.Unmarshal(buf, &tmp)
			if err != nil {
				t.Fatalf("failed to unmarshal JSON: %v", err)
			}

			if got := tmp.Bar; !got.Equals(x) {
				t.Fatalf("%#x does not equal itself after JSON decoding, got: %#x", x, got)
			}
		}
	})
}
