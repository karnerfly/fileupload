package files

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	maxSize  int64
	basePath string
}

func NewLocalStorage(b string, m int64) *LocalStorage {
	return &LocalStorage{
		basePath: b,
		maxSize:  m,
	}
}

func (ls *LocalStorage) Save(path string, body io.Reader) error {
	fp := ls.getFullPath(path)

	buf := &bytes.Buffer{}
	n, err := io.Copy(buf, body)
	if err != nil {
		return fmt.Errorf("cannot write data into buffer for file: %s", fp)
	}

	if n > ls.maxSize {
		return fmt.Errorf("too large file size: %d", n)
	}

	dir := filepath.Dir(fp)

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create directory: %s", dir)
	}

	if _, err = os.Stat(fp); err == nil {
		if err = os.Remove(fp); err != nil {
			return fmt.Errorf("cannot delete the file: %s", fp)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot get file info: %s", err.Error())
	}

	file, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("cannot create file: %s", fp)
	}

	defer file.Close()

	_, err = io.Copy(file, buf)
	if err != nil {
		return fmt.Errorf("cannot write into file: %s", fp)
	}

	return nil
}

func (ls *LocalStorage) getFullPath(path string) string {
	return filepath.Join(ls.basePath, path)
}
