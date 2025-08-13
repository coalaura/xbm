# xbm

Pure Go XBM image encoder/decoder.

## Install
```sh
go get -u github.com/coalaura/xbm
```

## Usage

### Decode
```go
package main

import (
	"fmt"
	"image"
	"os"

	"github.com/coalaura/xbm"
)

func main() {
	f, _ := os.Open("test.xbm")
	defer f.Close()

	img, err := xbm.Decode(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(img.Bounds())
}
```

### Encode
```go
package main

import (
	"image"
	"image/color"
	"os"

	"github.com/coalaura/xbm"
)

func main() {
	img := image.NewGray(image.Rect(0, 0, 8, 8))
	img.SetGray(3, 3, color.Gray{Y: 255})

	f, _ := os.Create("out.xbm")
	defer f.Close()

	xbm.Encode(f, img, xbm.XBMOptions{Name: "test"})
}
```