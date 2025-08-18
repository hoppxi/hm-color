package cmd

import (
	"os"
	"path/filepath"

	"hm-color/internal/config"
	"hm-color/internal/watcher"

	"github.com/spf13/cobra"
)

var cfg = &config.Config{}

var rootCmd = &cobra.Command{
	Use:   "hm-color",
	Short: "Dynamic theming tool for NixOS with swww wallpaper manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		
		if cfg.SwwwCache == "" {
			xdgCache := os.Getenv("XDG_CACHE_HOME")
			if xdgCache == "" {
				xdgCache = filepath.Join(os.Getenv("HOME"), ".cache")
			}
			cfg.SwwwCache = filepath.Join(xdgCache, "swww")
		}

		if !(cfg.JSONStdout || cfg.JSONOut != "" ||
			cfg.SCSSStdout || cfg.SCSSOut != "" ||
			cfg.CSSStdout || cfg.CSSOut != "" ||
			cfg.NixStdout || cfg.NixOut != "") {
			cfg.JSONStdout = true
		}


		return watcher.Start(cfg)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfg.Theme, "theme", "t", "system", "What base color to use [light|dark|system]")
	rootCmd.PersistentFlags().StringVar(&cfg.SwwwCache, "swww-cache", "", "Path to swww cache dir")
	rootCmd.PersistentFlags().BoolVarP(&cfg.NixStdout, "nix", "n", false, "Output theme as Nix to stdout")
	rootCmd.PersistentFlags().StringVar(&cfg.NixOut, "nix-out", "", "Path to nix output file")
	rootCmd.PersistentFlags().BoolVarP(&cfg.JSONStdout, "json", "j", false, "Output theme as JSON to stdout")
	rootCmd.PersistentFlags().StringVar(&cfg.JSONOut, "json-out", "", "Write JSON output to a file")
	rootCmd.PersistentFlags().BoolVarP(&cfg.SCSSStdout, "scss", "s", false, "Output theme as SCSS variables to stdout")
	rootCmd.PersistentFlags().StringVar(&cfg.SCSSOut, "scss-out", "", "Write SCSS variables output to a file")
	rootCmd.PersistentFlags().BoolVarP(&cfg.CSSStdout, "css", "c", false, "Output theme as CSS variables to stdout")
	rootCmd.PersistentFlags().StringVar(&cfg.CSSOut, "css-out", "", "Write CSS variables output to a file")
	rootCmd.PersistentFlags().StringVarP(&cfg.FlakePath, "flake", "f", "", "Run home-manager switch --flake <path> after generation")
	rootCmd.PersistentFlags().BoolVar(&cfg.Prune, "prune", false, "Delete the previous home-manager generation after switching")
	rootCmd.PersistentFlags().BoolVarP(&cfg.GitCommit, "git-commit", "g", false, "Commit nix config with message 'update color'")
}

func Execute() error {
	return rootCmd.Execute()
}
