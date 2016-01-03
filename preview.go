package preview

import (
	"io"
	"os"
)

const (
	// NotSupportedError is used when a given file is not supported.
	NotSupportedError = previewError("not supported file type.")
)

type previewError string

func (p previewError) Error() string {
	return string(p)
}

// Renderer draw a file to io.Writer.
type Renderer interface {
	Render(io.Writer) error
}

type support struct {
	tags   []string
	create RendererConstructor
}

var supports []support

// RendererConstructor represents function for creating Renderer.
type RendererConstructor func(*os.File, *Config) (Renderer, error)

// Register register Renderer constructor.
func Register(tags []string, c RendererConstructor) {
	supports = append(supports, support{tags, c})
}

// NewRenderer create Renderer.
func NewRenderer(f *os.File, conf *Config) (Renderer, error) {
	_, err := f.Stat()
	if err != nil {
		return nil, err
	}

	var r Renderer
	for _, s := range supports {
		f.Seek(0, 0)
		r, err = s.create(f, conf)
		if r == nil {
			continue
		}
		break
	}

	if r == nil {
		return nil, NotSupportedError
	}

	return r, nil
}
