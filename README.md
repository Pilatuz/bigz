# uint128 [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][reportcard-img]][reportcard]

`uint128` provides a high-performance `Uint128` type that supports standard arithmetic
operations. Unlike `math/big`, operations on `Uint128` values always produce new values
instead of modifying a pointer receiver. A `Uint128` value is therefore immutable, just
like `uint64` and friends.

The name `uint128.Uint128` stutters, so I recommend either using a "dot import"
or aliasing `uint128.Uint128` to give it a project-specific name. Embedding the type
is not recommended, because methods will still return `uint128.Uint128`; this means that,
if you want to extend the type with new methods, your best bet is probably to copy the
source code wholesale and rename the identifier. ¯\\\_(ツ)\_/¯

Released under the [MIT License](LICENSE).

## Installation

```shell
go get github.com/Pilatuz/uint128
```

## Differences

The key differences from [original package](https://github.com/lukechampine/uint128):

- No panics! All methods have wrap-around semantic!
- `Zero` and `Max` are functions to prevent modification of global variables.
- `New` was removed to encourage explicit `Uint128{Lo: ..., Hi: ...}` initialization.
- Trivial (via corresponding `big.Int.Format`) implementation of `Format` method to support for example hex output as `fmt.Sprintf("%X", u)`.
- Store/Load methods in little-endian and big-endian byte order.
- New `Not` and `AndNot` methods.

## Quick Start

TBD

See the [documentation][doc] for a complete API specification.

[doc-img]: https://godoc.org/github.com/Pilatuz/uint128?status.svg
[doc]: https://godoc.org/github.com/Pilatuz/uint128
[ci-img]: https://travis-ci.com/Pilatuz/uint128.svg?branch=master
[ci]: https://travis-ci.com/Pilatuz/uint128
[cov-img]: https://codecov.io/gh/Pilatuz/uint128/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/Pilatuz/uint128
[reportcard-img]: https://goreportcard.com/badge/github.com/Pilatuz/uint128
[reportcard]: https://goreportcard.com/report/github.com/Pilatuz/uint128
