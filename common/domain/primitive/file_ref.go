package primitive

type FileRef interface {
	FileRef() string
}

func NewCodeFileRef(v string) (FileRef, error) {
	// todo judge the length of ref
	return codeFileRef(v), nil
}

type codeFileRef string

func (r codeFileRef) FileRef() string {
	return string(r)
}
