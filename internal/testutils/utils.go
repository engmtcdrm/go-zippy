package testutils

import (
	"math/rand"
	"os"
	"strings"
)

// prefixAndSuffix splits pattern by the last wildcard "*", if applicable,
// returning prefix as the part before "*" and suffix as the part after "*".
// Lovingly stolen and modified from [os] package.
func prefixAndSuffix(pattern string) (prefix, suffix string) {
	for i := range pattern {
		if os.IsPathSeparator(pattern[i]) {
			return "", ""
		}
	}

	if pos := strings.LastIndexByte(pattern, '*'); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}

	return prefix, suffix
}

// genRandomDigits generates a random digit string.
func genRandomDigits(length int) string {
	if length < 1 {
		length = 10
	}

	digits := make([]byte, length)
	for i := range digits {
		digits[i] = byte(rand.Int63()%10 + '0') // Generate a random number between '0' and '9'
	}
	return string(digits)
}
