package primitive

import (
	"errors"

	"github.com/openmerlin/merlin-server/utils"
)

type CommentContent interface {
	CommentContent() string
}

func NewCommentContent(v string) (CommentContent, error) {
	if utils.StrLen(v) > cfg.MaxContentLength ||
		utils.StrLen(utils.XSSEscapeString(v)) > cfg.MaxContentLength {
		return nil, errors.New("title is too long")
	}

	return commentContent(v), nil
}

func CreateCommentContent(v string) CommentContent {
	return commentContent(v)
}

type commentContent string

func (c commentContent) CommentContent() string {
	return string(c)
}
