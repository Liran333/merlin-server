package emailimpl

import (
	"fmt"

	"github.com/openmerlin/merlin-server/discussion/domain/email"
	"github.com/openmerlin/merlin-server/utils"
)

type Email interface {
	Send(receiver []string, subject, content string) error
}

func NewEmailImpl(e Email, c *Config) *emailImpl {
	return &emailImpl{
		email: e,
		cfg:   c,
	}
}

type emailImpl struct {
	email Email
	cfg   *Config
}

func (impl *emailImpl) SendReportEmail(param email.ReportEmailParam) error {
	emailContent := impl.buildContent(param)
	return impl.email.Send(impl.cfg.ReportEmailReceiver, impl.cfg.ReportTitle, emailContent)
}

func (impl *emailImpl) buildContent(param email.ReportEmailParam) string {
	comment := utils.XSSEscapeString(param.Content.CommentContent())
	url := fmt.Sprintf("%s%ss/%s/%s/issues/detail/%d",
		impl.cfg.RootUrl,
		param.ResourceType,
		param.Index.Owner.Account(),
		param.Index.Name.MSDName(),
		param.IssueId,
	)

	template := `
<html>
<body>
<h3>评论所在评论区</h3>
<p>%s</p>
<h3>评论内容</h3>
<p>%s</p>
<h3>举报原因</h3>
<p>%s</p>
<p>%s</p>
<h3>举报用户</h3>
<p>%s</p>
</body>
</html>
`
	return fmt.Sprintf(template, url, comment, param.ReportType, param.ReportContent, param.User.Account())
}
