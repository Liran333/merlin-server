package coderepofileadapter

import (
	"code.gitea.io/sdk/gitea"
	"errors"
	"github.com/openmerlin/merlin-server/coderepo/domain"
	"net/http"
	"sort"
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
	gitAttributesFile, err := codeRepoFile.NewCodeRepoFileByUpdatePath(codeRepoFile, gitAttributesFile)

	if err != nil {
		return nil, nil, err
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

func (adapter *codeRepoFileAdapter) getLastCommit(codeRepoFile *domain.CodeRepoFile) (*gitea.Commit, error) {
	opt := gitea.ListCommitOptions{SHA: codeRepoFile.Ref.FileRef(), Path: codeRepoFile.FilePath.FilePath()}
	opt.PageSize = 1
	opt.Page = 1

	commits, resp, err := adapter.client.ListRepoCommits(codeRepoFile.Owner.Account(), codeRepoFile.Name.MSDName(), opt)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	if len(commits) == 0 {
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	return commits[0], nil

}

func (adapter *codeRepoFileAdapter) List(codeRepoFile *domain.CodeRepoFile) (*domain.ListFileInfo, error) {
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

	lastCommit, err := adapter.getLastCommit(codeRepoFile)

	if err != nil {
		return nil, err
	}

	listFileInfo := &domain.ListFileInfo{LastCommitInfo: domain.LastCommitInfo{
		FileCommit: domain.FileCommit{
			Message: lastCommit.RepoCommit.Message,
			Create:  lastCommit.Created,
		},
		FileAuthor: domain.FileAuthor{
			Name:      lastCommit.Author.UserName,
			AvatarURL: lastCommit.Author.AvatarURL,
		},
	}}

	matchStr, content, IsErrGetAttr := adapter.getGitAttribute(codeRepoFile)

	FileInfos := make([]domain.FileInfo, 0, len(crl))
	DirInfos := make([]domain.FileInfo, 0, len(crl))

	for _, c := range crl {
		fileInfo := domain.FileInfo{
			Name:  c.Name,
			Path:  c.Path,
			Type:  c.Type,
			Size:  c.Size,
			IsLfs: false,
		}

		if c.Type == fileType && IsErrGetAttr == nil {
			fileInfo.IsLfs = checkLfs(matchStr, *content, c.Path)
		}

		if c.Type == fileType {
			fileInfo.URL = *c.DownloadURL
		}

		curRepoFile, IsErrParsePath := codeRepoFile.NewCodeRepoFileByUpdatePath(codeRepoFile, c.Path)

		if fileInfo.IsLfs && IsErrParsePath == nil {
			lfsUrl, err := getLfsUrl(curRepoFile, *c.DownloadURL)

			if err == nil {
				fileInfo.URL = lfsUrl
			}

		}

		curLastCommit, err := adapter.getLastCommit(curRepoFile)

		if err == nil {
			fileInfo.FileCommit = domain.FileCommit{
				Create:  curLastCommit.Created,
				Message: curLastCommit.RepoCommit.Message,
			}
		}

		if c.Type == dirType {
			DirInfos = append(DirInfos, fileInfo)
			continue
		}

		FileInfos = append(FileInfos, fileInfo)

	}

	sort.Slice(DirInfos, func(i, j int) bool {
		return DirInfos[j].Name > DirInfos[i].Name
	})

	sort.Slice(FileInfos, func(i, j int) bool {
		return FileInfos[j].Name > FileInfos[i].Name
	})

	DirInfos = append(DirInfos, FileInfos...)

	listFileInfo.Tree = DirInfos

	return listFileInfo, nil

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

	if crl.Type == fileType && IsErrGetAttr == nil {
		fileInfo.IsLfs = checkLfs(matchStr, *content, crl.Path)
	}

	if crl.Type == fileType {
		fileInfo.URL = *crl.DownloadURL
	}

	if fileInfo.IsLfs {
		lfsUrl, err := getLfsUrl(codeRepoFile, *crl.DownloadURL)

		if err == nil {
			fileInfo.URL = lfsUrl
		}
	}

	lastCommit, err := adapter.getLastCommit(codeRepoFile)

	if err != nil {
		return nil, err
	}

	fileInfo.FileCommit = domain.FileCommit{
		Create:  lastCommit.Created,
		Message: lastCommit.RepoCommit.Message,
	}

	fileInfo.FileAuthor = domain.FileAuthor{
		Name:      lastCommit.Author.UserName,
		AvatarURL: lastCommit.Author.AvatarURL,
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
