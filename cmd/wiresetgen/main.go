package main

import (
	"github.com/GraphZC/go-wireset-gen/internal/command"
	"github.com/GraphZC/go-wireset-gen/internal/handlers"
	"github.com/GraphZC/go-wireset-gen/internal/services"
)

func main() {
	// Initialize services
	fileService := services.NewFileService()

	// Initialize handlers
	generateHandler := handlers.NewGenerateHandler(fileService)

	// Initialize commands
	rootCmd := command.NewRootCommand()
	rootCmd.AddCommand(command.NewGenerateCommand(generateHandler))

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
