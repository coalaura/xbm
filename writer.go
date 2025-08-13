package xbm

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
)

type XBMOptions struct {
	Name string
}

func Encode(w io.Writer, m image.Image, opts ...XBMOptions) error {
	var name string

	if len(opts) > 0 && opts[0].Name != "" {
		name = opts[0].Name
	}

	name = sanitizeName(name)

	if name == "" {
		name = "image"
	}

	bounds := m.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if _, err := fmt.Fprintf(w, "#define %s_width %d\n", name, width); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "#define %s_height %d\n", name, height); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "static unsigned char %s_bits[] = {\n", name); err != nil {
		return err
	}

	var (
		byteVal    uint8
		bitCount   uint
		pixelCount int
	)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(m.At(x, y)).(color.Gray)

			if gray.Y < 128 {
				byteVal |= (1 << bitCount)
			}

			bitCount++
			pixelCount++

			if bitCount == 8 {
				if _, err := fmt.Fprintf(w, "0x%02x,", byteVal); err != nil {
					return err
				}

				byteVal = 0
				bitCount = 0
			}
		}

		if bitCount > 0 {
			if _, err := fmt.Fprintf(w, "0x%02x,", byteVal); err != nil {
				return err
			}

			byteVal = 0
			bitCount = 0
		}
	}

	if _, err := fmt.Fprintln(w, "};"); err != nil {
		return err
	}

	return nil
}

func sanitizeName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "image"
	}

	var out strings.Builder

	for i, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (i > 0 && r >= '0' && r <= '9') {
			out.WriteRune(r)
		} else {
			out.WriteByte('_')
		}
	}

	return out.String()
}
