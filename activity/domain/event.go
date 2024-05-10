package domain

import (
	"encoding/json"

	"github.com/openmerlin/merlin-server/coderepo/domain"
	"github.com/openmerlin/merlin-server/utils"
)

// likeCreatedEvent
type likeCreatedEvent struct {
	Time      int64  `json:"time"`
	Owner     string `json:"owner"`
	RepoId    string `json:"repo_id"`
	RepoType  string `json:"repo_type"`
	RepoName  string `json:"repo_name"`
	CreatedBy string `json:"created_by"`
}

func (e *likeCreatedEvent) Message() ([]byte, error) {
	return json.Marshal(e)
}

// NewLikeCreatedEvent return a modelCreatedEvent
func NewLikeCreatedEvent(m *domain.CodeRepo, t string) likeCreatedEvent {
	return likeCreatedEvent{
		Time:      utils.Now(),
		Owner:     m.Owner.Account(),
		RepoId:    m.Id.Identity(),
		RepoType:  t,
		RepoName:  m.Name.MSDName(),
		CreatedBy: m.CreatedBy.Account(),
	}
}
