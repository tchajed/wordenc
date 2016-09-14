package shortdict

// Package shortdict implements word encoding for small power-of-2-sized //
// dictionaries.

import "strings"

const wordBits = 11

type bits struct {
	data []byte
}

func (b *bits) AddData(d uint, length uint) {
	for i := 0; i < int(length); i++ {
		highBit := d >> (length - 1)
		b.data = append(b.data, byte(highBit&1))
		d <<= 1
	}
}

func (b *bits) AddByte(d byte) {
	b.AddData(uint(d), 8)
}

// Length gives the number of bits in b
func (b *bits) Length() uint {
	return uint(len(b.data))
}

func (b *bits) PopWord(length uint) (d uint) {
	data := b.data[:length]
	for i := uint(0); i < length; i++ {
		offset := int(length - i - 1)
		d += uint(data[offset]) * (1 << i)
	}
	b.data = b.data[length:]
	return
}

func (b *bits) PadTo(length uint) {
	for len(b.data) < int(length) {
		b.data = append(b.data, 0)
	}
}

func getWord(index int) string {
	if index > len(wordList) {
		return "out-of-range"
	}
	return wordList[index][0]
}

// EncodeToString encodes data into a string of words separated by spaces.
func EncodeToString(data []byte) string {
	var words []string
	var partial bits
	for _, d := range data {
		partial.AddByte(d)
		if partial.Length() >= wordBits {
			index := int(partial.PopWord(wordBits))
			words = append(words, getWord(index))
		}
	}
	if partial.Length() > 0 {
		partial.PadTo(wordBits)
		index := int(partial.PopWord(wordBits))
		words = append(words, getWord(index))
	}
	return strings.Join(words, " ")
}
