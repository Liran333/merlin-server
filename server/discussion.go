package server

import (
	"github.com/gin-gonic/gin"

	"github.com/openmerlin/merlin-server/coderepo/infrastructure/resourceadapterimpl"
	"github.com/openmerlin/merlin-server/common/infrastructure/email"
	"github.com/openmerlin/merlin-server/common/infrastructure/postgresql"
	"github.com/openmerlin/merlin-server/config"
	"github.com/openmerlin/merlin-server/datasets/infrastructure/datasetrepositoryadapter"
	"github.com/openmerlin/merlin-server/discussion/app"
	"github.com/openmerlin/merlin-server/discussion/controller"
	"github.com/openmerlin/merlin-server/discussion/infrastructure/emailimpl"
	"github.com/openmerlin/merlin-server/discussion/infrastructure/messageimpl"
	"github.com/openmerlin/merlin-server/discussion/infrastructure/repositoryimpl"
	"github.com/openmerlin/merlin-server/models/infrastructure/modelrepositoryadapter"
	"github.com/openmerlin/merlin-server/space/infrastructure/spacerepositoryadapter"
)

func initDiscussion(cfg *config.Config, services *allServices) {
	issueRepoImpl := repositoryimpl.NewIssueImpl(postgresql.DAO(cfg.Discussion.Tables.Issue))
	commentRepoImpl := repositoryimpl.NewIssueCommentImpl(postgresql.DAO(cfg.Discussion.Tables.IssueComment))
	resourceImpl := resourceadapterimpl.NewResourceAdapterImpl(
		modelrepositoryadapter.ModelAdapter(),
		datasetrepositoryadapter.DatasetAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
	)

	services.discussionIssue = app.NewIssueService(
		resourceImpl,
		services.permissionApp,
		issueRepoImpl,
		issueRepoImpl,
		commentRepoImpl,
	)

	services.discussionComment = app.NewCommentService(
		resourceImpl,
		services.permissionApp,
		issueRepoImpl,
		commentRepoImpl,
		messageimpl.NewMessageImpl(cfg.Discussion.Topics),
		emailimpl.NewEmailImpl(email.GetEmailInst(), &cfg.Discussion.Report),
	)

	services.discussion = app.NewDiscussionService(
		resourceImpl,
		services.permissionApp,
		modelrepositoryadapter.ModelAdapter(),
		spacerepositoryadapter.SpaceAdapter(),
		datasetrepositoryadapter.DatasetAdapter(),
	)
}

func setRouterOfDiscussionWeb(rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForDiscussionWebController(
		rg,
		services.userMiddleWare,
		services.operationLog,
		services.discussionIssue,
		services.discussionComment,
		services.discussion,
	)
}

func setRouterOfDiscussionInternal(cfg *config.Config, rg *gin.RouterGroup, services *allServices) {
	controller.AddRouterForDiscussionInternalController(
		rg,
		services.userMiddleWare,
		app.NewIssueInternalService(
			repositoryimpl.NewIssueImpl(
				postgresql.DAO(cfg.Discussion.Tables.Issue),
			),
			repositoryimpl.NewIssueCommentImpl(
				postgresql.DAO(cfg.Discussion.Tables.IssueComment),
			),
		))
}
