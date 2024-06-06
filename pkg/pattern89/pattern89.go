package pattern89

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type Pattern struct {
	patBytes     []byte
	patWildcards []bool
}

func (p Pattern) Length() int {
	return len(p.patBytes)
}

func ParsePattern(pattern string) (Pattern, error) {
	// Remove commas so that we have just a byte stream.
	pattern = strings.ReplaceAll(pattern, ",", "")
	// Remove whitespace as well.
	pattern = strings.ReplaceAll(pattern, " ", "")
	if len(pattern)%2 != 0 {
		return Pattern{}, fmt.Errorf("invalid pattern")
	}
	patternBytes := make([]byte, len(pattern)/2)
	wildcards := make([]bool, len(pattern)/2)
	for i := 0; i < len(pattern); i += 2 {
		part := pattern[i : i+2]
		if part == "??" {
			patternBytes[i/2] = 0xFF
			wildcards[i/2] = true
		} else {
			b, _ := hex.DecodeString(part)
			patternBytes[i/2] = b[0]
			wildcards[i/2] = false
		}
	}
	return Pattern{patBytes: patternBytes, patWildcards: wildcards}, nil
}

func FindPattern(data []byte, pattern Pattern) int {
	for i := 0; i <= len(data)-len(pattern.patBytes); i++ {
		match := true
		for j := 0; j < len(pattern.patBytes); j++ {
			if !pattern.patWildcards[j] && data[i+j] != pattern.patBytes[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
