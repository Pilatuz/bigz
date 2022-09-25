package uint256_test

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Pilatuz/bigz/uint256"
)

// ExampleFromBig is an example for FromBig.
func ExampleFromBig() {
	fmt.Println(uint256.FromBig(nil))
	fmt.Println(uint256.FromBig(new(big.Int).SetInt64(12345)))
	// Output:
	// 0
	// 12345
}

// ExampleFromBigEx is an example for FromBigEx.
func ExampleFromBigEx() {
	one := new(big.Int).SetInt64(1)
	fmt.Println(uint256.FromBigEx(new(big.Int).SetInt64(-1))) // => Zero()
	fmt.Println(uint256.FromBigEx(one))
	fmt.Println(uint256.FromBigEx(one.Lsh(one, 256))) // 2^256, overflows => Max()
	// Output:
	// 0 false
	// 1 true
	// 115792089237316195423570985008687907853269984665640564039457584007913129639935 false
}

// ExampleFromString is an example for FromString.
func ExampleFromString() {
	u, _ := uint256.FromString("1")
	fmt.Println(u)
	_, err := uint256.FromString("-1")
	fmt.Println(err)
	// Output:
	// 1
	// out of 256-bit range
}

// ExampleUint256_String is an example for Uint256.String.
func ExampleUint256_String() {
	fmt.Println(uint256.Zero())
	fmt.Println(uint256.One())
	fmt.Println(uint256.Max())
	// Output:
	// 0
	// 1
	// 115792089237316195423570985008687907853269984665640564039457584007913129639935
}

// ExampleUint256_Format is an example for Uint256.Format.
func ExampleUint256_Format() {
	fmt.Printf("%08b\n", uint256.From64(42))
	fmt.Printf("%#O\n", uint256.From64(42))
	fmt.Printf("%#x\n", uint256.Max())
	// Output:
	// 00101010
	// 0o52
	// 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
}

// ExampleUint256_json is an example for JSON marshaling.
func ExampleUint256_json() {
	foo := map[string]interface{}{
		"bar": uint256.From64(12345),
	}

	buf, _ := json.Marshal(foo)
	fmt.Printf("%s", buf)
	// Output:
	// {"bar":"12345"}
}
