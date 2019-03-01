package app

import "github.com/stakater/Whitelister/internal/pkg/cmd"

// Run runs th Whitelister command
func Run() error {
	cmd := cmd.NewWhitelisterCommand()
	return cmd.Execute()
}
