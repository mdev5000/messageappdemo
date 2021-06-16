package messages

import "golang.org/x/text/unicode/norm"

func isPalindrome(msg string) bool {
	r := []rune(norm.NFC.String(msg))
	for start, end := 0, len(r)-1; start < end; start, end = start+1, end-1 {
		if r[start] != r[end] {
			return false
		}
	}
	return true
}
