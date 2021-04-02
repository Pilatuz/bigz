package uint128

import (
	"math/big"
	"testing"
)

// DummyOutput is exported to avoid unwanted optimizations
var DummyOutput int

// BenchmarkAdd performance tests for Add.
func BenchmarkAdd(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand128slice(K)
	yy := rand128slice(K)

	// Native: 64 + 64
	b.Run("Native_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo + yy[i%K].Lo
			DummyOutput += int(res & 1)
		}
	})

	// Uint128: 128 + 128
	b.Run("Uint128_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	// Uint128: 128 + 64
	b.Run("Uint128_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	// big.Int: 128 + 128
	b.Run("big.Int_128_128", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = yy[i].Big()
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Add(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})
}

// BenchmarkSub performance tests for Sub.
func BenchmarkSub(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand128slice(K)
	yy := rand128slice(K)

	// Native: 64 - 64
	b.Run("Native_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo - yy[i%K].Lo
			DummyOutput += int(res & 1)
		}
	})

	// Uint128: 128 - 128
	b.Run("Uint128_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	// Uint128: 128 - 64
	b.Run("Uint128_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	// big.Int: 128 - 128
	b.Run("big.Int_128_128", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = yy[i].Big()
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Sub(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})
}

// BenchmarkMul performance tests for Sub.
func BenchmarkMul(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand128slice(K)
	yy := rand128slice(K)

	// Native: 64 * 64
	b.Run("Native_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo * yy[i%K].Lo
			DummyOutput += int(res & 1)
		}
	})

	// Mul: 128 * 128
	b.Run("Mul_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hi, lo := Mul(xx[i%K], yy[i%K])
			DummyOutput += int(hi.Lo&1) + int(lo.Lo&1)
		}
	})

	// Uint128: 128 * 128
	b.Run("Uint128_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	// Uint128: 128 * 64
	b.Run("Uint128_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul64(yy[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	// big.Int: 128 * 64
	b.Run("big.Int_128_64", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = new(big.Int).SetUint64(yy[i].Lo)
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Mul(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})

	// big.Int: 128 + 128
	b.Run("big.Int_128_128", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = yy[i].Big()
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Mul(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})
}

// BenchmarkMisc performance tests for Lsh, Rsh, Cmp, etc.
func BenchmarkMisc(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand128slice(K)
	yy := rand128slice(K)
	zz := make([]uint, K)
	for i := 0; i < K; i++ {
		zz[i] = uint(yy[i].Lo & 0xFF)
	}

	b.Run("Uint128.Lsh_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lsh(zz[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Uint128.Rsh_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Rsh(zz[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Uint128.RotateLeft_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].RotateLeft(int(zz[i%K]))
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Uint128.RotateRight_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].RotateRight(int(zz[i%K]))
			DummyOutput += int(res.Lo & 1)
		}
	})

	b.Run("Uint128.Cmp_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp(yy[i%K])
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Uint128.Cmp64_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp64(yy[i%K].Lo)
			DummyOutput += int(res & 1)
		}
	})
}

// BenchmarkDiv performance tests for Div.
func BenchmarkDiv(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand128slice(K)
	yy := rand128slice(K)
	xh := rand128slice(K) // 64-bit half
	yh := rand128slice(K) // 64-bit half
	for i := 0; i < K; i++ {
		xh[i].Hi = 0
		yh[i].Hi = 0

		// avoid zeros
		if yy[i].IsZero() {
			yy[i].Lo += 13
		}
		if yh[i].IsZero() {
			yh[i].Lo += 17
		}
	}

	// native (just as a reference)
	b.Run("Native_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xh[i%K].Lo / yh[i%K].Lo
			DummyOutput += int(res & 1)
		}
	})

	// Uint128: 64 / 64
	b.Run("Uint128.Div64_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xh[i%K].Div64(yh[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	// Uint128: 64 / 64
	b.Run("Uint128.Div_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xh[i%K].Div(yh[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	// Uint128: 128 / 64
	b.Run("Uint128.Div64_128_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div64(yh[i%K].Lo)
			DummyOutput += int(res.Lo & 1)
		}
	})

	// Uint128: 128 / 128
	b.Run("Uint128.Div_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div(yy[i%K])
			DummyOutput += int(res.Lo & 1)
		}
	})

	// big.Int: 128 / 64
	b.Run("big.Int_128_64", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = yh[i].Big()
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Div(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})

	// big.Int: 128 / 128
	b.Run("big.Int_128_128", func(b *testing.B) {
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
