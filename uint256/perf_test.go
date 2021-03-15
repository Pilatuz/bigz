package uint256

import (
	"math/big"
	"testing"

	"github.com/Pilatuz/bigx/uint128"
)

// DummyOutput is exported to avoid unwanted optimizations
var DummyOutput int

// BenchmarkArithmetic performance tests for Add/Sub/Mul/...
func BenchmarkArithmetic(b *testing.B) {
	// prepare a set of values
	// just not to work with the same number
	const K = 1024 // should be power of 2
	xx := make([]Uint256, K)
	yy := make([]Uint256, K)
	zz := make([]uint, K)
	for i := 0; i < K; i++ {
		xx[i] = rand256()
		yy[i] = rand256()
		zz[i] = uint(yy[i].Lo.Lo & 0xFF)
	}

	// native (just as a reference)
	b.Run("Mul_64_64_native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo.Lo * yy[i%K].Lo.Lo
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Add_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add(yy[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Add64_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add128(yy[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Sub_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub(yy[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Sub64_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub128(yy[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Mul_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul(yy[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Mul64_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul128(yy[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Lsh_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lsh(zz[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Rsh_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Rsh(zz[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("RotateLeft_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].RotateLeft(int(zz[i%K]))
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("RotateRight_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].RotateRight(int(zz[i%K]))
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Cmp_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp(yy[i%K])
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Cmp64_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp128(yy[i%K].Lo)
			DummyOutput += int(res & 1)
		}
	})
}

// BenchmarkDivision performance tests for Div method
func BenchmarkDivision(b *testing.B) {
	// prepare a set of values
	// just not to work with the same number
	const K = 1024 // should be power of 2
	xx := make([]Uint256, K)
	yy := make([]Uint256, K)
	xx128 := make([]Uint256, K)
	yy128 := make([]Uint256, K)
	for i := 0; i < K; i++ {
		xx[i] = rand256()
		yy[i] = rand256()
		xx128[i] = rand256()
		yy128[i] = rand256()

		xx128[i].Hi = uint128.Zero()
		yy128[i].Hi = uint128.Zero()

		// avoid zeros
		if yy[i].Lo.IsZero() {
			yy[i].Lo.Lo += 13
		}
		if yy128[i].Lo.Lo == 0 {
			yy128[i].Lo.Lo += 17
		}
	}

	// native (just as a reference)
	b.Run("Div_64_64_native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx128[i%K].Lo.Lo / yy128[i%K].Lo.Lo
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Div128_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx128[i%K].Div128(yy128[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Div128_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div128(yy128[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Div_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx128[i%K].Div(yy128[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Div_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div(yy[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("big.Int.Div_256_128", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = yy128[i].Big()
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Div(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})

	b.Run("big.Int.Div_256_256", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = yy[i].Big()
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Div(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})
}
