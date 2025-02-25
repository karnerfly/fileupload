package files

import "io"

type Storage interface {
	Save(path string, body io.Reader) error
}
