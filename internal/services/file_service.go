package services

type FileService interface {
	ReadFile(filePath string) ([]byte, error)
}

type fileServiceImpl struct {
}

func NewFileService() FileService {
	return &fileServiceImpl{}
}

func (f *fileServiceImpl) ReadFile(filePath string) ([]byte, error) {
	panic("unimplemented")
}
