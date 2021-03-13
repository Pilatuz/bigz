package bigx_test

import (
	"fmt"

	"github.com/Pilatuz/bigx"
)

// Example_new an example of creating Uint128 values.
func Example_new() {
	fmt.Println(bigx.Uint128{Lo: 12345})
	// Output:
	// 12345
}
