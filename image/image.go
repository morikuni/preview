package image

import (
	"fmt"
	gimage "image"
	"math"
	// for support
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/morikuni/preview"
)

// https://github.com/itchyny/cam/blob/master/cam.c
func color3bit(r, g, b uint8) string {
	c := 40 + (r >> 7) | ((g >> 6) & 0x2) | ((b >> 5) & 0x4)
	return fmt.Sprintf("\x1b[%dm", c)
}

func color8bit(r, g, b uint8) string {
	c := 16 + 36*(r/43) + 6*(g/43) + b/43
	return fmt.Sprintf("\x1b[48;5;%dm", c)
}

type colorFunc func(r, g, b uint8) string

const reset = "\x1b[0m"

func asUint8(r, g, b uint32) (uint8, uint8, uint8) {
	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
}

type image struct {
	img       gimage.Image
	width     uint
	height    uint
	colorFunc colorFunc
}

// Render is implementation of Renderer
func (i *image) Render(w io.Writer) error {
	bounds := i.img.Bounds()
	minX := float64(bounds.Min.X)
	maxX := float64(bounds.Max.X)
	minY := float64(bounds.Min.Y)
	maxY := float64(bounds.Max.Y)
	iw := maxX - minX
	ih := maxY - minY
	cw := float64(i.width)
	ch := float64(i.height)
	var width, height float64
	if iw < cw {
		width = iw
		height = ih
	} else {
		width = cw
		height = ih * cw / iw
	}
	if ch < height {
		width = iw * ch / ih
		height = ch
	}

	dx := (maxX - minX) / float64(width)
	dy := (maxY - minY) / float64(height)
	for wy := float64(0); wy < height; wy++ {
		for wx := float64(0); wx < width; wx++ {
			lx, ly := minX+wx*dx, minY+(wy)*dy
			hx, hy := math.Min(minX+(wx+1)*dx, maxX), math.Min(minY+(wy+1)*dy, maxY)
			r, g, b := i.average(int(lx), int(ly), int(hx), int(hy))
			if _, err := fmt.Fprintf(w, "%s  ", i.colorFunc(r, g, b)); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, reset); err != nil {
			return err
		}
	}
	return nil
}

func (i *image) average(lx, ly, hx, hy int) (uint8, uint8, uint8) {
	var r, g, b float64
	var n uint32
	for y := ly; y < hy; y++ {
		for x := lx; x < hx; x++ {
			r16, g16, b16, _ := i.img.At(x, y).RGBA()
			r += float64(r16)
			g += float64(g16)
			b += float64(b16)
			n++
		}
	}
	r8, g8, b8 := asUint8(uint32(r)/n, uint32(g)/n, uint32(b)/n)
	return r8, g8, b8
}

// NewImage is RendererConstructor.
func NewImage(f *os.File, conf *preview.Config) (preview.Renderer, error) {
	img, _, err := gimage.Decode(f)
	if err != nil {
		return nil, err
	}

	var cf colorFunc
	switch conf.Color {
	case preview.Color8:
		cf = color3bit
	case preview.Color256:
		cf = color8bit
	}

	i := &image{img, conf.Width / 2, conf.Height, cf}

	return i, nil
}

func init() {
	preview.Register([]string{"jpg", "png", "gif", "tiff", "bmp"}, NewImage)
}
