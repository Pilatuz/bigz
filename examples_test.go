package uint128_test

import (
	"fmt"
	"math/big"
	"net"

	"github.com/Pilatuz/uint128"
)

// ExampleFromBig is an example for FromBig.
func ExampleFromBig() {
	fmt.Println(uint128.FromBig(nil))
	fmt.Println(uint128.FromBig(new(big.Int).SetInt64(12345)))
	// Output:
	// 0
	// 12345
}

// ExampleFromBigX is an example for FromBigX.
func ExampleFromBigX() {
	one := new(big.Int).SetInt64(1)
	fmt.Println(uint128.FromBigX(new(big.Int).SetInt64(-1))) // => Zero()
	fmt.Println(uint128.FromBigX(one))
	fmt.Println(uint128.FromBigX(one.Lsh(one, 128))) // 2^128, overflows => Max()
	// Output:
	// 0 false
	// 1 true
	// 340282366920938463463374607431768211455 false
}

// ExampleUint128_String is an example for Uint128.String.
func ExampleUint128_String() {
	fmt.Println(uint128.Zero())
	fmt.Println(uint128.One())
	fmt.Println(uint128.Max())
	// Output:
	// 0
	// 1
	// 340282366920938463463374607431768211455
}

// ExampleUint128_Format is an example for Uint128.Format.
func ExampleUint128_Format() {
	fmt.Printf("%08b\n", uint128.From64(42))
	fmt.Printf("%#O\n", uint128.From64(42))
	fmt.Printf("%#x\n", uint128.Max())
	// Output:
	// 00101010
	// 0o52
	// 0xffffffffffffffffffffffffffffffff
}

// ExampleLoadUint128BE is an example for LoadUint128BE.
func ExampleLoadUint128BE() {
	ip := net.ParseIP("cafe::dead:beaf")
	fmt.Printf("%032x\n", uint128.LoadUint128BE(ip.To16()))
	fmt.Printf("%032x\n", uint128.LoadUint128LE(ip.To16()))
	// Output:
	// cafe00000000000000000000deadbeaf
	// afbeadde00000000000000000000feca
}
