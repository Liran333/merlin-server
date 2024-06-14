package domain

import (
	"time"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	discussionprimitive "github.com/openmerlin/merlin-server/discussion/domain/primitive"
)

type IssueComment struct {
	Id        int64
	Author    primitive.Account
	emojis    []Emoji
	IssueId   int64
	Content   discussionprimitive.CommentContent
	CreatedAt time.Time
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

func (c *IssueComment) UpdateContent(content discussionprimitive.CommentContent) {
	c.Content = content
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

type IssueCommentReport struct {
	User      primitive.Account
	Type      string
	Content   string
	CommentId int64
}
