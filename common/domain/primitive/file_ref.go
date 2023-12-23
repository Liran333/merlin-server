package primitive

type FileRef interface {
	FileRef() string
}

func NewCodeFileRef(v string) (FileRef, error) {
	// todo judge the length of ref
	return codeFileRef(v), nil
}

func InitCodeFileRef() FileRef {
	return codeFileRef("main")
}

type codeFileRef string

func (r codeFileRef) FileRef() string {
	return string(r)
}
