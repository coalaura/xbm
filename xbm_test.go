package xbm

import (
	"crypto/sha256"
	"encoding/hex"
	"image/png"
	"os"
	"testing"
)

const (
	expectedWidth  = 206
	expectedHeight = 187
	expectedHash   = "e46b5718e44f2736d7d269e9c616128faa15558a302ea0a380f2eb1c2e64b8e7"
)

func Test_DecodeConfig(t *testing.T) {
	file, err := os.OpenFile("image.xbm", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	info, err := DecodeConfig(file)
	if err != nil {
		t.Fatal(err)
	}

	if info.Width != expectedWidth {
		t.Fatalf("expected width of %dpx got: %d", expectedWidth, info.Width)
	}

	if info.Height != expectedHeight {
		t.Fatalf("expected height of %dpx got: %d", expectedHeight, info.Height)
	}
}

func Test_Decode(t *testing.T) {
	file, err := os.OpenFile("image.xbm", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	img, err := Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	bounds := img.Bounds()

	if bounds.Dx() != expectedWidth {
		t.Fatalf("expected width of %dpx got: %d", expectedWidth, bounds.Dx())
	}

	if bounds.Dy() != expectedHeight {
		t.Fatalf("expected height of %dpx got: %d", expectedHeight, bounds.Dy())
	}
}

func Test_Encode(t *testing.T) {
	file, err := os.OpenFile("image.png", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.New()

	err = Encode(hash, img, XBMOptions{
		Name: "cats",
	})
	if err != nil {
		t.Fatal(err)
	}

	hexHash := hex.EncodeToString(hash.Sum(nil))
	if hexHash != expectedHash {
		t.Fatalf("hash mismatch: %s != %s", hexHash, expectedHash)
	}
}
