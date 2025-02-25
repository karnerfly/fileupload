package files

import (
	"bytes"
	"os"
	"testing"
)

func TestLocalStorage(t *testing.T) {

	data, err := os.ReadFile("./demo.jpg")
	if err != nil {
		t.Fail()
	}

	store := NewLocalStorage("../test", 1024*1024*5)

	err = store.Save("/d/testing.jpg", bytes.NewReader(data))
	if err != nil {
		t.Fail()
	}
}
