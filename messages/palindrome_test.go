package messages

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		{"word: atttta", "atttta", true},
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
	require.True(t, isPalindrome("\u0065\u0301\u0065\u0301"))
}

func TestIsPalindrome_extendedGraphemeClustersAreNotPalindromes(t *testing.T) {
	// see IsPalindrome docs for details.
	//lint:ignore ST1018 copied to display an emoji
	require.False(t, isPalindrome("🤦🏼‍♂️"))
}

func TestIsPalindrome_hiddenCharactersAreNotRemoved(t *testing.T) {
	require.False(t, isPalindrome("mee\u200Bm"))
}
