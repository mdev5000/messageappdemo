// +build propertyTests

package messages

import (
	"github.com/mdev5000/qlik_message/testutil/prop"
	"golang.org/x/text/unicode/norm"
	"testing"
)

const maxLenPropString = 1024

// Generate a palindrome to test against.
func genPal() string {
	for {
		s := prop.GenerateString(maxLenPropString)
		// If you don't convert it, it can take a very long time to find a proper palindrome.
		s = norm.NFKC.String(s)
		sr := []rune(s)
		r := rev(sr)
		out := make([]rune, len(sr)+len(r))
		copy(out, sr)
		copy(out[len(sr):], r)
		check := rev([]rune(norm.NFKC.String(string(out))))
		sout := string(out)
		if sout == string(check) {
			return sout
		}
	}
}

// Generate a non-palindrome to test against.
func genNonPal() string {
	for {
		s := prop.GenerateString(maxLenPropString)
		sr := []rune(s)
		r := rev(sr)
		if s != string(r) {
			return s
		}
	}
}

func rev(sr []rune) []rune {
	r := make([]rune, len(sr))
	for i, c := range sr {
		r[len(r)-1-i] = c
	}
	return r
}

func TestProp_IsPalindrome_correctWhenPalindrome(t *testing.T) {
	for i := 0; i < prop.NumCases(); i++ {
		pal := genPal()
		if !isPalindrome(pal) {
			t.Fatalf("expected palindrome %s returned false", pal)
		}
	}
}

func TestProp_IsPalindrome_correctWhenIsNotPalindrome(t *testing.T) {
	for i := 0; i < prop.NumCases(); i++ {
		pal := genNonPal()
		if isPalindrome(pal) {
			t.Fatalf("expected palindrome %s returned false", pal)
		}
	}
}
