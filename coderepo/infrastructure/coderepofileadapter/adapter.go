package coderepofileadapter

import (
	"errors"
	"net/http"
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/openmerlin/go-sdk/gitea"
	"github.com/openmerlin/merlin-server/coderepo/domain"
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
	crl, resp, err := adapter.client.ListCommitContents(
		codeRepoFile.Owner.Account(),
		codeRepoFile.Name.MSDName(),
		codeRepoFile.Ref.FileRef(),
		codeRepoFile.FilePath.FilePath(),
	)

	if err != nil {
		logrus.Errorf("ListCommitContents failed :%s", err.Error())

		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("ListCommitContents return invalid code : %d", resp.StatusCode)

		return nil, errors.New(resp.Status)
	}

	lastCommit, err := adapter.getLastCommit(codeRepoFile)

	if err != nil {
		logrus.Errorf("ListCommitContents getLastCommit failed :%s", err.Error())

		return nil, err
	}

	// if fork others directly, the Author will be nil
	fileAuthor := domain.FileAuthor{}
	if lastCommit.Author == nil {
		fileAuthor.Name = lastCommit.RepoCommit.Author.Name
		fileAuthor.AvatarURL = ``
	} else {
		fileAuthor.Name = lastCommit.Author.UserName
		fileAuthor.AvatarURL = lastCommit.Author.AvatarURL
	}

	listFileInfo := &domain.ListFileInfo{LastCommitInfo: domain.LastCommitInfo{
		FileCommit: domain.FileCommit{
			Message: lastCommit.RepoCommit.Message,
			Create:  lastCommit.Created,
		},
		FileAuthor: fileAuthor,
	}}

	FileInfos := make([]domain.FileInfo, 0, len(crl))
	DirInfos := make([]domain.FileInfo, 0, len(crl))

	for _, c := range crl {
		fileInfo := domain.FileInfo{
			Name:  c.Name,
			Path:  c.Path,
			Type:  c.Type,
			Size:  c.Size,
			IsLfs: c.IsLFS,
			FileCommit: domain.FileCommit{
				Message: c.LastCommitMessage,
				Create:  c.LastCommitCreate,
			},
		}

		if c.Type == fileType && c.DownloadURL != nil {

			fileInfo.URL = *c.DownloadURL

		}

		curRepoFile, IsErrParsePath := codeRepoFile.NewCodeRepoFileByUpdatePath(codeRepoFile, c.Path)

		if fileInfo.IsLfs && IsErrParsePath == nil {
			lfsUrl, err := getLfsUrl(curRepoFile, *c.DownloadURL)

			if err == nil {
				fileInfo.URL = lfsUrl
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
	crl, resp, err := adapter.client.GetCommitContents(
		codeRepoFile.Owner.Account(),
		codeRepoFile.Name.MSDName(),
		codeRepoFile.Ref.FileRef(),
		codeRepoFile.FilePath.FilePath(),
	)

	if err != nil {
		logrus.Errorf("GetContents failed :%s", err.Error())

		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("GetContents return invalid code : %d", resp.StatusCode)

		return nil, errors.New(resp.Status)
	}

	fileInfo := domain.DetailFileInfo{
		Name:  crl.Name,
		Path:  crl.Path,
		Type:  crl.Type,
		Size:  crl.Size,
		IsLfs: crl.IsLFS,
	}

	if crl.Type == fileType && crl.DownloadURL != nil {
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
		logrus.Errorf("getLastCommit failed :%s", err.Error())

		return nil, err
	}

	fileInfo.FileCommit = domain.FileCommit{
		Create:  lastCommit.Created,
		Message: lastCommit.RepoCommit.Message,
	}

	// if fork others directly, the Author will be nil
	fileAuthor := domain.FileAuthor{}
	if lastCommit.Author == nil {
		fileAuthor.Name = lastCommit.RepoCommit.Author.Name
		fileAuthor.AvatarURL = ``
	} else {
		fileAuthor.Name = lastCommit.Author.UserName
		fileAuthor.AvatarURL = lastCommit.Author.AvatarURL
	}

	fileInfo.FileAuthor = fileAuthor

	return &fileInfo, nil

}

func (adapter *codeRepoFileAdapter) Download(codeRepoFile *domain.CodeRepoFile) (*domain.DownLoadFileInfo, error) {
	fileInfo, err := adapter.Get(codeRepoFile)

	if err != nil {
		logrus.Errorf("Get failed :%s", err.Error())

		return nil, err
	}

	downloadFileInfo := domain.DownLoadFileInfo{
		IsLfs: fileInfo.IsLfs,
		URL:   fileInfo.URL,
	}

	if !downloadFileInfo.IsLfs {

		stream, err := adapter.getContent(codeRepoFile)

		if err != nil {
			logrus.Errorf("getContent failed :%s", err.Error())

			return nil, err
		}

		downloadFileInfo.Stream = string(stream)
	}

	return &downloadFileInfo, nil
}
