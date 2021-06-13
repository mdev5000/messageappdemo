package messages

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected bool
	}{
		{"empty string", "", true},
		{"single letter", "a", true},
		{"double letter (same)", "bb", true},
		{"double letter async", "bd", false},
		{"triple letter letter (same)", "bbb", true},
		{"word: cabac", "cabac", true},
		{"word: ducks", "ducks", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.expected, isPalindrome(c.value))
		})
	}
}

func TestIsPalindrome_lettersFollowedByCombiningCharactersAreStillPalindromes(t *testing.T) {
	// u0301 is a combining character which adds an accent. So this string is equivalent to "éé".
	// However is would need to be normalized to NFC for it to be recognized as a palindrome.
	s := "\u0065\u0301\u0065\u0301"
	require.True(t, isPalindrome(s))
}
