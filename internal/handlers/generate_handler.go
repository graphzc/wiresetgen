package handlers

import (
	"github.com/GraphZC/go-wireset-gen/internal/services/generator"
)

type GenerateHandler interface {
	GenerateWireSet(verbose bool) error
}

type generateHandlerImpl struct {
	generatorService generator.Service
}

func NewGenerateHandler(generatorService generator.Service) GenerateHandler {
	return &generateHandlerImpl{
		generatorService: generatorService,
	}
}

func (g *generateHandlerImpl) GenerateWireSet(verbose bool) error {
	return g.generatorService.GenerateWireSet(verbose)
}
