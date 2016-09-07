package wordenc

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func encode(data []byte, t *testing.T) string {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	n, err := enc.Write(data)
	if n < len(data) {
		t.Errorf("encoding did not write fully: %d/%d bytes", n, len(data))
		return ""
	}
	if err != nil {
		panic("writing to buffer failed")
	}
	err = enc.Close()
	if err != nil {
		panic("closing buffer failed")
	}
	return b.String()
}

func decode(s string) ([]byte, error) {
	data := []byte(s)
	dec := NewDecoder(bytes.NewReader(data))
	return ioutil.ReadAll(dec)
}

func roundtrip(t *testing.T, data []byte) {
	encoded := encode(data, t)
	decoded, err := decode(encoded)
	if err != nil {
		t.Errorf("decoding %v gives error %v", data, err)
		return
	}
	incorrect := false
	if len(data) != len(decoded) {
		incorrect = true
		t.Errorf("encoded/decoded %v has length %d, expected %d", data, len(decoded), len(data))
	} else {
		for i := 0; i < len(data); i++ {
			if data[i] != decoded[i] {
				incorrect = true
			}
		}
	}
	if incorrect {
		t.Errorf("encoded/decoded %v results in %v (%s)",
			data,
			decoded,
			encoded)
	}
}

func TestEncodeDecodeEmpty(t *testing.T) {
	roundtrip(t, []byte{})
}

func TestEncodeDecodeSingleByte(t *testing.T) {
	roundtrip(t, []byte{3})
	roundtrip(t, []byte{255})
	roundtrip(t, []byte{0})
}

func TestEncodeDecodeEvenBytes(t *testing.T) {
	roundtrip(t, []byte{2, 3})
	roundtrip(t, []byte{2, 4})
	roundtrip(t, []byte{123, 104, 12, 128})
	roundtrip(t, []byte{123, 104, 12, 86})
	roundtrip(t, []byte{123, 104, 12, 86, 100, 0})
}

func TestEncodeDecodeOddBytes(t *testing.T) {
	roundtrip(t, []byte{2, 3, 0})
	roundtrip(t, []byte{0, 2, 3})
	roundtrip(t, []byte{123, 4, 104, 12, 86})
	roundtrip(t, []byte{123, 104, 12, 255, 86, 100, 0})
}

func TestEncodeDecodeMultipleOfThree(t *testing.T) {
	roundtrip(t, []byte{2, 3, 0})
	roundtrip(t, []byte{0, 2, 43})
	roundtrip(t, []byte{32, 107, 65, 12, 204, 198})
}
