package generator

import (
	"testing"

	"github.com/graphzc/wiresetgen/internal/models"
	"github.com/stretchr/testify/require"
)

func Test_getModuleName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		goModFile    string
		expectedName string
		expectedErr  error
	}{
		{
			name:         "Valid go.mod file",
			goModFile:    "module github.com/graphzc/wiresetgen\n",
			expectedName: "github.com/graphzc/wiresetgen",
			expectedErr:  nil,
		},
		{
			name:         "Invalid go.mod file",
			goModFile:    "invalid content\n",
			expectedName: "",
			expectedErr:  ErrInvalidGoModFile,
		},
		{
			name:         "Empty go.mod file",
			goModFile:    "",
			expectedName: "",
			expectedErr:  ErrInvalidGoModFile,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			moduleName, err := getModuleName(tc.goModFile)

			require.Equal(t, tc.expectedName, moduleName)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func Test_extractWireGenLocation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		filePath     string
		fileContent  string
		expectedInfo *models.WireGenLocation
	}{
		{
			name:     "Valid file with wireinject",
			filePath: "internal/wire/wire.go",
			fileContent: `
			//go:build wireinject
			// +build wireinject

			package wire
			`,
			expectedInfo: &models.WireGenLocation{
				PackageName:   "internal/services/generator",
				DirectoryPath: "internal/services/generator",
			},
		},
		{
			name:         "No wireinject in file",
			filePath:     "internal/wire/not_wire_file.go",
			fileContent:  "package notwire",
			expectedInfo: nil,
		},
		{
			name:         "No package in file, only wireinject",
			filePath:     "internal/wire/not_wire_file.go",
			fileContent:  "//go:build wireinject",
			expectedInfo: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			wireGenLocation := extractWireGenLocation(tc.filePath, tc.fileContent)

			require.Equal(t, tc.expectedInfo, wireGenLocation)
		})
	}
}
