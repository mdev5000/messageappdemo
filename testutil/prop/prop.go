// +build propertyTests

// Contains help functions for property based testing.
package prop

import (
	"fmt"
	"os"
	"strconv"

	"github.com/chrismcguire/gobberish"
)

const defaultNumCases = 2000

// NumCases determines how many cases should be run for property tests. This can be configured via the NUM_PROP_TEST
// environment variable.
func NumCases() int {
	numCasesS := os.Getenv("NUM_PROP_TESTS")
	if numCasesS == "" {
		return defaultNumCases
	}
	numCases, err := strconv.Atoi(numCasesS)
	if err != nil {
		fmt.Printf("Invalid NUM_PROP_TESTS value, must be a number, was %s.\n", numCasesS)
		return defaultNumCases
	}
	return numCases
}

// Generate a random UTF-8 string.
func GenerateString(n int) string {
	return gobberish.GenerateString(n)
}
