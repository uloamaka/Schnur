package schnur

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"
)

// a function that returns the length of the input string
func stringLength(s string) int {
	return len(s)
}

// a function that checks if string is palindrome
func isPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return false
		}
	}
	return true
}

// a function that returns the count of the unique char
func uniqueCharCount(s string) int {
	charMap := make(map[rune]bool) 
	cnt := 0
	trimedS := trimSpace(s)
	lwcs := strings.ToLower(trimedS)
	for _, char := range lwcs {
		if !charMap[char] {
			charMap[char] = true
			cnt++
		}
	}
	return cnt
}

// a function that returns the count of word separated by whitespaces
func wordCount(s string) int {
	words := strings.Fields(s)
	return len(words)
}

// a function that returns a freq count of the char in the input string
func charFrequency(s string) map[string]int {
	freqMap := make(map[string]int)

	trimedS := trimSpace(s)
	if len(trimedS) == 0 {
		return freqMap
	}

	lwcs := strings.ToLower(trimedS)

	for _, char := range lwcs {
		freqMap[string(char)]++
	}
	return freqMap
}

// a function that performs SHA256 hashing on the input string
func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func trimSpace(str string) string {
    var b strings.Builder
    b.Grow(len(str))
    for _, ch := range str {
        if !unicode.IsSpace(ch) {
            b.WriteRune(ch)
        }
    }
    return b.String()
}