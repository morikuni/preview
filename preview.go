package preview

import (
	"io"
	"os"
)

const (
	// UnsupportedError is used when a given file is not supported.
	UnsupportedError = previewError("not supported file type.")
)

type previewError string

func (p previewError) Error() string {
	return string(p)
}

type previewFunc func(path string, out io.Writer, conf *Config) error

// Renderer draw a file to io.Writer.
type Renderer interface {
	Render(io.Writer) error
}

type support struct {
	tags []string
	pf   previewFunc
}

var supports []support

// Register register Renderer constructor.
func Register(tags []string, pf previewFunc) {
	supports = append(supports, support{tags, pf})
}

// Preview print a file to out with conf.
func Preview(path string, out io.Writer, conf *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, s := range supports {
		result := s.pf(path, out, conf)
		if result == UnsupportedError {
			continue
		}
		return result
	}

	return UnsupportedError
}
