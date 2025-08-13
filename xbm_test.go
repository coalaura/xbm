package xbm

import (
	"crypto/sha256"
	"encoding/hex"
	"image/png"
	"os"
	"testing"
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

	if info.Width != 800 {
		t.Fatalf("expected width of 800px got: %d", info.Width)
	}

	if info.Height != 500 {
		t.Fatalf("expected height of 500px got: %d", info.Height)
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

	if bounds.Dx() != 800 {
		t.Fatalf("expected width of 800px got: %d", bounds.Dx())
	}

	if bounds.Dy() != 500 {
		t.Fatalf("expected height of 500px got: %d", bounds.Dy())
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
	if hexHash != "a75f42d70168657910597e294501338018a36846e55f8355c88f5e0f49b5348e" {
		t.Fatalf("hash mismatch: %s != a75f42d70168657910597e294501338018a36846e55f8355c88f5e0f49b5348e", hexHash)
	}
}
