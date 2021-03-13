package uint128

import (
	"math/big"
	"testing"
)

// DummyOutput is exported to avoid unwanted optimizations
var DummyOutput int

// BenchmarkArithmetic performance tests for Add/Sub/Mul/...
func BenchmarkArithmetic(b *testing.B) {
	// prepare a set of values
	// just not to work with the same number
	const K = 1024 // should be power of 2
	xx := make([]Uint128, K)
	yy := make([]Uint128, K)
	zz := make([]uint, K)
	for i := 0; i < K; i++ {
		xx[i] = rand128()
		yy[i] = rand128()
		zz[i] = uint(yy[i].Lo & 0xFF)
	}

	// native (just as a reference)
	b.Run("Mul_64_64_native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo * yy[i%K].Lo
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Add_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Add64_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Sub_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Sub64_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Mul_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Mul64_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Lsh_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lsh(zz[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Rsh_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Rsh(zz[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("RotateLeft_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].RotateLeft(int(zz[i%K]))
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("RotateRight_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].RotateRight(int(zz[i%K]))
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Cmp_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp(yy[i%K])
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Cmp64_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp64(yy[i%K].Lo)
			DummyOutput += int(res & 1)
		}
	})
}

// BenchmarkDivision performance tests for Div method
func BenchmarkDivision(b *testing.B) {
	// prepare a set of values
	// just not to work with the same number
	const K = 1024 // should be power of 2
	xx := make([]Uint128, K)
	yy := make([]Uint128, K)
	xx64 := make([]Uint128, K)
	yy64 := make([]Uint128, K)
	for i := 0; i < K; i++ {
		xx[i] = rand128()
		yy[i] = rand128()
		xx64[i] = rand128()
		yy64[i] = rand128()

		xx64[i].Hi = 0
		yy64[i].Hi = 0

		// avoid zeros
		if yy[i].IsZero() {
			yy[i].Lo += 13
		}
		if yy64[i].IsZero() {
			yy64[i].Lo += 17
		}
	}

	// native (just as a reference)
	b.Run("Div_64_64_native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx64[i%K].Lo / yy64[i%K].Lo
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Div64_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx64[i%K].Div64(yy64[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Div64_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div64(yy64[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Div_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx64[i%K].Div(yy64[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Div_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("big.Int.Div_128_64", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = yy64[i].Big()
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Div(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})

	b.Run("big.Int.Div_128_128", func(b *testing.B) {
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
