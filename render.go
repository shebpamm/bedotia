// See https://sw.kovidgoyal.net/kitty/graphics-protocol.html.
package main

import (
	"fmt"
	"github.com/mazznoer/colorgrad"
	"golang.org/x/sys/unix"
	"image"
	"io"
	"os"
)

func screen_size() *unix.Winsize {
	sz, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
            panic(err)
	}

	return sz
}

func get_sizes_CSI() (width uint16, height uint16) {
  fmt.Print("\033[14t")
  fmt.Scanf("\033[8;%d;%dt", &height, &width)
  return width, height
}

func createImage(width int, height int, color_positions []float64) *image.RGBA {
	fw := float64(width)
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
        colors := get_colors()

        colors = colors[:len(color_positions)]

	grad, err := colorgrad.NewGradient().
        HtmlColors(colors...).
                Domain(color_positions...).
		Build()

	if err != nil {
		fmt.Println(err)
		return nil
	}

	for x := 0; x < width; x++ {
		col := grad.At(float64(x) / fw)
		for y := 0; y < height; y++ {
			img.Set(x, y, col)
		}
	}

	return img
}

func _render(w io.Writer, img image.Image) error {
	bounds := img.Bounds()

	// f=32 => RGBA
	_, err := fmt.Fprintf(w, "\033_Gq=1,a=T,z=-99999,C=1,f=32,s=%d,v=%d,t=d,", bounds.Dx(), bounds.Dy())
	if err != nil {
		return err
	}

	buf := make([]byte, 0, 16384) // Multiple of 4 (RGBA)

	// var p streamPayload
	var p zlibPayload
	p.Reset(w)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if len(buf) == cap(buf) {
				if _, err = p.Write(buf); err != nil {
					return err
				}
				buf = buf[:0]
			}
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			buf = append(buf, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}

	if _, err = p.Write(buf); err != nil {
		return err
	}
	return p.Close()
}

func render_gradient(length int, color_positions []float64) error {
	size := screen_size()

	column := int(size.Xpixel / size.Col)
	row := int(size.Ypixel / size.Row)

	img := createImage(column*length, row, color_positions)
	return _render(os.Stdout, img)
}
