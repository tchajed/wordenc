package shortdict

import (
	"fmt"
	"testing"
)

func testAddRemove(t *testing.T, b byte) {
	var partial bits
	partial.AddByte(b)
	partialStr := fmt.Sprintf("%v", partial)
	popped := partial.PopWord(8)
	if popped != uint(b) {
		t.Errorf("adding and removing %d results in %d [%s]", b, popped,
			partialStr)
	}
}

func TestAddPopByte(t *testing.T) {
	testAddRemove(t, 1)
	testAddRemove(t, 37)
	testAddRemove(t, 1<<7)
	testAddRemove(t, 1<<8-1)
}

func TestAddPopMultiple(t *testing.T) {
	bytes := []byte{34, 78, 156, 0, 128}
	var partial bits
	for _, b := range bytes {
		partial.AddByte(b)
	}
	for i, b := range bytes {
		popped := byte(partial.PopWord(8))
		if popped != b {
			t.Errorf("popped element %d (%d) results in %d", i,
				b,
				popped)
		}
	}
}
