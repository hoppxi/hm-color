package color

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GenerateMaterialPalette uses mcu-cli to generate a Material theme
// from a wallpaper image.
func GenerateMaterialPalette(wallpaper string, theme string) (map[string]string, error) {
	mode := strings.ToLower(theme)
	if mode != "light" && mode != "dark" {
		mode = "light"
	}

	// run mcuc with both themes and json output
	cmd := exec.Command("mcuc", "generate", "-i", wallpaper, "-f", "json", "-T", "both")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("mcuc generate error: %w", err)
	}

	// find JSON start (skip logs)
	start := bytes.IndexRune(output, '{')
	if start == -1 {
		return nil, fmt.Errorf("no JSON found in mcuc output:\n%s", output)
	}
	jsonBytes := output[start:]

	var palettes map[string]map[string]string
	if err := json.Unmarshal(jsonBytes, &palettes); err != nil {
		return nil, fmt.Errorf("parse json error: %w\nraw:\n%s", err, jsonBytes)
	}

	colors, ok := palettes[mode]
	if !ok {
		return nil, fmt.Errorf("theme %q not found in mcuc output", mode)
	}

	return colors, nil
}
