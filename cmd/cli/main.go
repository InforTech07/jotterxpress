package main

import (
	"jotterxpress/internal/adapters/cli"
	"os"
)

func main() {
	cliApp := cli.NewCLI()
	rootCmd := cliApp.SetupCommands()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
