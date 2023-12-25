package domain

import (
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"time"
)

type CodeRepoFile struct {
	Owner    primitive.Account
	Name     primitive.MSDName
	Ref      primitive.FileRef
	FilePath primitive.FilePath
}

func (c CodeRepoFile) NewCodeRepoFile(codeRepoFile *CodeRepoFile) CodeRepoFile {
	return CodeRepoFile{
		Owner:    codeRepoFile.Owner,
		Name:     codeRepoFile.Name,
		Ref:      codeRepoFile.Ref,
		FilePath: codeRepoFile.FilePath,
	}
}

func (c CodeRepoFile) NewCodeRepoFileByUpdatePath(codeRepoFile *CodeRepoFile, path string) (*CodeRepoFile, error) {
	filePath, err := primitive.NewCodeFilePath(path)

	if err != nil {
		return nil, err
	}

	newCodeRepoFile := &CodeRepoFile{
		Owner:    codeRepoFile.Owner,
		Name:     codeRepoFile.Name,
		Ref:      codeRepoFile.Ref,
		FilePath: filePath,
	}

	return newCodeRepoFile, nil

}

type FileCommit struct {
	Message string    `json:"message"`
	Create  time.Time `json:"created"`
}

type FileAuthor struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type FileInfo struct {
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	Type       string     `json:"type"`
	Size       int64      `json:"size"`
	IsLfs      bool       `json:"isLfs"`
	URL        string     `json:"url"`
	FileCommit FileCommit `json:"commit"`
}

type LastCommitInfo struct {
	FileCommit FileCommit `json:"commit"`
	FileAuthor FileAuthor `json:"author"`
}

type ListFileInfo struct {
	Tree           []FileInfo     `json:"tree"`
	LastCommitInfo LastCommitInfo `json:"last_commit"`
}

// TODO i want to use DetailFileInfo inherit BaseInfo, But error in initial struct.

type DetailFileInfo struct {
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	Type       string     `json:"type"`
	Size       int64      `json:"size"`
	IsLfs      bool       `json:"isLfs"`
	URL        string     `json:"url"`
	FileCommit FileCommit `json:"commit"`
	FileAuthor FileAuthor `json:"author"`
}

type DownLoadFileInfo struct {
	IsLfs  bool   `json:"isLfs"`
	Stream string `json:"stream"`
	URL    string `json:"download_link"`
}
