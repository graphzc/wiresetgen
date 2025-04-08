package main

import (
	"github.com/GraphZC/go-wireset-gen/internal/commands"
	"github.com/GraphZC/go-wireset-gen/internal/handlers"
	"github.com/GraphZC/go-wireset-gen/internal/repositories/files"
	"github.com/GraphZC/go-wireset-gen/internal/services/generator"
)

func main() {
	// Initialize repositories
	fileRepository := files.NewFileRepository()

	// Initialize services
	generatorService := generator.NewGenerateService(fileRepository)

	// Initialize handlers
	generateHandler := handlers.NewGenerateHandler(generatorService)

	// Initialize commands
	rootCmd := commands.NewRootCommand()
	rootCmd.AddCommand(commands.NewGenerateCommand(generateHandler))

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
