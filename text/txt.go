package text

import (
	"bufio"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/morikuni/preview"
)

type text struct {
	buf [][]byte
}

// Render is implementation of Renderer.
func (t *text) Render(w io.Writer) error {
	for _, b := range t.buf {
		if _, err := w.Write(b); err != nil {
			return nil
		}
	}
	return nil
}

// NewText is RendererConstructor.
func NewText(f *os.File, conf *preview.Config) (preview.Renderer, error) {
	ext := filepath.Ext(f.Name())

	if ext != ".txt" && ext != ".text" {
		return nil, preview.NotSupportedError
	}

	buf := make([][]byte, conf.Height)

	r := bufio.NewReader(f)

	for i := range buf {
		line, isP, err := r.ReadLine()

		if err != nil {
			if err == io.EOF {
				buf = buf[:i]
				break
			} else {
				return nil, err
			}
		}
		b := make([]byte, int(math.Min(float64(conf.Width), float64(len(line)))))
		copy(b, line)
		buf[i] = append(b, '\n')

		if isP {
			for _, x, _ := r.ReadLine(); x; {
			}
		}
	}

	return &text{buf}, nil
}

func init() {
	preview.Register([]string{"txt", "text"}, NewText)
}
