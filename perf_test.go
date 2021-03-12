package uint128

import (
	"math/big"
	"testing"
)

// DummyOutput is exported to avoid unwanted optimizations
var DummyOutput int

// BenchmarkArithmetic performance tests for Add/Sub/Mul/...
func BenchmarkArithmetic(b *testing.B) {
	const K = 1024
	xx := make([]Uint128, K)
	yy := make([]Uint128, K)
	zz := make([]uint, K)
	for i := 0; i < K; i++ {
		xx[i] = rand128()
		yy[i] = rand128()
		zz[i] = uint(yy[i].Lo & 0xFF)
	}

	b.Run("Mul native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo * yy[i%K].Lo
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Add64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Sub", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Sub64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Mul", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Mul64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Lsh", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lsh(zz[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Rsh", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Rsh(zz[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Cmp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp(yy[i%K])
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Cmp64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp64(yy[i%K].Lo)
			DummyOutput += int(res & 1)
		}
	})
}

// BenchmarkDivision performance tests for Div method
func BenchmarkDivision(b *testing.B) {
	const K = 1024
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

	b.Run("native 64/64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx64[i%K].Lo / yy64[i%K].Lo
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Div64 64/64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx64[i%K].Div64(yy64[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Div64 128/64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div64(yy64[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Div 64/64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx64[i%K].Div(yy64[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Div 128/64-Lo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := xx[i%K]
			y := yy64[i%K]
			x.Hi = y.Lo - 1
			res := x.Div(y)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Div 128/64-Hi", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := xx[i%K]
			y := yy64[i%K]
			x.Hi = y.Lo + 1
			res := x.Div(y)
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Div 128/128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("big.Int 128/64", func(b *testing.B) {
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

	b.Run("big.Int 128/128", func(b *testing.B) {
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
