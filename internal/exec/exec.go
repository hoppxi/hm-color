package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ApplyFlake(flake string, prune bool, gitCommit bool) {
	fmt.Println("Applying flake:", flake)

	// Extract repo path from flake string "path/to/flake#user@host"
	repoPath := flake
	if idx := strings.Index(flake, "#"); idx != -1 {
		repoPath = flake[:idx]
	}
	repoPath, _ = filepath.Abs(repoPath)

	exec.Command("git", "-C", repoPath, "add", ".").Run()

	if gitCommit {
		exec.Command("git", "-C", repoPath, "commit", "-m", "update color").Run()
	}

	cmd := exec.Command("home-manager", "switch", "--flake", flake)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	if prune {
		prunePreviousGeneration()
	}
}

// Finds the previous home-manager generation and deletes it
func prunePreviousGeneration() {
	cmd := exec.Command("home-manager", "generations")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to list generations:", err)
		return
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		fmt.Println("No previous generation found to prune")
		return
	}

	// Line format: "2025-08-18 21:44 : id 123"
	prevLine := lines[1]
	parts := strings.Fields(prevLine)
	if len(parts) < 4 {
		fmt.Println("Could not parse generation line:", prevLine)
		return
	}

	id := parts[len(parts)-1] // last token is generation ID
	fmt.Println("Pruning generation:", id)

	exec.Command("home-manager", "remove-generations", id).Run()
}
