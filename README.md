# bigz [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][reportcard-img]][reportcard]

`bigz/uint128` provides a high-performance `Uint128` type that supports standard arithmetic
operations. Unlike `math/big`, operations on `Uint128` always produce new values
instead of modifying a pointer receiver. A `Uint128` value is therefore immutable, just
like `uint64` and friends.

`bigz/uint256` provides similar `Uint256` type (note, 256-bit division is still subject for some optimizations).

Released under the [MIT License](LICENSE).


## Installation

```shell
go get github.com/Pilatuz/bigz
```

The name `uint128.Uint128` and `uint256.Uint256` stutter, so it is recommended either using a "facade" package:

```go
import (
    "github.com/Pilatuz/bigz"
)

// then use bigz.Uint128 type
// then use bigz.Uint256 type
```

or type aliasing to give it a project-specific name:

```go
import (
    "github.com/Pilatuz/bigz/uint128"
    "github.com/Pilatuz/bigz/uint256"
)

type U128 = uint128.Uint128
type U256 = uint256.Uint256
```


## What's new

The key differences from [original package](https://github.com/lukechampine/uint128):

- No panics! All methods have wrap-around semantic!
- `Zero` and `Max` are functions to prevent modification of global variables.
- `New` was removed to encourage explicit `Uint128{Lo: ..., Hi: ...}` initialization.
- Trivial (via `big.Int`) implementation of fmt.Formatter interface to support for example hex output as `fmt.Sprintf("%X", u)`.
- Trivial (via `big.Int`) implementation of TextMarshaller and TextUnmarshaler interfaces to support JSON encoding.
- Store/Load methods support little-endian and big-endian byte order.
- New `Not` and `AndNot` methods.
- New `uint256.Uint256` type.


## Quick Start

The 128-bit or 256-bit integer can be initialized in the following ways:

| `uint128` 128-bit package          | `uint256` 256-bit package            | Description                                       |
|------------------------------------|--------------------------------------|---------------------------------------------------|
| `u := Uint128{Lo: lo64, Hi: hi64}` | `u := Uint256{Lo: lo128, Hi: hi128}` | Set both lower half and upper half.               |
| `u := From64(lo64)`                | `u := From128(lo128)`                | Set only lower half.                              |
|                                    | `u := From64(lo64)`                  | Set only lower 64-bit.                            |
| `u := Zero()`                      | `u := Zero()`                        | The same as `From64(0)`.                          |
| `u := One()`                       | `u := One()`                         | The same as `From64(1)`.                          |
| `u := Max()`                       | `u := Max()`                         | The largest possible value.                       |
| `u := FromBig(big)`                | `u := FromBig(big)`                  | Convert from `*big.Int` with saturation.          |
| `u := FromBigX(big)`               | `u := FromBigX(big)`                 | The same as `FromBig` but provides `ok` flag.     |

The following arithmetic operations are supported:

| `bigz.Uint128`           | `bigz.Uint256`            | Standard `*big.Int` equivalent                                  |
|--------------------------|---------------------------|-----------------------------------------------------------------|
| `u.Add`, `u.Add64`       | `u.Add`, `u.Add128`       | [`big.Int.Add`](https://golang.org/pkg/math/big/#Int.Add)       |
| `u.Sub`, `u.Sub64`       | `u.Sub`, `u.Sub128`       | [`big.Int.Sub`](https://golang.org/pkg/math/big/#Int.Sub)       |
| `u.Mul`, `u.Mul64`       | `u.Mul`, `u.Mul128`       | [`big.Int.Mul`](https://golang.org/pkg/math/big/#Int.Mul)       |
| `u.Div`, `u.Div64`       | `u.Div`, `u.Div128`       | [`big.Int.Div`](https://golang.org/pkg/math/big/#Int.Div)       |
| `u.Mod`, `u.Mod64`       | `u.Mod`, `u.Mod128`       | [`big.Int.Mod`](https://golang.org/pkg/math/big/#Int.Mod)       |
| `u.QuoRem`, `u.QuoRem64` | `u.QuoRem`, `u.QuoRem128` | [`big.Int.QuoRem`](https://golang.org/pkg/math/big/#Int.QuoRem) |

The following logical and comparison operations are supported:

| `bigz.Uint128`           | `bigz.Uint1256`           | Standard `*big.Int` equivalent                                  |
|--------------------------|---------------------------|-----------------------------------------------------------------|
| `u.Equals`, `u.Equals64` | `u.Equals`, `u.Equals128` | [`big.Int.Cmp == 0`](https://golang.org/pkg/math/big/#Int.Cmp)  |
| `u.Cmp`, `u.Cmp64`       | `u.Cmp`     `u.Cmp64`     | [`big.Int.Cmp`](https://golang.org/pkg/math/big/#Int.Cmp)       |
| `u.Not`                  | `u.Not`                   | [`big.Int.Not`](https://golang.org/pkg/math/big/#Int.Not)       |
| `u.AndNot`, `u.AndNot64` | `u.AndNot`, `u.AndNot128` | [`big.Int.AndNot`](https://golang.org/pkg/math/big/#Int.AndNot) |
| `u.And`, `u.And64`       | `u.And`, `u.And128`       | [`big.Int.And`](https://golang.org/pkg/math/big/#Int.And)       |
| `u.Or`, `u.Or64`         | `u.Or`, `u.Or128`         | [`big.Int.Or`](https://golang.org/pkg/math/big/#Int.Or)         |
| `u.Xor`, `u.Xor64`       | `u.Xor`, `u.Xor128`       | [`big.Int.Xor`](https://golang.org/pkg/math/big/#Int.Xor)       |
| `u.Lsh`                  | `u.Lsh`                   | [`big.Int.Lsh`](https://golang.org/pkg/math/big/#Int.Lsh)       |
| `u.Rsh`                  | `u.Rsh`                   | [`big.Int.Rsh`](https://golang.org/pkg/math/big/#Int.Rsh)       |

The following bit operations are supported:

| `bigz.Uint128`    | `bigz.Uint256`    | Standard 64-bit equivalent                                                  |
|-------------------|-------------------|-----------------------------------------------------------------------------|
| `u.RotateLeft`    | `u.RotateLeft`    | [`bits.RotateLeft64`](https://golang.org/pkg/math/bits/#RotateLeft64)       |
| `u.RotateRight`   | `u.RotateRight`   | [`bits.RotateRight64`](https://golang.org/pkg/math/bits/#RotateRight64)     |
| `u.BitLen`        | `u.BitLen`        | [`bits.Len64`](https://golang.org/pkg/math/bits/#Len64) or [`big.Int.BitLen`](https://golang.org/pkg/math/big/#Int.BitLen) |
| `u.LeadingZeros`  | `u.LeadingZeros`  | [`bits.LeadingZeros64`](https://golang.org/pkg/math/bits/#LeadingZeros64)   |
| `u.TrailingZeros` | `u.TrailingZeros` | [`bits.TrailingZeros64`](https://golang.org/pkg/math/bits/#TrailingZeros64) |
| `u.OnesCount`     | `u.OnesCount`     | [`bits.OnesCount64`](https://golang.org/pkg/math/bits/#OnesCount64)         |
| `u.Reverse`       | `u.Reverse`       | [`bits.Reverse64`](https://golang.org/pkg/math/bits/#Reverse64)             |
| `u.ReverseBytes`  | `u.ReverseBytes`  | [`bits.ReverseBytes64`](https://golang.org/pkg/math/bits/#ReverseBytes64)   |

The following miscellaneous operations are supported:

| `bigz.Uint128`      | `bigz.Uint256`      | Standard equivalent                                                                  |
|---------------------|---------------------|--------------------------------------------------------------------------------------|
| `u.String`          | `u.String`          | [`big.Int.String`](https://golang.org/pkg/math/big/#Int.String)                      |
| `u.Format`          | `u.Format`          | [`big.Int.Format`](https://golang.org/pkg/math/big/#Int.Format)                      |
| `u.MarshalText`     | `u.MarshalText`     | [`big.Int.MarshalText`](https://golang.org/pkg/math/big/#Int.MarshalText)            |
| `u.UnmarshalText`   | `u.UnmarshalText`   | [`big.Int.UnmarshalText`](https://golang.org/pkg/math/big/#Int.UnmarshalText)        |
| `StoreLittleEndian` | `StoreLittleEndian` | [`binary.LittleEndian.PutUint64`](https://golang.org/pkg/encoding/binary/#ByteOrder) |
| `LoadLittleEndian`  | `LoadLittleEndian`  | [`binary.LittleEndian.Uint64`](https://golang.org/pkg/encoding/binary/#ByteOrder)    |
| `StoreBigEndian`    | `StoreBigEndian`    | [`binary.BigEndian.PutUint64`](https://golang.org/pkg/encoding/binary/#ByteOrder)    |
| `LoadBigEndian`     | `LoadBigEndian`     | [`binary.BigEndian.Uint64`](https://golang.org/pkg/encoding/binary/#ByteOrder)       |

See the [documentation][doc] for a complete API specification.


[doc-img]: https://godoc.org/github.com/Pilatuz/bigz?status.svg
[doc]: https://godoc.org/github.com/Pilatuz/bigz
[ci-img]: https://travis-ci.com/Pilatuz/bigz.svg?branch=master
[ci]: https://travis-ci.com/Pilatuz/bigz
[cov-img]: https://codecov.io/gh/Pilatuz/bigz/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/Pilatuz/bigz
[reportcard-img]: https://goreportcard.com/badge/github.com/Pilatuz/bigz
[reportcard]: https://goreportcard.com/report/github.com/Pilatuz/bigz
