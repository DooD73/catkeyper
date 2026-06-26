package iconrender

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestFillRoundedRectInset(t *testing.T) {
	const size = 32
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	fill := color.NRGBA{R: 255, G: 128, B: 0, A: 255}

	FillRoundedRectInset(img, fill, size, 4, 6)

	corner := img.At(0, 0)
	if _, _, _, a := corner.RGBA(); a != 0 {
		t.Fatalf("expected transparent corner, got alpha %d", a>>8)
	}

	center := color.RGBAModel.Convert(img.At(size/2, size/2)).(color.RGBA)
	want := color.RGBAModel.Convert(fill).(color.RGBA)
	if center != want {
		t.Fatalf("expected center pixel %#v, got %#v", want, center)
	}
}

func TestRenderIcon(t *testing.T) {
	svgPath := filepath.Join("..", "..", "assets", "openmoji-cat-face.svg")
	const size = 64

	img, err := RenderIcon(svgPath, size)
	if err != nil {
		t.Fatalf("RenderIcon: %v", err)
	}

	if img.Bounds().Dx() != size || img.Bounds().Dy() != size {
		t.Fatalf("expected %dx%d image, got %dx%d", size, size, img.Bounds().Dx(), img.Bounds().Dy())
	}

	center := img.At(size/2, size/2)
	if _, _, _, a := center.RGBA(); a == 0 {
		t.Fatal("expected rendered icon center to be opaque")
	}
}

func TestWritePNG(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	fill := color.NRGBA{R: 10, G: 20, B: 30, A: 255}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, fill)
		}
	}

	path := filepath.Join(t.TempDir(), "icon.png")
	if err := WritePNG(path, img); err != nil {
		t.Fatalf("WritePNG: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open png: %v", err)
	}
	defer f.Close()

	decoded, err := png.Decode(f)
	if err != nil {
		t.Fatalf("png.Decode: %v", err)
	}

	if decoded.Bounds().Dx() != 8 || decoded.Bounds().Dy() != 8 {
		t.Fatalf("expected 8x8 decoded image, got %dx%d", decoded.Bounds().Dx(), decoded.Bounds().Dy())
	}
}
