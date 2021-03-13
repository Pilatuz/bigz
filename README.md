# uint128 [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][reportcard-img]][reportcard]

`uint128` provides a high-performance `Uint128` type that supports standard arithmetic
operations. Unlike `math/big`, operations on `Uint128` always produce new values
instead of modifying a pointer receiver. A `Uint128` value is therefore immutable, just
like `uint64` and friends.

Released under the [MIT License](LICENSE).


## Installation

```shell
go get github.com/Pilatuz/uint128
```

The name `uint128.Uint128` stutters, so it recommended either using a "import alias":

```go
import (
    bigx "github.com/Pilatuz/uint128"
)

// use it as bigx.Uint128
```

or type aliasing `uint128.Uint128` to give it a project-specific name:

```go
import (
    "github.com/Pilatuz/uint128"
)

type Uint128 = uint128.Uint128
```


## What's new

The key differences from [original package](https://github.com/lukechampine/uint128):

- No panics! All methods have wrap-around semantic!
- `Zero` and `Max` are functions to prevent modification of global variables.
- `New` was removed to encourage explicit `Uint128{Lo: ..., Hi: ...}` initialization.
- Trivial (via corresponding `big.Int.Format`) implementation of `Format` method to support for example hex output as `fmt.Sprintf("%X", u)`.
- Store/Load methods in little-endian and big-endian byte order.
- New `Not` and `AndNot` methods.


## Quick Start

The 128-bit integer can be initialized in the following ways:

| Method                             | Description                                       |
|------------------------------------|---------------------------------------------------|
| `u := Uint128{Lo: lo64, Hi: hi64}` | Set both low and high 64-bit halfs.               |
| `u := From64(lo64)`                | Set only low 64-bit half.                         |
| `u := Zero()`                      | The same as `From64(0)`.                          |
| `u := One()`                       | The same as `From64(1)`.                          |
| `u := Max()`                       | The largest possible 128-bit value (`2^128 - 1`). |
| `u := FromBig(big)`                | Convert from `*big.Int` with saturation.          |
| `u := FromBigX(big)`               | The same as `FromBig` but provides `ok` flag.     |

The following arithmetic operations are supported:

| 128-bit    | 64-bit       | Standard `*big.Int` equivalent                                  |
|------------|--------------|-----------------------------------------------------------------|
| `u.Add`    | `u.Add64`    | [`big.Int.Add`](https://golang.org/pkg/math/big/#Int.Add)       |
| `u.Sub`    | `u.Sub64`    | [`big.Int.Sub`](https://golang.org/pkg/math/big/#Int.Sub)       |
| `u.Mul`    | `u.Mul64`    | [`big.Int.Mul`](https://golang.org/pkg/math/big/#Int.Mul)       |
| `u.Div`    | `u.Div64`    | [`big.Int.Div`](https://golang.org/pkg/math/big/#Int.Div)       |
| `u.Mod`    | `u.Mod64`    | [`big.Int.Mod`](https://golang.org/pkg/math/big/#Int.Mod)       |
| `u.QuoRem` | `u.QuoRem64` | [`big.Int.QuoRem`](https://golang.org/pkg/math/big/#Int.QuoRem) |

The following logical and comparison operations are supported:

| 128-bit    | 64-bit       | Standard `*big.Int` equivalent                                  |
|------------|--------------|-----------------------------------------------------------------|
| `u.Equals` | `u.Equals64` | [`big.Int.Cmp == 0`](https://golang.org/pkg/math/big/#Int.Cmp)  |
| `u.Cmp`    | `u.Cmp64`    | [`big.Int.Cmp`](https://golang.org/pkg/math/big/#Int.Cmp)       |
| `u.Not`    |              | [`big.Int.Not`](https://golang.org/pkg/math/big/#Int.Not)       |
| `u.AndNot` | `u.AndNot64` | [`big.Int.AndNot`](https://golang.org/pkg/math/big/#Int.AndNot) |
| `u.And`    | `u.And64`    | [`big.Int.And`](https://golang.org/pkg/math/big/#Int.And)       |
| `u.Or`     | `u.Or64`     | [`big.Int.Or`](https://golang.org/pkg/math/big/#Int.Or)         |
| `u.Xor`    | `u.Xor64`    | [`big.Int.Xor`](https://golang.org/pkg/math/big/#Int.Xor)       |
| `u.Lsh`    |              | [`big.Int.Lsh`](https://golang.org/pkg/math/big/#Int.Lsh)       |
| `u.Rsh`    |              | [`big.Int.Rsh`](https://golang.org/pkg/math/big/#Int.Rsh)       |

The following bit operations are supported:

| 128-bit           | Standard 64-bit equivalent                                                  |
|-------------------|-----------------------------------------------------------------------------|
| `u.RotateLeft`    | [`bits.RotateLeft64`](https://golang.org/pkg/math/bits/#RotateLeft64)       |
| `u.RotateRight`   | [`bits.RotateRight64`](https://golang.org/pkg/math/bits/#RotateRight64)     |
| `u.BitLen`        | [`bits.Len64`](https://golang.org/pkg/math/bits/#Len64) or [`big.Int.BitLen`](https://golang.org/pkg/math/big/#Int.BitLen) |
| `u.LeadingZeros`  | [`bits.LeadingZeros64`](https://golang.org/pkg/math/bits/#LeadingZeros64)   |
| `u.TrailingZeros` | [`bits.TrailingZeros64`](https://golang.org/pkg/math/bits/#TrailingZeros64) |
| `u.OnesCount`     | [`bits.OnesCount64`](https://golang.org/pkg/math/bits/#OnesCount64)         |
| `u.Reverse`       | [`bits.Reverse64`](https://golang.org/pkg/math/bits/#Reverse64)             |
| `u.ReverseBytes`  | [`bits.ReverseBytes64`](https://golang.org/pkg/math/bits/#ReverseBytes64)   |

The following miscellaneous operations are supported:

| 128-bit            | Standard equivalent                                                                  |
|--------------------|--------------------------------------------------------------------------------------|
| `u.String`         | [`big.Int.String`](https://golang.org/pkg/math/big/#Int.String)                      |
| `u.Format`         | [`big.Int.Format`](https://golang.org/pkg/math/big/#Int.Format)                      |
| `u.StoreUint128LE` | [`binary.LittleEndian.PutUint64`](https://golang.org/pkg/encoding/binary/#ByteOrder) |
| `u.LoadUint128LE`  | [`binary.LittleEndian.Uint64`](https://golang.org/pkg/encoding/binary/#ByteOrder)    |
| `u.StoreUint128BE` | [`binary.BigEndian.PutUint64`](https://golang.org/pkg/encoding/binary/#ByteOrder)    |
| `u.LoadUint128BE`  | [`binary.BigEndian.Uint64`](https://golang.org/pkg/encoding/binary/#ByteOrder)       |

See the [documentation][doc] for a complete API specification.


[doc-img]: https://godoc.org/github.com/Pilatuz/uint128?status.svg
[doc]: https://godoc.org/github.com/Pilatuz/uint128
[ci-img]: https://travis-ci.com/Pilatuz/uint128.svg?branch=master
[ci]: https://travis-ci.com/Pilatuz/uint128
[cov-img]: https://codecov.io/gh/Pilatuz/uint128/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/Pilatuz/uint128
[reportcard-img]: https://goreportcard.com/badge/github.com/Pilatuz/uint128
[reportcard]: https://goreportcard.com/report/github.com/Pilatuz/uint128
