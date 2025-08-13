package xbm

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"strconv"
)

func DecodeConfig(r io.Reader) (image.Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return image.Config{}, err
	}

	width, height, err := parseHeader(data)
	if err != nil {
		return image.Config{}, err
	}

	return image.Config{
		ColorModel: color.GrayModel,
		Width:      width,
		Height:     height,
	}, nil
}

func Decode(r io.Reader) (image.Image, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	width, height, err := parseHeader(data)
	if err != nil {
		return nil, err
	}

	hexBytes := extractHexBytes(data)
	if len(hexBytes) == 0 {
		return nil, fmt.Errorf("no pixel data found")
	}

	img := image.NewGray(image.Rect(0, 0, width, height))

	for i := 0; i < width*height; i++ {
		byteIndex := i / 8
		bitIndex := uint(i % 8)

		if byteIndex < len(hexBytes) {
			if (hexBytes[byteIndex] & (1 << bitIndex)) != 0 {
				img.Pix[i] = 255
			}
		}
	}

	return img, nil
}

func parseHeader(data []byte) (width, height int, err error) {
	for line := range bytes.SplitSeq(data, []byte{'\n'}) {
		if !bytes.HasPrefix(line, []byte("#define")) {
			continue
		}

		fields := bytes.Fields(line)

		if len(fields) >= 3 {
			val, _ := strconv.Atoi(string(fields[2]))

			if bytes.HasSuffix(fields[1], []byte("_width")) {
				width = val
			} else if bytes.HasSuffix(fields[1], []byte("_height")) {
				height = val
			}
		}
	}

	if width <= 0 || height <= 0 {
		return 0, 0, fmt.Errorf("invalid width/height")
	}

	return width, height, nil
}

func extractHexBytes(data []byte) []byte {
	start := bytes.IndexByte(data, '{')
	end := bytes.LastIndexByte(data, '}')

	if start == -1 || end == -1 || end <= start {
		return nil
	}

	fields := bytes.Fields(data[start+1 : end])
	out := make([]byte, 0, len(fields))

	for _, f := range fields {
		if len(f) > 2 && (f[0:2][0] == '0' && (f[1] == 'x' || f[1] == 'X')) {
			if v, err := strconv.ParseUint(string(f), 0, 8); err == nil {
				out = append(out, byte(v))
			}
		}
	}

	return out
}
