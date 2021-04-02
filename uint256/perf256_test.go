package uint256

import (
	"math/big"
	"testing"

	"github.com/Pilatuz/bigx/v2/uint128"
)

// DummyOutput is exported to avoid unwanted optimizations
var DummyOutput int

// BenchmarkAdd performance tests for Add.
func BenchmarkAdd(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand256slice(K)
	yy := rand256slice(K)

	// Native: 64 + 64
	b.Run("Native_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo.Lo + yy[i%K].Lo.Lo
			DummyOutput += int(res & 1)
		}
	})

	// Uint256: 256 + 256
	b.Run("Uint256_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add(yy[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// Uint256: 256 + 128
	b.Run("Uint256_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Add128(yy[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// big.Int: 256 + 256
	b.Run("big.Int_256_256", func(b *testing.B) {
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
	xx := rand256slice(K)
	yy := rand256slice(K)

	// Native: 64 - 64
	b.Run("Native_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo.Lo - yy[i%K].Lo.Lo
			DummyOutput += int(res & 1)
		}
	})

	// Uint256: 256 - 256
	b.Run("Uint256_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub(yy[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// Uint256: 256 - 128
	b.Run("Uint256_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Sub128(yy[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// big.Int: 256 + 256
	b.Run("big.Int_256_256", func(b *testing.B) {
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

// BenchmarkMul performance tests for Mul.
func BenchmarkMul(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand256slice(K)
	yy := rand256slice(K)

	// Native: 64 * 64
	b.Run("Native_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lo.Lo * yy[i%K].Lo.Lo
			DummyOutput += int(res & 1)
		}
	})

	// Mul: 256 * 256
	b.Run("Mul_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hi, lo := Mul(xx[i%K], yy[i%K])
			DummyOutput += int(hi.Lo.Lo&1) + int(lo.Lo.Lo&1)
		}
	})

	// Uint256: 256 * 256
	b.Run("Uint256_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul(yy[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// Uint256: 256 * 128
	b.Run("Uint256_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Mul128(yy[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// big.Int: 256 * 128
	b.Run("big.Int_256_128", func(b *testing.B) {
		xb := make([]*big.Int, K)
		yb := make([]*big.Int, K)
		for i := 0; i < K; i++ {
			xb[i] = xx[i].Big()
			yb[i] = yy[i].Lo.Big()
		}
		q := new(big.Int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			q = q.Mul(xb[i%K], yb[i%K])
		}
		DummyOutput += int(q.Uint64() & 1)
	})

	// big.Int: 256 * 256
	b.Run("big.Int_256_256", func(b *testing.B) {
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
	xx := rand256slice(K)
	yy := rand256slice(K)
	zz := make([]uint, K)
	for i := 0; i < K; i++ {
		zz[i] = uint(yy[i].Lo.Lo & 0xFF)
	}

	b.Run("Uint256.Lsh_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Lsh(zz[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Uint256.Rsh_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Rsh(zz[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Uint256.RotateLeft_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].RotateLeft(int(zz[i%K]))
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Uint256.RotateRight_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].RotateRight(int(zz[i%K]))
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	b.Run("Uint256.Cmp_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp(yy[i%K])
			DummyOutput += int(res & 1)
		}
	})

	b.Run("Uint256.Cmp64_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Cmp128(yy[i%K].Lo)
			DummyOutput += int(res & 1)
		}
	})
}

// BenchmarkDiv performance tests for Div.
func BenchmarkDiv(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand256slice(K)
	yy := rand256slice(K)
	xh := rand256slice(K) // 128-bit half
	yh := rand256slice(K) // 128-bit half
	for i := 0; i < K; i++ {
		xh[i].Hi = uint128.Zero()
		yh[i].Hi = uint128.Zero()

		// avoid zeros
		if yy[i].Lo.IsZero() {
			yy[i].Lo.Lo += 13
		}
		if yh[i].Lo.Lo == 0 {
			yh[i].Lo.Lo += 17
		}
	}

	// native (just as a reference)
	b.Run("Native_64_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xh[i%K].Lo.Lo / yh[i%K].Lo.Lo
			DummyOutput += int(res & 1)
		}
	})

	// Uint256: 256 / 64
	b.Run("Uint256.Div64_256_64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xh[i%K].Div64(yh[i%K].Lo.Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// Uint256: 128 / 128
	b.Run("Uint256.Div128_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xh[i%K].Div128(yh[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// Uint256: 128 / 128
	b.Run("Uint256.Div_128_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xh[i%K].Div(yh[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// Uint256: 256 / 128
	b.Run("Uint256.Div128_256_128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div128(yh[i%K].Lo)
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// Uint256: 256 / 256
	b.Run("Uint256.Div_256_256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res := xx[i%K].Div(yy[i%K])
			DummyOutput += int(res.Lo.Lo & 1)
		}
	})

	// big.Int: 256 / 128
	b.Run("big.Int.Div_256_128", func(b *testing.B) {
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

	// big.Int: 256 / 256
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
