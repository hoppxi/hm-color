package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ApplyFlake() {
	fmt.Println("Activating home-manager")

	activatePath := filepath.Join(os.Getenv("HOME"), ".local/state/nix/profiles/home-manager/activate")
	cmd := exec.Command(activatePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

}
