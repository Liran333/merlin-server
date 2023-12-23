package coderepofileadapter

import (
	"code.gitea.io/sdk/gitea"
	"errors"
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"net/http"
)

type codeRepoFileAdapter struct {
	client *gitea.Client
}

func NewCodeRepoFileAdapter(c *gitea.Client) *codeRepoFileAdapter {
	return &codeRepoFileAdapter{client: c}
}

func (adapter *codeRepoFileAdapter) getContent(codeRepoFile *domain.CodeRepoFile) ([]byte, error) {
	crl, resp, err := adapter.client.GetFile(
		codeRepoFile.Owner.Account(),
		codeRepoFile.Name.MSDName(),
		codeRepoFile.Ref.FileRef(),
		codeRepoFile.FilePath.FilePath(),
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return crl, nil

}

func (adapter *codeRepoFileAdapter) getGitAttribute(codeRepoFile *domain.CodeRepoFile) ([]string, *string, error) {
	filePath, err := primitive.NewCodeFilePath(gitAttributesFile)

	if err != nil {
		return nil, nil, err
	}

	gitAttributesFile := &domain.CodeRepoFile{
		Owner:    codeRepoFile.Owner,
		Name:     codeRepoFile.Name,
		Ref:      codeRepoFile.Ref,
		FilePath: filePath,
	}

	// todo if not exist, there have error about 404
	fileStream, err := adapter.getContent(gitAttributesFile)

	if err != nil {
		return nil, nil, err
	}

	fileContent := string(fileStream)

	matchStr := parseGitAttributesFile(fileContent)

	return matchStr, &fileContent, nil

}

func (adapter *codeRepoFileAdapter) List(codeRepoFile *domain.CodeRepoFile) ([]domain.ListFileInfo, error) {
	crl, resp, err := adapter.client.ListContents(
		codeRepoFile.Owner.Account(),
		codeRepoFile.Name.MSDName(),
		codeRepoFile.Ref.FileRef(),
		codeRepoFile.FilePath.FilePath(),
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	matchStr, content, IsErrGetAttr := adapter.getGitAttribute(codeRepoFile)

	codeRepoFiles := make([]domain.ListFileInfo, 0, len(crl))

	for _, c := range crl {

		fileInfo := domain.ListFileInfo{
			Name:  c.Name,
			Path:  c.Path,
			Type:  c.Type,
			Size:  c.Size,
			IsLfs: false,
		}

		if c.Type == fileType {
			fileInfo.URL = *c.DownloadURL

		}

		if c.Type == fileType && IsErrGetAttr == nil {
			fileInfo.IsLfs = checkLfs(matchStr, *content, c.Path)
		}

		if fileInfo.IsLfs {
			filePath, err := primitive.NewCodeFilePath(c.Path)
			if err == nil {

				lfsRepoFile := codeRepoFile.NewCodeRepoFile(codeRepoFile)
				lfsRepoFile.FilePath = filePath

				lfsUrl, err := getLfsUrl(&lfsRepoFile, *c.DownloadURL)
				if err == nil {
					fileInfo.URL = lfsUrl
				}

			}

		}

		opt := gitea.ListCommitOptions{SHA: codeRepoFile.Ref.FileRef(), Path: c.Path}
		opt.PageSize = 1
		opt.Page = 1

		commits, resp, err := adapter.client.ListRepoCommits(codeRepoFile.Owner.Account(), codeRepoFile.Name.MSDName(), opt)

		if err != nil || resp.StatusCode != http.StatusOK {
			codeRepoFiles = append(codeRepoFiles, fileInfo)
			continue
		}

		if len(commits) > 0 {
			message := domain.FileCommitInfo{
				Create:  commits[0].Created,
				Message: commits[0].RepoCommit.Message,
			}
			fileInfo.FileCommit = message
		}

		codeRepoFiles = append(codeRepoFiles, fileInfo)
	}

	return codeRepoFiles, nil
}

func (adapter *codeRepoFileAdapter) Get(codeRepoFile *domain.CodeRepoFile) (*domain.DetailFileInfo, error) {
	crl, resp, err := adapter.client.GetContents(
		codeRepoFile.Owner.Account(),
		codeRepoFile.Name.MSDName(),
		codeRepoFile.Ref.FileRef(),
		codeRepoFile.FilePath.FilePath(),
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	fileInfo := domain.DetailFileInfo{
		Name:  crl.Name,
		Path:  crl.Path,
		Type:  crl.Type,
		Size:  crl.Size,
		IsLfs: false,
	}

	matchStr, content, IsErrGetAttr := adapter.getGitAttribute(codeRepoFile)

	if crl.Type == fileType {
		fileInfo.URL = *crl.DownloadURL
	}

	if crl.Type == fileType && IsErrGetAttr == nil {
		fileInfo.IsLfs = checkLfs(matchStr, *content, crl.Path)
	}

	if fileInfo.IsLfs {
		filePath, err := primitive.NewCodeFilePath(fileInfo.Path)

		if err == nil {
			lfsRepoFile := codeRepoFile.NewCodeRepoFile(codeRepoFile)
			lfsRepoFile.FilePath = filePath

			lfsUrl, err := getLfsUrl(&lfsRepoFile, *crl.DownloadURL)
			if err == nil {
				fileInfo.URL = lfsUrl
			}
		}

	}

	opt := gitea.ListCommitOptions{SHA: codeRepoFile.Ref.FileRef(), Path: crl.Path}
	opt.PageSize = 1
	opt.Page = 1

	commits, resp, err := adapter.client.ListRepoCommits(codeRepoFile.Owner.Account(), codeRepoFile.Name.MSDName(), opt)

	if err == nil && resp.StatusCode == http.StatusOK && len(commits) > 0 {

		lastCommit := commits[0]

		fileInfo.FileCommit = domain.FileCommitInfo{
			Create:  lastCommit.Created,
			Message: lastCommit.RepoCommit.Message,
		}

		fileInfo.FileAuthor = domain.FileAuthor{
			Name: lastCommit.RepoCommit.Author.Name,
		}

	}

	return &fileInfo, nil
}

func (adapter *codeRepoFileAdapter) Download(codeRepoFile *domain.CodeRepoFile) (*domain.DownLoadFileInfo, error) {
	fileInfo, err := adapter.Get(codeRepoFile)

	if err != nil {
		return nil, err
	}

	downloadFileInfo := domain.DownLoadFileInfo{
		IsLfs: fileInfo.IsLfs,
		URL:   fileInfo.URL,
	}

	if !downloadFileInfo.IsLfs {

		stream, err := adapter.getContent(codeRepoFile)

		if err != nil {
			return nil, err
		}

		downloadFileInfo.Stream = string(stream)
	}

	return &downloadFileInfo, nil
}
