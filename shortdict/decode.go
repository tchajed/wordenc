package shortdict

import (
	"fmt"
	"strings"
)

var wordMapping = make(map[string]int)

func init() {
	for i, words := range wordList {
		for _, word := range words {
			wordMapping[word] = i
		}
	}
}

func lookupWord(w string) (int, error) {
	index, ok := wordMapping[w]
	if !ok {
		return 0, fmt.Errorf("invalid word %s", w)
	}
	return index, nil
}

// decodeWords decodes the sequence words to a byte array of length bytes.
func decodeWords(words []string, length int) (p []byte, err error) {
	// remove unneccessary words; this is an optimization - no more than length
	// bytes are ever added to the output
	numWordsNeeded := length * 8 / wordBits
	if numWordsNeeded*wordBits < length*8 {
		numWordsNeeded++
	}
	if numWordsNeeded < len(words) {
		words = words[:numWordsNeeded]
	}
	var partial bits
	for _, w := range words {
		index, err := lookupWord(w)
		if err != nil {
			return nil, err
		}
		partial.AddData(uint(index), wordBits)
		for partial.Length() >= 8 && len(p) < length {
			p = append(p, byte(partial.PopWord(8)))
		}
	}
	return
}

// DecodeString decodes s, assumed to be length bytes long, into a byte array.
func DecodeString(s string, length int) ([]byte, error) {
	return decodeWords(strings.Fields(s), length)
}
