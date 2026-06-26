package iconrender

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/fyne-io/oksvg"
	"github.com/srwiley/rasterx"
)

var IconSizes = map[string]int{
	"icon_16x16.png":      16,
	"icon_16x16@2x.png":   32,
	"icon_32x32.png":      32,
	"icon_32x32@2x.png":   64,
	"icon_128x128.png":    128,
	"icon_128x128@2x.png": 256,
	"icon_256x256.png":    256,
	"icon_256x256@2x.png": 512,
	"icon_512x512.png":    512,
	"icon_512x512@2x.png": 1024,
}

func RenderIcon(svgPath string, size int) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	FillRoundedRect(img, color.NRGBA{R: 255, G: 246, B: 226, A: 255}, size, size/5)
	FillRoundedRectInset(img, color.NRGBA{R: 255, G: 128, B: 0, A: 255}, size, size/12, size/5)
	FillRoundedRectInset(img, color.NRGBA{R: 255, G: 235, B: 196, A: 255}, size, size/7, size/6)

	icon, err := oksvg.ReadIcon(svgPath, oksvg.WarnErrorMode)
	if err != nil {
		return nil, err
	}

	padding := float64(size) * 0.20
	icon.SetTarget(padding, padding*0.82, float64(size)-padding*2, float64(size)-padding*1.76)
	scanner := rasterx.NewScannerGV(size, size, img, img.Bounds())
	dasher := rasterx.NewDasher(size, size, scanner)
	icon.Draw(dasher, 1)

	return img, nil
}

func FillRoundedRect(img *image.RGBA, c color.NRGBA, size, radius int) {
	FillRoundedRectInset(img, c, size, 0, radius)
}

func FillRoundedRectInset(img *image.RGBA, c color.NRGBA, size, inset, radius int) {
	min := inset
	max := size - inset
	r2 := radius * radius

	for y := min; y < max; y++ {
		for x := min; x < max; x++ {
			dx, dy := 0, 0
			if x < min+radius {
				dx = min + radius - x
			} else if x >= max-radius {
				dx = x - (max - radius - 1)
			}
			if y < min+radius {
				dy = min + radius - y
			} else if y >= max-radius {
				dy = y - (max - radius - 1)
			}
			if dx*dx+dy*dy <= r2 {
				img.Set(x, y, c)
			}
		}
	}
}

func WritePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
