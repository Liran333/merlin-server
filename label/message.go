package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/webhook"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	"github.com/openmerlin/merlin-server/label/app"
	"github.com/openmerlin/merlin-server/label/utils"
	"github.com/openmerlin/merlin-server/models/domain"
	"github.com/openmerlin/merlin-server/models/messageapp"
)

const (
	msgHeaderUUID      = "X-Gitea-Event-UUID"
	msgHeaderUserAgent = "User-Agent"
	msgHeaderEventType = "X-Gitea-Event"
)

type MessageServer struct {
	userAgent     string
	labelService  app.LabelService
	modelsService messageapp.ModelAppService
}

func NewMessageServer(l app.LabelService, m messageapp.ModelAppService, ua string) *MessageServer {
	return &MessageServer{
		labelService:  l,
		modelsService: m,
		userAgent:     ua,
	}
}

func (m *MessageServer) handle(payload []byte, header map[string]string) error {
	eventType, err := m.parseRequest(header)
	if err != nil {
		return fmt.Errorf("invalid msg, err:%s", err.Error())
	}

	if eventType != webhook.HookEventPush {
		return nil
	}

	var p structs.PushPayload
	if err = json.Unmarshal(payload, &p); err != nil {
		return err
	}

	labels, ok, err := m.labelService.GetLabels(&p)
	if !ok {
		return err
	}

	index, err := m.toModelIndex(&p)
	if err != nil {
		return err
	}

	return m.modelsService.ResetLabels(&index, &labels)
}

func (m *MessageServer) parseRequest(header map[string]string) (eventType webhook.HookEventType, err error) {
	if header == nil {
		err = errors.New("no header")

		return
	}

	if header[msgHeaderUserAgent] != m.userAgent {
		err = errors.New("unknown " + msgHeaderUserAgent)

		return
	}

	if eventType = webhook.HookEventType(header[msgHeaderEventType]); eventType == "" {
		err = errors.New("missing " + msgHeaderEventType)

		return
	}

	if header[msgHeaderUUID] == "" {
		err = errors.New("missing " + msgHeaderUUID)
	}

	return
}

func (m *MessageServer) toModelIndex(p *structs.PushPayload) (index domain.ModelIndex, err error) {
	org, repo := utils.GetOrgRepo(p.Repo)
	owner, err := primitive.NewAccount(org)
	if err != nil {
		return
	}

	name, err := primitive.NewMSDName(repo)
	if err != nil {
		return
	}

	return domain.ModelIndex{
		Owner: owner,
		Name:  name,
	}, nil
}
