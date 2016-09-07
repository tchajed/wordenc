package wordenc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

var wordMapping = make(map[string]int)

func init() {
	for i, w := range wordList {
		wordMapping[w] = i
	}
}

// wordStream wraps a word-splitting scanner, converting to word indices.
type wordStream struct {
	scanner *bufio.Scanner
}

func newWordStream(r io.Reader) wordStream {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	return wordStream{scanner}
}

func (s wordStream) NextWord() (index int, err error) {
	if !s.scanner.Scan() {
		return 0, io.EOF
	}
	word := s.scanner.Text()
	index, ok := wordMapping[word]
	if !ok {
		return 0, fmt.Errorf("invalid word %s found", word)
	}
	return
}

// Considers the concatenation half || w and splits it into two bytes.
func joinHalfWith12Bits(half byte, w int) (byte, byte) {
	firstByte := half<<4 + byte(w>>8&(1<<8-1))
	secondByte := byte(w & (1<<8 - 1))
	return firstByte, secondByte
}

// Returns the first byte in the 12-bit w and the leftover half byte.
func split12Bits(w int) (full byte, half byte) {
	full = byte(w >> 4 & (1<<8 - 1))
	half = byte(w & (1<<4 - 1))
	return
}

type wordDecoder struct {
	ws        wordStream
	halfBytes []byte
}

func (wd *wordDecoder) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(wd.halfBytes) >= 2 {
		p[0] = wd.halfBytes[0]<<4 + wd.halfBytes[1]
		wd.halfBytes = wd.halfBytes[2:]
		return 1, nil
	}
	w, err := wd.ws.NextWord()
	if err != nil {
		return 0, err
	}
	if w < 1<<12 {
		if len(wd.halfBytes) == 1 {
			b1, b2 := joinHalfWith12Bits(wd.halfBytes[0], w)
			if len(p) >= 2 {
				p[0] = b1
				p[1] = b2
				wd.halfBytes = nil
				return 2, nil
			}
			p[0] = b1
			wd.halfBytes = []byte{upperHalf(b2), lowerHalf(b2)}
			return 1, nil
		}
		b, h := split12Bits(w)
		p[0] = b
		wd.halfBytes = []byte{h}
		return 1, nil
	}

	// add single half byte
	h := byte((w - 1<<12) & (1<<8 - 1))
	wd.halfBytes = append(wd.halfBytes, h)
	return wd.Read(p)
}

// NewDecoder constructs a new word stream decoder.
func NewDecoder(r io.Reader) io.Reader {
	s := newWordStream(r)
	return &wordDecoder{s, nil}
}

// DecodeFromString decodes whitespace-separated word-encoded data to bytes.
func DecodeFromString(s string) ([]byte, error) {
	data := []byte(s)
	dec := NewDecoder(bytes.NewReader(data))
	return ioutil.ReadAll(dec)
}
