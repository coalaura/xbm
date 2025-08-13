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

	bytesPerRow := (width + 7) / 8
	byteIndex := 0

	for y := range height {
		for x := range width {
			b := hexBytes[byteIndex+x/8]
			bit := (b >> uint(x%8)) & 1

			if bit == 0 {
				img.SetGray(x, y, color.Gray{Y: 255})
			} else {
				img.SetGray(x, y, color.Gray{Y: 0})
			}
		}

		byteIndex += bytesPerRow
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

	fields := bytes.Split(data[start+1:end], []byte(","))
	out := make([]byte, 0, len(fields))

	for _, f := range fields {
		f = bytes.TrimSpace(f)

		if len(f) > 2 && (f[0:2][0] == '0' && (f[1] == 'x' || f[1] == 'X')) {
			if v, err := strconv.ParseUint(string(f), 0, 8); err == nil {
				out = append(out, byte(v))
			}
		}
	}

	return out
}
