package text

import (
	"bufio"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/morikuni/preview"
)

// PreviewTxt print text file.
func PreviewTxt(path string, out io.Writer, conf *preview.Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	ext := filepath.Ext(f.Name())

	if ext != ".txt" && ext != ".text" {
		return preview.NotSupportedError
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
				return err
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

	for _, b := range buf {
		if _, err := out.Write(b); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	preview.Register([]string{"txt"}, PreviewTxt)
}
