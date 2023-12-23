package primitive

type FilePath interface {
	FilePath() string
}

func NewCodeFilePath(v string) (FilePath, error) {
	// todo judge the length of path
	return codeFilePath(v), nil
}

type codeFilePath string

func (r codeFilePath) FilePath() string {
	return string(r)
}
