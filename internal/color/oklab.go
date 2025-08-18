package color

import (
	"fmt"
	"math"
)

// -------- sRGB <-> Linear --------

func srgbToLinear(u float64) float64 {
	if u <= 0.04045 {
		return u / 12.92
	}
	return math.Pow((u+0.055)/1.055, 2.4)
}

func linearToSrgb(u float64) float64 {
	if u <= 0.0031308 {
		return 12.92 * u
	}
	return 1.055*math.Pow(u, 1/2.4) - 0.055
}

// -------- sRGB (0-255) <-> OKLab --------

// OKLab forward per Björn Ottosson
func rgbToOKLab(r8, g8, b8 uint8) (L, a, c float64) {
	r := srgbToLinear(float64(r8) / 255.0)
	g := srgbToLinear(float64(g8) / 255.0)
	b := srgbToLinear(float64(b8) / 255.0)

	// Linear RGB -> LMS (nonlinear space for OKLab)
	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b

	l_ := math.Cbrt(l)
	m_ := math.Cbrt(m)
	s_ := math.Cbrt(s)

	L = 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	a = 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	c = 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_
	return
}

func oklabToRGB(L, a, b float64) (r8, g8, b8 uint8) {
	l_ := L + 0.3963377774*a + 0.2158037573*b
	m_ := L - 0.1055613458*a - 0.0638541728*b
	s_ := L - 0.0894841775*a - 1.2914855480*b

	l := l_ * l_ * l_
	m := m_ * m_ * m_
	s := s_ * s_ * s_

	r := +4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	g := -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	bc := +0.0041960863*l - 0.7034186147*m + 1.7076147010*s

	r = clamp01(r)
	g = clamp01(g)
	bc = clamp01(bc)

	r = linearToSrgb(r)
	g = linearToSrgb(g)
	bc = linearToSrgb(bc)

	r = clamp01(r)
	g = clamp01(g)
	bc = clamp01(bc)

	return uint8(math.Round(r * 255)), uint8(math.Round(g * 255)), uint8(math.Round(bc * 255))
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// -------- OKLCH helpers --------

type OKLCH struct {
	L float64 // 0..1
	C float64 // 0..~0.322
	H float64 // degrees 0..360
}

func rgbToOKLCH(r, g, b uint8) OKLCH {
	L, a, bb := rgbToOKLab(r, g, b)
	h := math.Atan2(bb, a) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	c := math.Sqrt(a*a + bb*bb)
	return OKLCH{L: L, C: c, H: h}
}

func oklchToRGB(o OKLCH) (uint8, uint8, uint8) {
	a := o.C * math.Cos(o.H*math.Pi/180)
	b := o.C * math.Sin(o.H*math.Pi/180)
	return oklabToRGB(o.L, a, b)
}

func hexFromRGB(r, g, b uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// Tone is 0..100 (Material-ish). Convert to OKLCH L ~ 0..1 nonlinearly.
// OKLab L is already perceptual; a simple mapping works well:
func toneToOKL(Ltone float64) float64 {
	// 0..100 -> 0..1
	return clamp01(Ltone / 100.0)
}
