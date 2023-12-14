// See https://sw.kovidgoyal.net/kitty/graphics-protocol.html.
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/mazznoer/colorgrad"
	"golang.org/x/sys/unix"
	"image"
	"image/png"
	"io"
	"math/rand"
	"os"
)

func screenSize() *unix.Winsize {
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
func _place_gradient(w io.Writer) {
	id := rand.Intn(1000000)
	fmt.Fprintf(w, "\033_Ga=p,C=1,z=-999,i=123,H=5,p=%d\033", id)
}

func _transmit(w io.Writer, img image.Image) {
	const chunkEncSize = 4096
	// const chunkEncSize = 48
	const chunkRawSize = (chunkEncSize / 4) * 3

	bounds := img.Bounds()

	// f=32 => RGBA
	fmt.Fprintf(w, "\033_Gq=1,a=t,C=0,f=32,i=123,s=%d,v=%d,t=d,", bounds.Dx(), bounds.Dy())

	bufRaw := make([]byte, 0, chunkRawSize)
	bufEnc := make([]byte, chunkEncSize)

	flush := func(last bool) {
		if len(bufRaw) == 0 {
			w.Write([]byte("m=0;\033\\"))
			return
		}
		if last {
			w.Write([]byte("m=0;"))
		} else {
			w.Write([]byte("m=1;"))
		}

		// fmt.Fprintln(os.Stderr, len(bufRaw), "=>", (len(bufRaw)+2)/3*4)

		base64.StdEncoding.Encode(bufEnc, bufRaw)
		w.Write(bufEnc[:(len(bufRaw)+2)/3*4])

		if last {
			w.Write([]byte("\033\\"))
		} else {
			w.Write([]byte("\033\\\033_G"))
			bufRaw = bufRaw[:0]
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if len(bufRaw)+4 > chunkRawSize {
				flush(false)
			}
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			bufRaw = append(bufRaw, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}
	flush(true)
}

func render_gradient(length int, color_positions []float64) error {
	size := screenSize()

	column := int(size.Xpixel / size.Col)
	row := int(size.Ypixel / size.Row)

	img := createImage(column*length, row, color_positions)

	// write image to file
	f, _ := os.Create("image.png")
	defer f.Close()
	_ = png.Encode(f, img)

	_transmit(os.Stdout, img)
	_place_gradient(os.Stdout)

	return nil
}
