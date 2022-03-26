package helpers

import (
	"image"
	"image/draw"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

func LoadImage(path string) ([]uint8, int32, int32) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba.Pix, width, height
}

func ReflectImageVertically(imageData []uint8, width int32, alfa bool) []uint8 {
	reflected := make([]uint8, 0, len(imageData))
	var stride int
	if alfa {
		stride = int(width * 4)
	} else {
		stride = int(width * 3)
	}

	for i := len(imageData) - stride; i >= 0; i = i - stride {
		for j := i; j < stride+i; j++ {
			reflected = append(reflected, imageData[j])
		}
	}
	return reflected
}
