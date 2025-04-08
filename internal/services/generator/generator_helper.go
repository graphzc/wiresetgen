package generator

import (
	"os"
	"path"
	"strings"

	"github.com/GraphZC/go-wireset-gen/internal/models"
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
			// Cut the latest / section from the file path
			pathParts := strings.Split(filePath, string(os.PathSeparator))

			if len(pathParts) > 0 {
				pathParts = pathParts[:len(pathParts)-1]
				cuttedFilePath := strings.Join(pathParts, "/")

				fullImportPath := path.Join(moduleName, cuttedFilePath)

				setInfo.ImportPath = strings.TrimSpace(fullImportPath)
			}

			setInfos = append(setInfos, &setInfo)
		}
	}

	return setInfos
}

func extractWireGenLocation(filePath string, fileContent string) *models.WireGenInfo {
	lines := strings.Split(fileContent, "\n")

	isFound := false

	for i := range lines {
		line := strings.TrimSpace(lines[i])

		if strings.Contains(line, "//go:build wireinject") {
			isFound = true
			continue
		}

		if isFound {
			strings.Contains(line, "package")
			packageParts := strings.Split(line, "package")
			if len(packageParts) > 1 {
				packageParts = strings.Split(packageParts[1], " ")
				if len(packageParts) > 1 {
					return &models.WireGenInfo{
						PackageName:   strings.TrimSpace(packageParts[1]),
						DirectoryPath: path.Dir(filePath),
					}
				}
			}
		}
	}

	return nil
}
