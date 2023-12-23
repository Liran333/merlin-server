package giteaimpl

import (
	"bytes"

	"code.gitea.io/sdk/gitea"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"k8s.io/apimachinery/pkg/util/sets"

	commongitea "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	localgitea "github.com/openmerlin/merlin-server/label/domain/gitea"
	"github.com/openmerlin/merlin-server/models/domain"
)

type iClient interface {
	GetFile(owner, repo, ref, filepath string, resolveLFS ...bool) ([]byte, *gitea.Response, error)
}

type giteaImpl struct {
	cli iClient
}

func NewGiteaImpl(cfg *commongitea.Config) *giteaImpl {
	cli, _ := gitea.NewClient(cfg.URL, gitea.SetToken(cfg.Token))

	return &giteaImpl{
		cli: cli,
	}
}

func (impl *giteaImpl) GetLabels(option *localgitea.Option) (label domain.ModelLabels, err error) {
	content, _, err := impl.cli.GetFile(option.Org, option.Repo, option.Ref, option.Path)
	if err != nil {
		return
	}

	meta, err := impl.parseTags(content)
	if err != nil {
		return
	}

	return domain.ModelLabels{
		Task:       meta.Task,
		Others:     sets.New(meta.Tags...),
		Frameworks: sets.New(meta.Frameworks...),
	}, nil
}

type MetaData struct {
	Task       string   `yaml:"task"`
	Tags       []string `yaml:"tags"`
	Frameworks []string `yaml:"frameworks"`
}

func (impl *giteaImpl) parseTags(content []byte) (meta MetaData, err error) {
	md := goldmark.New(
		goldmark.WithExtensions(&frontmatter.Extender{}),
	)

	ctx := parser.NewContext()
	var buf bytes.Buffer
	if err = md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return
	}

	if err = frontmatter.Get(ctx).Decode(&meta); err != nil {
		return
	}

	return meta, nil
}
