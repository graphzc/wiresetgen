package generator

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/GraphZC/go-wireset-gen/internal/models"
	"github.com/GraphZC/go-wireset-gen/internal/repositories/files"
	"github.com/GraphZC/go-wireset-gen/internal/templates"
	"github.com/sirupsen/logrus"
)

var (
	ErrIsNotProjectRoot = errors.New("is not in project root directory")
	ErrInvalidGoModFile = errors.New("invalid go.mod file")
)

type Service interface {
	GenerateWireSet(verbose bool) error
}

type generatorServiceImpl struct {
	fileRepository files.Repository
}

func NewGenerateService(fileRepository files.Repository) Service {
	return &generatorServiceImpl{
		fileRepository: fileRepository,
	}
}

func (g *generatorServiceImpl) GenerateWireSet(verbose bool) error {
	goModFile, err := g.fileRepository.GetGoModFile()
	if err != nil {
		if errors.Is(err, files.ErrFileNotFound) {
			return ErrIsNotProjectRoot
		}

		return err
	}

	moduleName, err := getModuleName(string(goModFile))
	if err != nil {
		return ErrInvalidGoModFile
	}

	goFiles, err := g.fileRepository.ListAllGoFiles()
	if err != nil {
		return err
	}

	allSetInfo := make([]*models.WireSetInfo, 0)
	allWireGenInfo := make([]*models.WireGenInfo, 0)

	for _, file := range goFiles {
		fileContent, err := g.fileRepository.ReadFile(file)
		if err != nil {
			return err
		}

		extractedWireGenInfo := extractWireGenLocation(file, string(fileContent))
		if extractedWireGenInfo != nil {
			allWireGenInfo = append(allWireGenInfo, extractedWireGenInfo)

			if verbose {
				logrus.Info("Found wire gen file at", file)
			}

			// If this fils is wire gen file, skip extracting set info
			continue
		}

		extractedSetInfos := extractSetInfo(moduleName, file, string(fileContent))
		if len(extractedSetInfos) > 0 {
			if verbose {
				for _, setInfo := range extractedSetInfos {
					logrus.Infof("Found wire set %s for function %s\n", setInfo.SetName, setInfo.FunctionName)
				}
			}

			allSetInfo = append(allSetInfo, extractedSetInfos...)
		}
	}

	// Create import map
	// Create a map[importPath]alias
	importMap := make(map[string]string)
	aliasCounts := make(map[string]int)
	for _, setInfo := range allSetInfo {
		if _, exists := importMap[setInfo.ImportPath]; exists {
			continue
		}

		importPathParts := strings.Split(setInfo.ImportPath, "/")

		alias := importPathParts[len(importPathParts)-1]
		if _, exists := aliasCounts[alias]; exists {
			aliasCounts[alias]++
			alias = fmt.Sprintf("%s%d", alias, aliasCounts[alias])
		} else {
			aliasCounts[alias] = 1
		}

		importMap[setInfo.ImportPath] = alias
	}

	// Convert import map to []models.ImportTemplate
	imports := make([]*models.ImportTemplate, 0, len(importMap))
	for importPath, alias := range importMap {
		imports = append(imports, &models.ImportTemplate{
			Alias: alias,
			Path:  importPath,
		})
	}

	// Sort imports by import path
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})

	// Convert allSetInfo to map[setName][]*wireSetInfo
	setInfoMap := make(map[string][]*models.WireSetInfo)
	for _, setInfo := range allSetInfo {
		setInfoMap[setInfo.SetName] = append(setInfoMap[setInfo.SetName], setInfo)
	}

	for _, wireGenInfo := range allWireGenInfo {
		if verbose {
			logrus.Infof("Generating wire set for %s\n", wireGenInfo.DirectoryPath)
		}

		wireSetGenTemplate := templates.WireSetGenTemplate
		tmpl, err := template.New("wireSetGen").Parse(wireSetGenTemplate)
		if err != nil {
			return err
		}

		wireSetsMap := make(map[string]*models.WireSet, 0)
		for setName, setInfos := range setInfoMap {
			wireSetsMap[setName] = &models.WireSet{
				SetName: setName,
			}

			for _, info := range setInfos {
				wireSetsMap[setName].FuncPath = append(wireSetsMap[setName].FuncPath, fmt.Sprintf("%s.%s", info.PackageName, info.FunctionName))
			}
		}

		// Convert wireSetsMap to slice and sort by set name
		wireSets := make([]*models.WireSet, 0, len(wireSetsMap))
		for _, wireSet := range wireSetsMap {
			wireSets = append(wireSets, wireSet)
		}
		sort.Slice(wireSets, func(i, j int) bool {
			return wireSets[i].SetName < wireSets[j].SetName
		})

		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, models.WireSetGenTemplateModel{
			PackageName: wireGenInfo.PackageName,
			Imports:     imports,
			WireSets:    wireSets,
		}); err != nil {
			return err
		}

		logrus.Infoln("wireGenInfo", wireGenInfo)
		// Write the generated file
		err = g.fileRepository.WriteFile(wireGenInfo.DirectoryPath, "wire_set_gen.go", string(buf.Bytes()))
		if err != nil {
			return err
		}

		if verbose {
			logrus.Infof("Generated wire set file at %s/wire_set_gen.go\n", wireGenInfo.DirectoryPath)
		}
	}

	return nil
}
