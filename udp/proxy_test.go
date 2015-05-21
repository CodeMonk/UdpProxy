package udp

import (
	"bytes"

	"testing"
)

func TestReverse(t *testing.T) {
	one := []byte("One")
	reversed := reverse(one)

	if bytes.Compare(reversed, []byte("enO")) != 0 {
		t.Errorf("Reverse did not work. %s -> %s\n", one, reversed)
	}

	back := reverse(reversed)
	if bytes.Compare(back, one) != 0 {
		t.Errorf("Reverse did not work. %s -> %s\n", reversed, back)
	}
}
