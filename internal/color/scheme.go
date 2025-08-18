package color

import (
	"math"
	"strings"
)

// Build core palettes from a seed OKLCH.
// We follow Material's idea: core palettes for primary/secondary/tertiary/neutral/neutral-variant/error.
// Here we vary chroma per palette using multipliers, and we generate a standard tone set.
type palettes struct {
	primary  map[int]string
	secondary map[int]string
	tertiary map[int]string
	neutral  map[int]string
	neutralV map[int]string
	errorP   map[int]string
}

var toneSet = []int{0, 4, 6, 10, 12, 17, 20, 22, 24, 30, 40, 50, 60, 70, 80, 87, 90, 92, 94, 96, 98, 99, 100}

// clamp chroma to a safe displayable range for given lightness
func clampChromaForDisplay(l, c float64) float64 {
	// Empirical safe bound to avoid gamut clipping
	maxC := 0.4 * (1 - 0.5*math.Abs(l-0.5))
	if maxC < 0.02 {
		maxC = 0.02
	}
	if c > maxC {
		return maxC
	}
	if c < 0 {
		return 0
	}
	return c
}

func buildTonalPalette(h float64, cBase float64) map[int]string {
	out := make(map[int]string, len(toneSet))
	for _, t := range toneSet {
		L := toneToOKL(float64(t))
		C := clampChromaForDisplay(L, cBase)
		col := OKLCH{L: L, C: C, H: normalizeHue(h)}
		r, g, b := oklchToRGB(col)
		out[t] = hexFromRGB(r, g, b)
	}
	return out
}

func normalizeHue(h float64) float64 {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	return h
}

func buildCorePalettes(seed OKLCH) palettes {
	// Base chroma multipliers akin to M3 behavior
	cPri := seed.C
	if cPri < 0.06 {
		cPri = 0.06
	}
	cSec := cPri * 0.6
	cTer := cPri * 0.8
	cNeu := cPri * 0.15
	cNuv := cPri * 0.2

	p := palettes{
		primary:   buildTonalPalette(seed.H, cPri),
		secondary: buildTonalPalette(seed.H+60, cSec),
		tertiary:  buildTonalPalette(seed.H-60, cTer),
		neutral:   buildTonalPalette(seed.H, cNeu),
		neutralV:  buildTonalPalette(seed.H, cNuv),
		errorP:    buildTonalPalette(25, 0.25), // stable error palette
	}
	return p
}

func pick(t map[int]string, tone int) string { return t[tone] }

// map to Material-like tokens; both light & dark sets return the full kebab-case map
func tokensFromPalettes(p palettes, mode string) map[string]string {
	mode = strings.ToLower(mode)
	out := map[string]string{
		"shadow":        "#000000",
		"scrim":         "#000000",
		"surface-tint":  pick(p.primary, 40),
		"surface-variant": pick(p.neutralV, 90),
		"on-surface-variant": pick(p.neutralV, 30),
		"outline":           pick(p.neutralV, 50),
		"outline-variant":   pick(p.neutralV, 80),
	}

	if mode == "dark" {
		out["background"] = pick(p.neutral, 6)
		out["on-background"] = pick(p.neutral, 90)
		out["surface"] = pick(p.neutral, 6)
		out["surface-dim"] = pick(p.neutral, 6)
		out["surface-bright"] = pick(p.neutral, 24)
		out["surface-container-lowest"] = pick(p.neutral, 4)
		out["surface-container-low"] = pick(p.neutral, 10)
		out["surface-container"] = pick(p.neutral, 12)
		out["surface-container-high"] = pick(p.neutral, 17)
		out["surface-container-highest"] = pick(p.neutral, 22)
		out["on-surface"] = pick(p.neutral, 90)
		out["inverse-surface"] = pick(p.neutral, 90)
		out["inverse-on-surface"] = pick(p.neutral, 20)

		out["primary"] = pick(p.primary, 80)
		out["on-primary"] = pick(p.primary, 20)
		out["primary-container"] = pick(p.primary, 30)
		out["on-primary-container"] = pick(p.primary, 90)
		out["inverse-primary"] = pick(p.primary, 40)

		out["secondary"] = pick(p.secondary, 80)
		out["on-secondary"] = pick(p.secondary, 20)
		out["secondary-container"] = pick(p.secondary, 30)
		out["on-secondary-container"] = pick(p.secondary, 90)

		out["tertiary"] = pick(p.tertiary, 80)
		out["on-tertiary"] = pick(p.tertiary, 20)
		out["tertiary-container"] = pick(p.tertiary, 30)
		out["on-tertiary-container"] = pick(p.tertiary, 90)

		out["error"] = pick(p.errorP, 80)
		out["on-error"] = pick(p.errorP, 20)
		out["error-container"] = pick(p.errorP, 30)
		out["on-error-container"] = pick(p.errorP, 90)
		return out
	}

	// light
	out["background"] = pick(p.neutral, 99)
	out["on-background"] = pick(p.neutral, 10)
	out["surface"] = pick(p.neutral, 99)
	out["surface-dim"] = pick(p.neutral, 87)
	out["surface-bright"] = pick(p.neutral, 98)
	out["surface-container-lowest"] = pick(p.neutral, 100)
	out["surface-container-low"] = pick(p.neutral, 96)
	out["surface-container"] = pick(p.neutral, 94)
	out["surface-container-high"] = pick(p.neutral, 92)
	out["surface-container-highest"] = pick(p.neutral, 90)
	out["on-surface"] = pick(p.neutral, 10)
	out["inverse-surface"] = pick(p.neutral, 20)
	out["inverse-on-surface"] = pick(p.neutral, 95)

	out["primary"] = pick(p.primary, 40)
	out["on-primary"] = pick(p.primary, 100)
	out["primary-container"] = pick(p.primary, 90)
	out["on-primary-container"] = pick(p.primary, 10)
	out["inverse-primary"] = pick(p.primary, 80)

	out["secondary"] = pick(p.secondary, 40)
	out["on-secondary"] = pick(p.secondary, 100)
	out["secondary-container"] = pick(p.secondary, 90)
	out["on-secondary-container"] = pick(p.secondary, 10)

	out["tertiary"] = pick(p.tertiary, 40)
	out["on-tertiary"] = pick(p.tertiary, 100)
	out["tertiary-container"] = pick(p.tertiary, 90)
	out["on-tertiary-container"] = pick(p.tertiary, 10)

	out["error"] = pick(p.errorP, 40)
	out["on-error"] = pick(p.errorP, 100)
	out["error-container"] = pick(p.errorP, 90)
	out["on-error-container"] = pick(p.errorP, 10)

	return out
}
