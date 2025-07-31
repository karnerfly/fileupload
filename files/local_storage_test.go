package files

import (
	"bytes"
	"os"
	"testing"
)

func TestLocalStorage(t *testing.T) {

	data, err := os.ReadFile("../tests/car.jpg")
	if err != nil {
		t.Fail()
	}

	store := NewLocalStorage("../tests/store", 1024*1024*5)

	err = store.Save("upload/testing.jpg", bytes.NewReader(data))
	if err != nil {
		t.Fail()
	}
}

func TestFilePath(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{input: "test", output: "/home/karnerfly/codes/web_projects/fileupload/files/base/test"},
		{input: "./test", output: "/home/karnerfly/codes/web_projects/fileupload/files/base/test"},
		{input: "test/test2", output: "/home/karnerfly/codes/web_projects/fileupload/files/base/test/test2"},
		{input: "../test", output: "/home/karnerfly/codes/web_projects/fileupload/files/test"},
		{input: "../test/test2", output: "/home/karnerfly/codes/web_projects/fileupload/files/test/test2"},
		{input: "/", output: "/home/karnerfly/codes/web_projects/fileupload/files/base"},
		{input: ".", output: "/home/karnerfly/codes/web_projects/fileupload/files/base"},
		{input: "./", output: "/home/karnerfly/codes/web_projects/fileupload/files/base"},
		{input: "", output: "/home/karnerfly/codes/web_projects/fileupload/files/base"},
		{input: "..", output: "/home/karnerfly/codes/web_projects/fileupload/files"},
		{input: "../", output: "/home/karnerfly/codes/web_projects/fileupload/files"},
	}

	store := NewLocalStorage("base", 1024)

	for _, v := range tests {
		get, err := store.getFullPath(v.input)
		if err != nil {
			t.Error(err)
		}
		if v.output != get {
			t.Errorf("[Test Failed]: Want %s but get %s", v.output, get)
		}
	}
}
