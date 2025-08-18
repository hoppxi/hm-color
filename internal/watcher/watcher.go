package watcher

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"hm-color/internal/color"
	"hm-color/internal/config"
	"hm-color/internal/exec"
	"hm-color/internal/formats"
	"hm-color/internal/output"

	"github.com/fsnotify/fsnotify"
)

func Start(cfg *config.Config) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	if err := w.Add(cfg.SwwwCache); err != nil {
		return err
	}

	log.Printf("Watching %s for wallpaper changes...\n", cfg.SwwwCache)

	for {
		select {
		case event := <-w.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				data, err := os.ReadFile(event.Name)
				if err != nil {
					log.Printf("failed to read cache file %s: %v", event.Name, err)
					continue
				}

				line := strings.TrimSpace(string(data))
				if line == "" {
					continue
				}

				parts := strings.Fields(line)
				if len(parts) == 0 {
					continue
				}

				// Always take the last field as the wallpaper path
				wallpaper := parts[len(parts)-1]

				// Ensure path is absolute (swww sometimes writes relative paths)
				if !filepath.IsAbs(wallpaper) {
					wallpaper = filepath.Join(filepath.Dir(event.Name), wallpaper)
				}

				// Validate the file exists
				if _, err := os.Stat(wallpaper); err != nil {
					log.Printf("skipping invalid wallpaper path %q: %v", wallpaper, err)
					continue
				}

				log.Printf("Detected new wallpaper: %s", wallpaper)

				colors, err := color.GenerateMaterialPalette(wallpaper, cfg.Theme)
				if err != nil {
					log.Println("Color extraction failed:", err)
					continue
				}

				// Format outputs
				if cfg.JSONStdout || cfg.JSONOut != "" {
					out := formats.FormatJSON(colors)
					output.Handle("json", out, cfg.JSONStdout, cfg.JSONOut)
				}
				if cfg.SCSSStdout || cfg.SCSSOut != "" {
					out := formats.FormatSCSS(colors)
					output.Handle("scss", out, cfg.SCSSStdout, cfg.SCSSOut)
				}
				if cfg.CSSStdout || cfg.CSSOut != "" {
					out := formats.FormatCSS(colors)
					output.Handle("css", out, cfg.CSSStdout, cfg.CSSOut)
				}
				if cfg.NixStdout || cfg.NixOut != "" {
					out := formats.FormatNix(colors)
					output.Handle("nix", out, cfg.NixStdout, cfg.NixOut)
				}

				// Optional flake handling
				if cfg.FlakePath != "" {
					exec.ApplyFlake(cfg.FlakePath, cfg.Prune, cfg.GitCommit)
				}
			}
		case err := <-w.Errors:
			log.Println("watcher error:", err)
		}
	}
}
