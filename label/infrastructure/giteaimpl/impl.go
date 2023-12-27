package giteaimpl

import (
	"bytes"
	"errors"
	"fmt"

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
		err = fmt.Errorf("parse tags of [%s %s] err: %s", option.Org, option.Repo, err.Error())

		return
	}

	return domain.ModelLabels{
		Task:       meta.Task,
		Others:     sets.New(meta.Tags...),
		Frameworks: sets.New(meta.Frameworks),
	}, nil
}

type MetaData struct {
	Task       string   `yaml:"pipeline_tag"`
	Tags       []string `yaml:"tags"`
	Frameworks string   `yaml:"framework"`
}

func (impl *giteaImpl) parseTags(content []byte) (meta MetaData, err error) {
	if len(content) == 0 {
		err = errors.New("README.md is empty")

		return
	}

	md := goldmark.New(
		goldmark.WithExtensions(&frontmatter.Extender{}),
	)

	ctx := parser.NewContext()
	var buf bytes.Buffer
	if err = md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return
	}

	data := frontmatter.Get(ctx)
	if data == nil {
		err = errors.New("frontmatter is empty")

		return
	}

	if err = data.Decode(&meta); err != nil {
		return
	}

	return meta, nil
}
