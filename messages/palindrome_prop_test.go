// +build propertyTests

package messages

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"testing"

	"github.com/mdev5000/messageappdemo/testutil/prop"
	"golang.org/x/text/unicode/norm"
)

const maxLenPropString = 1024

func init() {
	var b [8]byte
	_, err := cryptorand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

// Generate a palindrome to test against.
func genPal() string {
	for {
		s := prop.GenerateString(maxLenPropString)
		// If you don't convert it, it can take a very long time to find a proper palindrome.
		s = norm.NFC.String(s)
		sr := []rune(s)
		r := rev(sr)
		out := make([]rune, len(sr)+len(r))
		copy(out, sr)
		copy(out[len(sr):], r)
		check := rev([]rune(norm.NFC.String(string(out))))
		sout := string(out)
		// NFC conversion does not always result in the same value when the string is reversed so make sure this
		// string is actually a palindrome.
		if sout == string(check) {
			// Randomly choose either NFC or NFD format.
			if rand.Intn(2) == 1 {
				return norm.NFD.String(sout)
			}
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
