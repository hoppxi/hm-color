package color

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// DominantColor returns the most frequent sRGB color (8-bit per channel) sampled sparsely.
func DominantColor(path string) (r, g, b uint8, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, 0, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return 0, 0, 0, err
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	stepX := 4
	stepY := 4
	if w*h > 2_000_000 {
		stepX, stepY = 8, 8
	} else if w*h < 200_000 {
		stepX, stepY = 2, 2
	}

	counts := make(map[uint32]int, (w/stepX)*(h/stepY))
	var maxColor uint32
	maxCount := -1

	for y := bounds.Min.Y; y < bounds.Max.Y; y += stepY {
		for x := bounds.Min.X; x < bounds.Max.X; x += stepX {
			R, G, B, _ := img.At(x, y).RGBA()
			r8 := uint8(R >> 8)
			g8 := uint8(G >> 8)
			b8 := uint8(B >> 8)

			key := (uint32(r8) << 16) | (uint32(g8) << 8) | uint32(b8)
			c := counts[key] + 1
			counts[key] = c
			if c > maxCount {
				maxCount = c
				maxColor = key
			}
		}
	}

	return uint8(maxColor >> 16), uint8((maxColor >> 8) & 0xFF), uint8(maxColor & 0xFF), nil
}
