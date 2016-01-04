package text

import (
	"bufio"
	"fmt"
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
		return preview.UnsupportedError
	}

	buf := make([][]byte, conf.Height)

	sc := bufio.NewScanner(f)

	for i := uint(0); i < conf.Height && sc.Scan(); i++ {
		line := sc.Text()

		if err != nil {
			if err == io.EOF {
				buf = buf[:i]
				break
			} else {
				return err
			}
		}
		b := line[:int(math.Min(float64(conf.Width), float64(len(line))))]

		if _, err := fmt.Fprintln(out, string(b)); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	preview.Register([]string{"txt"}, PreviewTxt)
}
