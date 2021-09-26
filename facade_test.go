package bigz_test

import (
	"testing"

	"github.com/Pilatuz/bigz"
)

// note, this is just to make codecov happy.
// all real tests are done in sub-dirs.

// TestUint128 dummy tests for Uint128 helpers.
func TestUint128(t *testing.T) {
	if got := bigz.Zero128().String(); got != "0" {
		t.Errorf("Zero128 failed: %v", got)
	}
	if got := bigz.One128().String(); got != "1" {
		t.Errorf("One128 failed: %v", got)
	}
	if got := bigz.Max128().String(); got != "340282366920938463463374607431768211455" {
		t.Errorf("Max128 failed: %v", got)
	}
}

// TestUint256 dummy tests for Uint256 helpers.
func TestUint256(t *testing.T) {
	if got := bigz.Zero256().String(); got != "0" {
		t.Errorf("Zero256 failed: %v", got)
	}
	if got := bigz.One256().String(); got != "1" {
		t.Errorf("One256 failed: %v", got)
	}
	if got := bigz.Max256().String(); got != "115792089237316195423570985008687907853269984665640564039457584007913129639935" {
		t.Errorf("Max256 failed: %v", got)
	}
}
