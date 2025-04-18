package generator

import (
	"errors"
	"testing"

	"github.com/graphzc/wiresetgen/internal/models"
	"github.com/graphzc/wiresetgen/pkg/utils"
	"github.com/stretchr/testify/assert"
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
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()

			moduleName, err := getModuleName(tc.goModFile)

			assert.Equal(tt, tc.expectedName, moduleName)
			assert.ErrorIs(tt, err, tc.expectedErr)
		})
	}
}

func Test_extractPackageName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		line          string
		expectedName  *string
		expectedError error
	}{
		{
			name:          "Valid package line",
			line:          "package main",
			expectedName:  utils.ToPointer("main"),
			expectedError: nil,
		},
		{
			name:          "Invalid package line",
			line:          "invalid package line",
			expectedName:  nil,
			expectedError: nil,
		},
		{
			name:          "Empty package line",
			line:          "",
			expectedName:  nil,
			expectedError: nil,
		},
		{
			name:          "Package line with comment",
			line:          "// package main",
			expectedName:  nil,
			expectedError: nil,
		},
		{
			name:          "Package line with extra spaces",
			line:          "    package main   ",
			expectedName:  utils.ToPointer("main"),
			expectedError: nil,
		},
		{
			name:          "Package line with multiple spaces",
			line:          "package    main",
			expectedName:  utils.ToPointer("main"),
			expectedError: nil,
		},
		{
			name:          "Package line with tabs",
			line:          "\tpackage\tmain",
			expectedName:  utils.ToPointer("main"),
			expectedError: nil,
		},
		{
			name:          "Package line with tabs and spaces",
			line:          "\t package\t main",
			expectedName:  utils.ToPointer("main"),
			expectedError: nil,
		},
		{
			name:          "Package line with multiple words",
			line:          "package main package",
			expectedName:  nil,
			expectedError: ErrInvalidPackageName,
		},
		{
			name:          "Package line followed by comment",
			line:          "package main // This is a comment",
			expectedName:  utils.ToPointer("main"),
			expectedError: nil,
		},
		{
			name:          "Package with no name",
			line:          "package",
			expectedName:  nil,
			expectedError: ErrInvalidPackageName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()

			result, err := extractPackageName(tc.line)

			assert.Equal(tt, tc.expectedName, result)
			if tc.expectedError != nil {
				assert.EqualError(tt, err, tc.expectedError.Error())
			} else {
				assert.NoError(tt, err)
			}
		})
	}
}

func Test_extractWireGenLocation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		filePath      string
		fileContent   string
		expectedInfo  *models.WireGenLocation
		expectedError error
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
				PackageName:   "wire",
				DirectoryPath: "internal/wire",
			},
			expectedError: nil,
		},
		{
			name:     "Valid file with wireinject multiple blank line",
			filePath: "internal/wire/wire.go",
			fileContent: `
			//go:build wireinject
			// +build wireinject



			package wire
			`,
			expectedInfo: &models.WireGenLocation{
				PackageName:   "wire",
				DirectoryPath: "internal/wire",
			},
			expectedError: nil,
		},
		{
			name:          "No wireinject in file",
			filePath:      "internal/wire/not_wire_file.go",
			fileContent:   "package notwire",
			expectedInfo:  nil,
			expectedError: nil,
		},
		{
			name:          "No package in file, only wireinject",
			filePath:      "internal/wire/not_wire_file.go",
			fileContent:   "//go:build wireinject",
			expectedInfo:  nil,
			expectedError: errors.New("no package in wiregen file: internal/wire/not_wire_file.go"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()

			wireGenLocation, err := extractWireGenLocation(tc.filePath, tc.fileContent)

			assert.Equal(tt, tc.expectedInfo, wireGenLocation)
			if tc.expectedError != nil {
				assert.EqualError(tt, tc.expectedError, err.Error())
			} else {
				assert.NoError(tt, err)
			}
		})
	}
}
