package color

import (
	"fmt"
	"strings"
)

func GenerateMaterialPalette(wallpaper string, theme string) (map[string]string, error) {
	r, g, b, err := DominantColor(wallpaper)
	if err != nil {
		return nil, fmt.Errorf("dominant color: %w", err)
	}

	seed := rgbToOKLCH(r, g, b)

	mode := strings.ToLower(theme)
	if mode == "system" || mode == "" {
		mode = DetectSystemTheme()
	}
	if mode != "light" && mode != "dark" {
		mode = "light"
	}

	pals := buildCorePalettes(seed)
	return tokensFromPalettes(pals, mode), nil
}
