package handlers

import (
	"os"
	"testing"

	"github.com/Pureparadise56b/fileupload/files"
)

func TestSaveFile(t *testing.T) {
	store := files.NewLocalStorage("../test", 1024*1024*5)
	fh := NewFileHandler(store)

	file, err := os.Open("./demo.jpg")
	if err != nil {
		t.Fail()
	}

	err = fh.saveFile("123", "hello.jpg", file)
	if err != nil {
		t.Fail()
	}
}
