package handlers

import (
	"log"

	"github.com/GraphZC/go-wireset-gen/internal/services"
)

type GenerateHandler interface {
	GenerateWireSet() error
}

type generateHandlerImpl struct {
	fileService services.FileService
}

func NewGenerateHandler(fileService services.FileService) GenerateHandler {
	return &generateHandlerImpl{
		fileService: fileService,
	}
}

func (g *generateHandlerImpl) GenerateWireSet() error {
	log.Println("Generating wire set...")
	return nil
}
