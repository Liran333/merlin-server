/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"time"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

// CodeRepoFile represents a code repository file with its owner, name, reference, and file path.
type CodeRepoFile struct {
	Owner    primitive.Account
	Name     primitive.MSDName
	Ref      primitive.FileRef
	FilePath primitive.FilePath
}

// NewCodeRepoFile creates a new CodeRepoFile instance with the given values.
func (c CodeRepoFile) NewCodeRepoFile(codeRepoFile *CodeRepoFile) CodeRepoFile {
	return CodeRepoFile{
		Owner:    codeRepoFile.Owner,
		Name:     codeRepoFile.Name,
		Ref:      codeRepoFile.Ref,
		FilePath: codeRepoFile.FilePath,
	}
}

// NewCodeRepoFileByUpdatePath creates a new CodeRepoFile instance by updating the file path.
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

// FileCommit represents a file commit with its message and creation time.
type FileCommit struct {
	Message string    `json:"message"`
	Create  time.Time `json:"created"`
}

// FileAuthor represents an author of a file with their name and avatar URL.
type FileAuthor struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// FileInfo represents information about a file including its
// name, path, type, size, LFS status, URL, and commit details.
type FileInfo struct {
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	Type       string     `json:"type"`
	Size       int64      `json:"size"`
	IsLfs      bool       `json:"isLfs"`
	URL        string     `json:"url"`
	FileCommit FileCommit `json:"commit"`
}

// LastCommitInfo represents the last commit information for a file including the commit details and author details.
type LastCommitInfo struct {
	FileCommit FileCommit `json:"commit"`
	FileAuthor FileAuthor `json:"author"`
}

// ListFileInfo represents a list of files with their tree structure and last commit information.
type ListFileInfo struct {
	Tree           []FileInfo     `json:"tree"`
	LastCommitInfo LastCommitInfo `json:"last_commit"`
}

// DetailFileInfo represents detailed information about a file including its
// name, path, type, size, LFS status, URL, commit details, and author details.
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

// DownLoadFileInfo represents download information for a file including its LFS status, stream, and download link.
type DownLoadFileInfo struct {
	IsLfs  bool   `json:"isLfs"`
	Stream string `json:"stream"`
	URL    string `json:"download_link"`
}
