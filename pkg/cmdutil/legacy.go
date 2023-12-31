package cmdutil

import (
	"fmt"
	"os"

	"github.com/khulnasoft-lab/goctl/v2/internal/config"
)

// TODO: consider passing via Factory
// TODO: support per-hostname settings
func DetermineEditor(cf func() (config.Config, error)) (string, error) {
	editorCommand := os.Getenv("GOCTL_EDITOR")
	if editorCommand == "" {
		cfg, err := cf()
		if err != nil {
			return "", fmt.Errorf("could not read config: %w", err)
		}
		editorCommand = cfg.Editor("")
	}

	return editorCommand, nil
}
