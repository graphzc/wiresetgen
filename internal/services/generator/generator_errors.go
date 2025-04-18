package generator

import "errors"

var (
	ErrIsNotProjectRoot   = errors.New("is not in project root directory")
	ErrInvalidGoModFile   = errors.New("invalid go.mod file")
	ErrInvalidPackageName = errors.New("invalid package name")
)
