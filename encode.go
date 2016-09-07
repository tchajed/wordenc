package wordenc

import (
	"bytes"
	"errors"
	"io"
)

// should encode all 2^12 12-bit sequences (1.5 bytes) and all 2^4 half bytes = 4112 words

type wordEncoder struct {
	out       io.Writer
	halfBytes []byte
	empty     bool
}

func fullAndHalfIndex(full byte, half byte) int {
	return int(full)*(1<<4) + int(half)
}

func halfAndFullIndex(half byte, full byte) int {
	return int(half)*(1<<8) + int(full)
}

func halfIndex(half byte) int {
	return 1<<12 + int(half)
}

func lowerHalf(d byte) byte {
	return d & (1<<4 - 1)
}

func upperHalf(d byte) byte {
	return d >> 4
}

func (we *wordEncoder) writeWord(index int) error {
	if !we.empty {
		if _, err := we.out.Write([]byte{' '}); err != nil {
			return err
		}
	}
	_, err := we.out.Write([]byte(wordList[index]))
	if err != nil {
		return err
	}
	we.empty = false
	return nil
}

func (we *wordEncoder) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if len(we.halfBytes) == 1 {
		err := we.writeWord(halfAndFullIndex(we.halfBytes[0], b[0]))
		if err != nil {
			return 0, err
		}
		b = b[1:]
		we.halfBytes = nil
		n, err := we.Write(b)
		return n + 1, err
	}
	if len(we.halfBytes) == 2 {
		firstByte := we.halfBytes[0]<<4 + we.halfBytes[1]
		index := fullAndHalfIndex(firstByte, upperHalf(b[0]))
		err := we.writeWord(index)
		if err != nil {
			return 0, err
		}
		we.halfBytes = []byte{lowerHalf(b[0])}
		b = b[1:]
		n, err := we.Write(b)
		return n + 1, err
	}
	for bitOffset := 0; bitOffset < len(b)*8; bitOffset += 12 {
		byteOffset := bitOffset / 8
		currByte := b[byteOffset]
		if byteOffset == len(b)-1 && bitOffset%8 == 0 {
			we.halfBytes = []byte{upperHalf(currByte), lowerHalf(currByte)}
		} else if byteOffset == len(b)-1 && bitOffset%8 != 0 {
			we.halfBytes = []byte{lowerHalf(currByte)}
		} else {
			var index int
			if bitOffset%8 == 0 {
				index = fullAndHalfIndex(b[byteOffset], b[byteOffset+1]>>4)
			} else {
				index = halfAndFullIndex(b[byteOffset]&(1<<4-1), b[byteOffset+1])
			}
			err := we.writeWord(index)
			if err != nil {
				return byteOffset, err
			}
		}
	}
	return len(b), nil
}

func (we *wordEncoder) Close() error {
	for _, b := range we.halfBytes {
		if err := we.writeWord(halfIndex(b)); err != nil {
			return err
		}
	}
	return nil
}

// NewEncoder returns a WriteCloser that encodes data written to it and outputs
// the encoded words to w. Partially buffered words are flushed upon calling
// Close on the returned Writecloser.
func NewEncoder(w io.Writer) io.WriteCloser {
	return &wordEncoder{out: w,
		halfBytes: nil,
		empty:     true,
	}
}

// EncodeToString encodes data as space-separated words.
func EncodeToString(data []byte) (string, error) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	n, err := enc.Write(data)
	if n < len(data) {
		return "", errors.New("encoding did not process all bytes")
	}
	if err != nil {
		return "", err
	}
	err = enc.Close()
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
