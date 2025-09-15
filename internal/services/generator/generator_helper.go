package generator

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/graphzc/wiresetgen/internal/models"
)

// For get the module name from go.mod file
func getModuleName(goModFile string) (string, error) {
	lines := strings.Split(goModFile, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[1], nil
			}
		}
	}
	return "", ErrInvalidGoModFile
}

// For extract string from package line
// Return packageName, nil when found package name
// Return nil, err when has an error
// Return nil, nil when no package in that line
func extractPackageName(line string) (*string, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, nil
	}

	if !strings.HasPrefix(line, "package") {
		return nil, nil
	}

	packageParts := strings.Fields(line)

	if len(packageParts) < 2 {
		return nil, ErrInvalidPackageName
	}

	if len(packageParts) > 2 && !strings.HasPrefix(packageParts[2], "//") {
		return nil, ErrInvalidPackageName
	}

	return &packageParts[1], nil
}

// For indicates @WireSet("name") annotation and extracts the data
func extractSetInfo(moduleName string, filePath string, fileContent string) []*models.WireSetInfo {
	setInfos := make([]*models.WireSetInfo, 0)

	lines := strings.Split(fileContent, "\n")
	packageName := ""
	for i := range lines {
		line := strings.TrimSpace(lines[i])

		if strings.Contains(line, "package") {
			packageParts := strings.Split(line, " ")
			if len(packageParts) > 1 {
				packageName = strings.TrimSpace(packageParts[1])
			}
		}

		if strings.Contains(line, "@WireSet(\"") {
			setInfo := models.WireSetInfo{}

			// Extract the annotation value
			annotationParts := strings.Split(line, "\"")
			if len(annotationParts) > 1 {
				setInfo.SetName = strings.TrimSpace(annotationParts[1])
			}

			// Extract the function name
			// Go to the next line
			if i+1 < len(lines) {
				nextLine := lines[i+1]
				functionParts := strings.Split(nextLine, " ")

				var fullFunctionName string
				if len(functionParts) > 1 {
					fullFunctionName = strings.TrimSpace(functionParts[1])
				}

				// Cut only the function name
				functionParts = strings.Split(fullFunctionName, "(")
				if len(functionParts) > 0 {
					setInfo.FunctionName = strings.TrimSpace(functionParts[0])
				}
			}

			// Set the package name
			setInfo.PackageName = packageName

			// Set the file path
			// Cut the latest path separator section from the file path
			pathParts := strings.Split(filePath, string(filepath.Separator))

			if len(pathParts) > 0 {
				pathParts = pathParts[:len(pathParts)-1]
				// Convert back to import path format (always use forward slashes for Go imports)
				cuttedFilePath := strings.Join(pathParts, "/")

				fullImportPath := path.Join(moduleName, cuttedFilePath)

				setInfo.ImportPath = strings.TrimSpace(fullImportPath)
			}

			setInfos = append(setInfos, &setInfo)
		}
	}

	return setInfos
}

// For extract the wiregen location from the file
// Return the wiregen location when found
// Return nil, nil when no wiregen location in that file
// Return nil, err when has an error
func extractWireGenLocation(filePath string, fileContent string) (*models.WireGenLocation, error) {
	lines := strings.Split(fileContent, "\n")

	isFound := false

	for i := range lines {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, "//go:build wireinject") {
			isFound = true
			continue
		}

		if isFound {
			packageName, err := extractPackageName(line)
			if err != nil {
				return nil, err
			}

			if packageName != nil {
				return &models.WireGenLocation{
					PackageName:   *packageName,
					DirectoryPath: filepath.Dir(filePath),
				}, nil
			}
		}
	}

	if isFound {
		return nil, fmt.Errorf("no package in wiregen file: %s", filePath)
	}

	return nil, nil
}
