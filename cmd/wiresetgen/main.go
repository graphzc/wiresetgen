package main

import (
	"github.com/graphzc/wiresetgen/internal/commands"
	"github.com/graphzc/wiresetgen/internal/handlers"
	"github.com/graphzc/wiresetgen/internal/repositories/files"
	"github.com/graphzc/wiresetgen/internal/services/generator"
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
