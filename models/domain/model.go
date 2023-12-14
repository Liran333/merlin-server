package domain

import (
	coderepo "github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
)

type Model struct {
	coderepo.CodeRepo

	Desc     primitive.MSDDesc
	Fullname primitive.MSDFullname

	Labels        []string
	Version       int
	CreatedAt     int64
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}
