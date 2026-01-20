package main

import (
	"github.com/JoseMariaMicoli/VaporTrace/cmd"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic" // Add this import
)

func main() {
	// Run environment sensing BEFORE starting the CLI
	logic.SenseEnvironment()
	
	// Execute the CLI commands
	cmd.Execute()
}