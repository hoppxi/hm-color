package color

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// Detect system theme best-effort; fallback to "light".
func DetectSystemTheme() string {
	// 1) GNOME (gsettings)
	if out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme").Output(); err == nil {
		txt := strings.ToLower(string(bytes.TrimSpace(out)))
		if strings.Contains(txt, "dark") {
			return "dark"
		}
		if strings.Contains(txt, "light") {
			return "light"
		}
	}

	// 2) KDE
	if v := os.Getenv("KDE_COLOR_SCHEME"); strings.Contains(strings.ToLower(v), "dark") {
		return "dark"
	}

	// 3) GTK_THEME env
	if v := os.Getenv("GTK_THEME"); strings.Contains(strings.ToLower(v), "dark") {
		return "dark"
	}

	// 4) Hyprland common env marker
	if strings.Contains(strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP")), "dark") {
		return "dark"
	}

	return "light"
}
