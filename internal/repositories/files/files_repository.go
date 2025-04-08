package files

import (
	"os"
	"path"
)

const BASE_DIR = "./"

type Repository interface {
	GetGoModFile() (string, error)
	ListAllGoFiles() ([]string, error)
	ReadFile(filePath string) (string, error)
	WriteFile(directory string, fileName string, data string) error
}

type repositoryImpl struct{}

func NewFileRepository() Repository {
	return &repositoryImpl{}
}

func (f *repositoryImpl) ReadFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrFileNotFound
		}
		return "", err
	}

	return string(data), nil
}

func (f *repositoryImpl) GetGoModFile() (string, error) {
	return f.ReadFile(path.Join(BASE_DIR, "go.mod"))
}

func (f *repositoryImpl) ListAllGoFiles() ([]string, error) {
	pendingDirectory := []string{BASE_DIR}
	goFiles := make([]string, 0)

	for len(pendingDirectory) > 0 {
		currentDir := pendingDirectory[0]
		pendingDirectory = pendingDirectory[1:]

		files, err := os.ReadDir(currentDir)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			filePath := path.Join(currentDir, file.Name())
			if file.IsDir() {
				pendingDirectory = append(pendingDirectory, filePath)
			} else if path.Ext(file.Name()) == ".go" {
				goFiles = append(goFiles, filePath)
			}
		}
	}

	return goFiles, nil
}

func (f *repositoryImpl) WriteFile(directory string, fileName string, data string) error {
	// Convert string to byte slice
	dataBytes := []byte(data)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return err
	}
	// Create the file path
	filePath := path.Join(directory, fileName)

	return os.WriteFile(filePath, dataBytes, 0644)
}
