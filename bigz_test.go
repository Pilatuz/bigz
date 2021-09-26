package bigz_test

import (
	"fmt"

	"github.com/Pilatuz/bigz"
)

// Example_new an example of creating Uint128 values.
func Example_new() {
	fmt.Println(bigz.Uint128{Lo: 12345})
	// Output:
	// 12345
}
