package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
)

type IssueComment struct {
	Id             int64
	Author         primitive.Account
	emojis         []Emoji
	IssueId        int64
	Content        discussionprimitive.CommentContent
	CreatedAt      time.Time
	IsFirstComment bool
}

func NewFirstIssueComment(author primitive.Account, issueId int64, content discussionprimitive.CommentContent,
) IssueComment {
	return IssueComment{
		Author:         author,
		IssueId:        issueId,
		Content:        content,
		IsFirstComment: true,
	}
}

func NewIssueComment(author primitive.Account, issueId int64, content discussionprimitive.CommentContent,
) IssueComment {
	return IssueComment{
		Author:  author,
		IssueId: issueId,
		Content: content,
	}
}

func (c *IssueComment) AddEmoji(emojiType discussionprimitive.EmojiType, user primitive.Account) {
	for k := range c.emojis {
		if c.emojis[k].addUser(emojiType, user) {
			return
		}
	}

	c.emojis = append(c.emojis, Emoji{
		Type:  emojiType,
		Users: []primitive.Account{user},
	})
}

func (c *IssueComment) IsCommentOwner(user primitive.Account) bool {
	return c.Author == user
}

func (c *IssueComment) UpdateContent(user primitive.Account, content discussionprimitive.CommentContent) error {
	if !c.IsCommentOwner(user) {
		return allerror.NewNoPermission("no permission", errors.New("not comment owner"))
	}

	c.Content = content

	return nil
}

func (c *IssueComment) IsFirstCommentOfIssue() bool {
	return c.IsFirstComment
}

type Emoji struct {
	Type  discussionprimitive.EmojiType
	Users []primitive.Account
}

func (e *Emoji) addUser(t discussionprimitive.EmojiType, user primitive.Account) bool {
	if e.Type == t {
		e.Users = append(e.Users, user)

		return true
	}

	return false
}

type updateCommentCountEvent struct {
	IssueId              int64 `json:"issue_id"`
	IncreaseCommentCount int64 `json:"increase_comment_count"`
}

func (u updateCommentCountEvent) Message() ([]byte, error) {
	return json.Marshal(u)
}

func NewUpdateCommentCountEvent(issueId, count int64) updateCommentCountEvent {
	return updateCommentCountEvent{
		IssueId:              issueId,
		IncreaseCommentCount: count,
	}
}
