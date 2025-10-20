package schnur

import (
	"testing"
)

func expectEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper() 
	if got != want {
		t.Errorf("%s: got %v; want %v", name, got, want)
	}
}


func TestStringLength(t *testing.T) {
	test := []struct {
		input    string
		expected int
	}{
	{"hello", 5},
	{"", 0},
	{"string to analyze", 17},
	{"12345", 5},
	}
	for _, tc := range test {
		got := stringLength(tc.input)
		expectEqual(t, "StringLength("+tc.input+")", got, tc.expected)
	}	
}

func TestIsPalindrome(t *testing.T) {
	test := []struct {
		input    string
		expected bool
	}{
	{"racecar", true},
	{"hello", false},
	{"A man, a plan, a canal, Panama!", false},
	{"Was it a car or a cat I saw", false},
	{"Never Odd Or Even", false},
	}
	for _, tc := range test {
		got := isPalindrome(tc.input)
		expectEqual(t, "IsPalindrome("+tc.input+")", got, tc.expected)
	}
}

func TestUniqueCharCount(t *testing.T) {
	test := []struct {
		input    string
		expected int
	}{
	{"hello", 4},
	{"", 0},
	{"abcABC", 3},
	{"Never Odd Or Even", 6},
	}
	for _, tc := range test {
		got := uniqueCharCount(tc.input)
		expectEqual(t, "UniqueCharCount("+tc.input+")", got, tc.expected)
	}
}

func TestWordCount(t *testing.T) {
	test := []struct {
		input    string
		expected int
	}{
	{"hello world", 2},
	{"", 0},
	{"one two three four five", 5},
	{"   leading and trailing spaces   ", 4},
	{"string to analyze", 3},
	}
	for _, tc := range test {
		got := wordCount(tc.input)
		expectEqual(t, "WordCount("+tc.input+")", got, tc.expected)
	}
}

func TestCharFrequency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[rune]int
	}{
		{"normal string", "hello", map[rune]int{'h': 1, 'e': 1, 'l': 2, 'o': 1}},
		{"empty string", "", map[rune]int{}},
		{"repeated letters", "aabbcc", map[rune]int{'a': 2, 'b': 2, 'c': 2}},
	}

	for _, tc := range tests {
		got := charFrequency(tc.input)

		if len(got) != len(tc.expected) {
			t.Errorf("%s: expected length %d, got %d", tc.name, len(tc.expected), len(got))
			continue
		}
		
		for k, v := range tc.expected {
			if got[k] != v {
				t.Errorf("%s: for key %q expected %d, got %d", tc.name, k, v, got[k])
			}
		}
	}
}

func TestHashString(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
	{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	{"string to analyze", "94b4087035c47dc5ec70499327758a792a6a4db132313a67143ec61dc489c33f"},
	}
	for _, tc := range test {
		got := hashString(tc.input)
		expectEqual(t, "HashString("+tc.input+")", got, tc.expected)
	}	
}