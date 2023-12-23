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

type FileCommitInfo struct {
	Message string    `json:"message"`
	Create  time.Time `json:"created"`
}

type ListFileInfo struct {
	Name       string         `json:"name"`
	Path       string         `json:"path"`
	Type       string         `json:"type"`
	Size       int64          `json:"size"`
	IsLfs      bool           `json:"isLfs"`
	URL        string         `json:"url"`
	FileCommit FileCommitInfo `json:"message"`
}

type FileAuthor struct {
	Name string `json:"name"`
}

// TODO i want to use DetailFileInfo inherit BaseInfo, But error in initial struct.

type DetailFileInfo struct {
	Name       string         `json:"name"`
	Path       string         `json:"path"`
	Type       string         `json:"type"`
	Size       int64          `json:"size"`
	IsLfs      bool           `json:"isLfs"`
	URL        string         `json:"url"`
	FileCommit FileCommitInfo `json:"message"`
	FileAuthor FileAuthor     `json:"author"`
}

type DownLoadFileInfo struct {
	IsLfs  bool   `json:"isLfs"`
	Stream []byte `json:"stream"`
	URL    string `json:"download_link"`
}
