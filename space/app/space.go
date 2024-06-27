/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package app

import (
	"context"
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	coderepoapp "github.com/openmerlin/merlin-server/coderepo/app"
	commonapp "github.com/openmerlin/merlin-server/common/app"
	"github.com/openmerlin/merlin-server/common/domain/allerror"
	"github.com/openmerlin/merlin-server/common/domain/primitive"
	commonrepo "github.com/openmerlin/merlin-server/common/domain/repository"
	computilityapp "github.com/openmerlin/merlin-server/computility/app"
	computilitydomain "github.com/openmerlin/merlin-server/computility/domain"
	orgapp "github.com/openmerlin/merlin-server/organization/app"
	orgrepo "github.com/openmerlin/merlin-server/organization/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain"
	"github.com/openmerlin/merlin-server/space/domain/email"
	"github.com/openmerlin/merlin-server/space/domain/message"
	"github.com/openmerlin/merlin-server/space/domain/obs"
	"github.com/openmerlin/merlin-server/space/domain/repository"
	"github.com/openmerlin/merlin-server/space/domain/securestorage"
	spaceappApp "github.com/openmerlin/merlin-server/spaceapp/app"
	spaceappRepository "github.com/openmerlin/merlin-server/spaceapp/domain/repository"
	userapp "github.com/openmerlin/merlin-server/user/app"
	"github.com/openmerlin/merlin-server/utils"
)

const (
	variableTypeName = "variable"
	secretTypeName   = "secret"
)

func newSpaceNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceNotFound, "not found", err)
}

func newModelNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "not found", err)
}

func newSpaceSecretNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceSecretNotFound, "not found", err)
}

func newSpaceSecretCountExceeded(err error) error {
	return allerror.NewCountExceeded("space secret count exceed", err)
}

func newSpaceVariableNotFound(err error) error {
	return allerror.NewNotFound(allerror.ErrorCodeSpaceVariableNotFound, "not found", err)
}

func newSpaceVariableCountExceeded(err error) error {
	return allerror.NewCountExceeded("space variable count exceed", err)
}

// Permission is an interface for checking permissions.
type Permission interface {
	Check(primitive.Account, primitive.Account, primitive.ObjType, primitive.Action) error
}

// SpaceAppService is an interface for space application services.
type SpaceAppService interface {
	Create(context.Context, primitive.Account, *CmdToCreateSpace) (string, error)
	Delete(context.Context, primitive.Account, primitive.Identity) (string, error)
	Update(context.Context, primitive.Account, primitive.Identity, *CmdToUpdateSpace) (string, error)
	Disable(context.Context, primitive.Account, primitive.Identity, *CmdToDisableSpace) (string, error)
	GetByName(context.Context, primitive.Account, *domain.SpaceIndex) (SpaceDTO, error)
	List(context.Context, primitive.Account, *CmdToListSpaces) (SpacesDTO, error)
	AddLike(primitive.Identity) error
	DeleteLike(primitive.Identity) error
	Recommend(context.Context, primitive.Account) []repository.SpaceSummary
	Boutique(context.Context, primitive.Account) []repository.SpaceSummary
	setSpacesStatus(ctx context.Context, spacesDTO []repository.SpaceSummary) []repository.SpaceSummary
	SendSpaceReportEmail(primitive.Account, *CmdToReportDatasetEmail) error
	UploadCover(*CmdToUploadCover) (SpaceCoverDTO, error)
}

// NewSpaceAppService creates a new instance of SpaceAppService.
func NewSpaceAppService(
	permission commonapp.ResourcePermissionAppService,
	msgAdapter message.SpaceMessage,
	codeRepoApp coderepoapp.CodeRepoAppService,
	spaceappRepository spaceappRepository.Repository,
	variableAdapter repository.SpaceVariableRepositoryAdapter,
	secretAdapter repository.SpaceSecretRepositoryAdapter,
	secureStorageAdapter securestorage.SpaceSecureManager,
	repoAdapter repository.SpaceRepositoryAdapter,
	npuGatekeeper orgapp.PrivilegeOrg,
	member orgrepo.OrgMember,
	disableOrg orgapp.PrivilegeOrg,
	computilityApp computilityapp.ComputilityInternalAppService,
	spaceappApp spaceappApp.SpaceappAppService,
	user userapp.UserService,
	obs obs.ObsService,
	email email.Email,
) SpaceAppService {
	return &spaceAppService{
		permission:           permission,
		msgAdapter:           msgAdapter,
		codeRepoApp:          codeRepoApp,
		spaceappRepository:   spaceappRepository,
		variableAdapter:      variableAdapter,
		secretAdapter:        secretAdapter,
		secureStorageAdapter: secureStorageAdapter,
		repoAdapter:          repoAdapter,
		npuGatekeeper:        npuGatekeeper,
		member:               member,
		disableOrg:           disableOrg,
		computilityApp:       computilityApp,
		spaceappApp:          spaceappApp,
		user:                 user,
		obs:                  obs,
		email:                email,
	}
}

type spaceAppService struct {
	permission           commonapp.ResourcePermissionAppService
	msgAdapter           message.SpaceMessage
	codeRepoApp          coderepoapp.CodeRepoAppService
	spaceappRepository   spaceappRepository.Repository
	variableAdapter      repository.SpaceVariableRepositoryAdapter
	secretAdapter        repository.SpaceSecretRepositoryAdapter
	secureStorageAdapter securestorage.SpaceSecureManager
	repoAdapter          repository.SpaceRepositoryAdapter
	npuGatekeeper        orgapp.PrivilegeOrg
	member               orgrepo.OrgMember
	disableOrg           orgapp.PrivilegeOrg
	computilityApp       computilityapp.ComputilityInternalAppService
	spaceappApp          spaceappApp.SpaceappAppService
	user                 userapp.UserService
	obs                  obs.ObsService
	email                email.Email
}

// Create creates a new space with the given command and returns the ID of the created space.
func (s *spaceAppService) Create(ctx context.Context, user primitive.Account, cmd *CmdToCreateSpace) (string, error) {
	err := s.permission.CanCreate(ctx, user, cmd.Owner, primitive.ObjTypeSpace)
	if err != nil {
		return "", xerrors.Errorf("failed to create space: %w", err)
	}

	now := utils.Now()
	space := cmd.toSpace()

	id := primitive.CreateIdentity(primitive.GetId())
	hdType := space.GetComputeType()
	count := space.GetQuotaCount()
	compCmd := computilityapp.CmdToUserQuotaUpdate{
		Index: computilitydomain.ComputilityAccountRecordIndex{
			UserName:    user,
			ComputeType: hdType,
			SpaceId:     id,
		},
		QuotaCount: count,
	}

	if err := s.spaceCountCheck(ctx, cmd.Owner); err != nil {
		return "", err
	}

	coderepo, err := s.codeRepoApp.Create(ctx, user, &cmd.CmdToCreateRepo)
	if err != nil {
		return "", err
	}

	err = s.computilityApp.UserQuotaConsume(compCmd)
	if err != nil {
		logrus.Errorf("space create error | call api for quota consume failed | user:%s ,err: %s", user, err)

		return "", err
	}

	space.UpdatedAt = now
	space.CodeRepo = coderepo
	space.CreatedAt = now
	space.NoApplicationFile = true
	space.Exception = primitive.ExceptionNoApplicationFile

	if cmd.Hardware.IsNpu() {
		space.CompPowerAllocated = true
		space.Labels.HardwareType = "npu"
	} else {
		space.Labels.HardwareType = "cpu"
	}

	if err = s.repoAdapter.Add(&space); err != nil {
		err = xerrors.Errorf("space create failed | release user:%s quota | err: %w", user, err)
		logrus.Error(err)

		ierr := s.computilityApp.UserQuotaRelease(compCmd)
		if ierr != nil {
			return "", xerrors.Errorf("release user:%s quota failed after add space failed: %w", user, ierr)
		}

		if err = s.codeRepoApp.Delete(space.RepoIndex()); err != nil {
			logrus.Errorf("delete user:%s space repo:%v failed after add space failed: %v",
				user, space.RepoIndex(), ierr)

			return "", err
		}

		return "", err
	}

	if err = s.computilityApp.SpaceCreateSupply(computilityapp.CmdToSupplyRecord{
		Index: computilitydomain.ComputilityAccountRecordIndex{
			UserName:    user,
			ComputeType: hdType,
			SpaceId:     id,
		},
		QuotaCount: count,
		NewSpaceId: space.Id,
	}); err != nil {
		logrus.Errorf("add space id supplyment failed | user: %s, err: %s", user, err)

		_, err = s.Delete(ctx, user, space.Id)
		if err != nil {
			logrus.Errorf("delete space after add space id supplyment failed | user: %s, err: %s", user, err)

			return "", xerrors.Errorf("add space id supplyment failed: %w", err)
		}

		err = s.computilityApp.UserQuotaRelease(compCmd)
		if err != nil {
			logrus.Errorf("release quota after add space id supplyment failed | user: %s, err: %s", user, err)

			return "", xerrors.Errorf("add space id supplyment failed: %w", err)
		}

		return "", xerrors.Errorf("add space id supplyment failed: %w", err)
	}

	e := domain.NewSpaceCreatedEvent(&space)
	if err1 := s.msgAdapter.SendSpaceCreatedEvent(&e); err1 != nil {
		err1 = xerrors.Errorf("failed to send space created event, space id: %s err: %s", space.Id.Identity(), err1)
		logrus.Error(err1)
	}

	return space.Id.Identity(), nil
}

// Delete deletes the space with the given space ID and returns the action performed.
func (s *spaceAppService) Delete(
	ctx context.Context, user primitive.Account, spaceId primitive.Identity) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	action = fmt.Sprintf(
		"delete space of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), space.Name.MSDName(),
	)

	notFound, err := commonapp.CanDeleteOrNotFound(ctx, user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(fmt.Errorf("%s not found", spaceId.Identity()))

		return
	}
	if !s.codeRepoApp.IsNotFound(space.Id) {
		if err = s.codeRepoApp.Delete(space.RepoIndex()); err != nil {
			return
		}
	}
	// del space app
	if err = s.spaceappRepository.DeleteBySpaceId(space.Id); err != nil {
		return
	}

	// del space variable secret
	if err = s.delSpaceVariableSecret(space.Id); err != nil {
		return
	}

	if err = s.repoAdapter.Delete(space.Id); err != nil {
		return
	}

	e := domain.NewSpaceDeletedEvent(user, &space)
	if err1 := s.msgAdapter.SendSpaceDeletedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space deleted event, space id:%s", spaceId.Identity())
	}

	if space.Hardware.IsNpu() && space.CompPowerAllocated {
		logrus.Infof("release quota after user:%s npu space:%s delete", user, spaceId.Identity())

		c := computilityapp.CmdToUserQuotaUpdate{
			Index: computilitydomain.ComputilityAccountRecordIndex{
				UserName:    space.CreatedBy,
				ComputeType: space.GetComputeType(),
				SpaceId:     space.Id,
			},
			QuotaCount: space.GetQuotaCount(),
		}

		err = s.computilityApp.UserQuotaRelease(c)
		if err != nil {
			logrus.Errorf("failed to release user:%s quota after space:%s delete: %s",
				user.Account(), spaceId.Identity(), err)

			return "", nil
		}
	}

	return
}

func (s *spaceAppService) delSpaceVariableSecret(spaceId primitive.Identity) error {
	spaceVariableSecretList, err := s.variableAdapter.ListVariableSecret(spaceId.Identity())
	if err != nil {
		return err
	}
	for _, envSecret := range spaceVariableSecretList {
		envSecretId, err := primitive.NewIdentity(envSecret.Id)
		if err != nil {
			logrus.Errorf("failed to get envSecretId, err:%s", err)
			continue
		}
		if envSecret.Type == variableTypeName {
			variable, err := s.variableAdapter.FindVariableById(envSecretId)
			if err != nil {
				logrus.Errorf("failed to get variable, err:%s", err)
				continue
			}
			if err = s.secureStorageAdapter.DeleteSpaceEnvSecret(
				variable.GetVariablePath(), variable.Name.ENVName()); err != nil {
				logrus.Errorf("failed to delete variable, err:%s", err)
				continue
			}
			if err = s.variableAdapter.DeleteVariable(envSecretId); err != nil {
				logrus.Errorf("failed to delete variable db, err:%s", err)
				continue
			}
		}
		if envSecret.Type == secretTypeName {
			secret, err := s.secretAdapter.FindSecretById(envSecretId)
			if err != nil {
				logrus.Errorf("failed to get secret, err:%s", err)
				continue
			}
			if err = s.secureStorageAdapter.DeleteSpaceEnvSecret(
				secret.GetSecretPath(), secret.Name.ENVName()); err != nil {
				logrus.Errorf("failed to delete secret, err:%s", err)
				continue
			}
			if err = s.secretAdapter.DeleteSecret(envSecretId); err != nil {
				logrus.Errorf("failed to delete secret db, err:%s", err)
				continue
			}
		}
	}
	return nil
}

// Update updates the space with the given space ID using the provided command and returns the action performed.
func (s *spaceAppService) Update(
	ctx context.Context, user primitive.Account, spaceId primitive.Identity, cmd *CmdToUpdateSpace,
) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"update space of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), space.Name.MSDName(),
	)

	notFound, err := commonapp.CanUpdateOrNotFound(ctx, user, &space, s.permission)
	if err != nil {
		return
	}
	if notFound {
		err = newSpaceNotFound(fmt.Errorf("%s not found", spaceId.Identity()))

		return
	}

	if space.IsDisable() {
		err = allerror.NewResourceDisabled(allerror.ErrorCodeResourceDisabled,
			"resource was disabled, cant be modified.",
			fmt.Errorf("cant change resource to public"))
		return
	}

	oldVisibility := space.Visibility.Visibility()
	isPrivateToPublic := space.IsPrivate() && cmd.Visibility.IsPublic()

	b, err := s.codeRepoApp.Update(&space.CodeRepo, &cmd.CmdToUpdateRepo)
	if err != nil {
		return
	}

	b1 := cmd.toSpace(&space)
	if !b && !b1 {
		return
	}

	if err = s.repoAdapter.Save(&space); err != nil {
		return
	}

	if oldVisibility != space.Visibility.Visibility() {
		if _, err := s.spaceappApp.WakeupSpaceApp(ctx, user,
			&domain.SpaceIndex{Name: space.Name, Owner: space.Owner}); err != nil {
			logrus.Errorf("failed to walk up space app, space id:%s", spaceId.Identity())
		}
	}

	e := domain.NewSpaceUpdatedEvent(domain.SpaceUpdateEventParam{
		IsPriToPub:    isPrivateToPublic,
		Space:         &space,
		User:          user,
		OldVisibility: oldVisibility,
	})
	if err1 := s.msgAdapter.SendSpaceUpdatedEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space updated event, space id:%s", spaceId.Identity())
	}

	return
}

// Disable disable the space with the given space ID using the provided command and returns the action performed.
func (s *spaceAppService) Disable(
	ctx context.Context, user primitive.Account, spaceId primitive.Identity, cmd *CmdToDisableSpace,
) (action string, err error) {
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return
	}

	action = fmt.Sprintf(
		"disable space of %s:%s/%s",
		spaceId.Identity(), space.Owner.Account(), space.Name.MSDName(),
	)

	err = s.canDisable(ctx, user)
	if err != nil {
		return
	}

	if space.IsDisable() {
		logrus.Errorf("space %s already been disabled", space.Name.MSDName())
		err = allerror.NewResourceDisabled(
			allerror.ErrorCodeResourceAlreadyDisabled,
			"already been disabled",
			fmt.Errorf("already been disabled"))
		return
	}

	// del space app
	_, err = s.spaceappRepository.FindBySpaceId(ctx, space.Id)
	if err != nil && !commonrepo.IsErrorResourceNotExists(err) {
		logrus.Errorf("get space app by id %v failed, err:%v", space.Id, err)
		return
	} else if err == nil {
		if err = s.spaceappRepository.DeleteBySpaceId(space.Id); err != nil {
			logrus.Errorf("delete space app by id %v failed, err:%v", space.Id, err)
			return
		}
	}

	if space.Hardware.IsNpu() && space.CompPowerAllocated {
		logrus.Infof("release quota after npu space:%s delete", spaceId.Identity())

		c := computilityapp.CmdToUserQuotaUpdate{
			Index: computilitydomain.ComputilityAccountRecordIndex{
				UserName:    space.CreatedBy,
				ComputeType: space.GetComputeType(),
				SpaceId:     space.Id,
			},
			QuotaCount: space.GetQuotaCount(),
		}

		err = s.computilityApp.UserQuotaRelease(c)
		if err != nil {
			logrus.Errorf("failed to release user:%s quota after space:%s delete: %s",
				user.Account(), spaceId.Identity(), err)

			return
		}

		space.CompPowerAllocated = false
	}

	cmdRepo := coderepoapp.CmdToUpdateRepo{
		Visibility: primitive.VisibilityPrivate,
	}
	_, err = s.codeRepoApp.Update(&space.CodeRepo, &cmdRepo)
	if err != nil {
		return
	}

	cmd.toSpace(&space)

	if err = s.repoAdapter.Save(&space); err != nil {
		return
	}

	e := domain.NewSpaceDisableEvent(user, &space)
	if err1 := s.msgAdapter.SendSpaceDisableEvent(&e); err1 != nil {
		logrus.Errorf("failed to send space diabale event, space id:%s", spaceId.Identity())
	}

	logrus.Infof("send space diabale event success, space id:%s", spaceId.Identity())

	return
}

func (s *spaceAppService) canDisable(ctx context.Context, user primitive.Account) error {
	if s.disableOrg != nil {
		if err := s.disableOrg.Contains(ctx, user); err != nil {
			logrus.Errorf("user:%s cant disable space err:%s", user.Account(), err)
			return allerror.NewNoPermission("no permission", fmt.Errorf("cant disable"))
		}
	} else {
		logrus.Errorf("do not config disable org, no permit to disable")
		return allerror.NewNoPermission("no permission", fmt.Errorf("cant disable"))
	}

	return nil
}

// GetByName retrieves a space by its name and returns the corresponding SpaceDTO.
func (s *spaceAppService) GetByName(
	ctx context.Context, user primitive.Account, index *domain.SpaceIndex) (SpaceDTO, error) {
	var dto SpaceDTO

	space, err := s.repoAdapter.FindByName(index)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = newSpaceNotFound(err)
		}

		return dto, err
	}

	if err := s.permission.CanRead(ctx, user, &space); err != nil {
		if allerror.IsNoPermission(err) {
			err = newSpaceNotFound(err)
		}

		return dto, err
	}

	return toSpaceDTO(&space), nil
}

// List retrieves a list of spaces based on the provided command parameters and returns the corresponding SpacesDTO.
func (s *spaceAppService) List(ctx context.Context, user primitive.Account, cmd *CmdToListSpaces) (
	SpacesDTO, error,
) {
	if user == nil {
		cmd.Visibility = primitive.VisibilityPublic
	} else {
		if cmd.Owner == nil {
			// It can list the private spaces of user,
			// but it maybe no need to do it.
			cmd.Visibility = primitive.VisibilityPublic
		} else {
			if user != cmd.Owner {
				err := s.permission.CanListOrgResource(
					ctx, user, cmd.Owner, primitive.ObjTypeSpace,
				)
				if err != nil {
					cmd.Visibility = primitive.VisibilityPublic
				}
			}
		}
	}

	spaceLists, total, err := s.repoAdapter.List(cmd, user, s.member)

	spaceLists = s.setSpacesStatus(ctx, spaceLists)

	return SpacesDTO{
		Total:  total,
		Spaces: spaceLists,
	}, err
}

// DeleteById is an example for restful API.
func (s *spaceAppService) DeleteById(user primitive.Account, spaceId string) error {
	// get space by space id
	// check if user can delete it
	// delete it
	return nil
}

func (s *spaceAppService) spaceCountCheck(ctx context.Context, owner primitive.Account) error {
	cmdToList := CmdToListSpaces{
		Owner: owner,
	}

	total, err := s.repoAdapter.Count(&cmdToList)
	if err != nil {
		return allerror.NewCommonRespError("failed to count spaces", xerrors.Errorf("failed to count spaces: %w", err))
	}

	if s.user.IsOrganization(ctx, owner) && total >= config.MaxCountPerOrg {
		return allerror.NewCountExceeded("space count exceed",
			xerrors.Errorf("space count(now:%d max:%d) exceed", total, config.MaxCountPerOrg))
	}

	if !s.user.IsOrganization(ctx, owner) && total >= config.MaxCountPerUser {
		return allerror.NewCountExceeded("space count exceed",
			xerrors.Errorf("space count(now:%d max:%d) exceed", total, config.MaxCountPerUser))
	}

	return nil
}

func (s *spaceAppService) AddLike(spaceId primitive.Identity) error {
	// Retrieve the code repository information.
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		return err
	}

	// Only proceed if the repository is public.
	isPublic := space.IsPublic()
	if !isPublic {
		return nil
	}

	if err := s.repoAdapter.AddLike(space); err != nil {
		return err
	}
	return nil
}

func (s *spaceAppService) DeleteLike(spaceId primitive.Identity) error {
	// Retrieve the code repository information.
	space, err := s.repoAdapter.FindById(spaceId)
	if err != nil {
		return err
	}

	// Only proceed if the repository is public.
	isPublic := space.IsPublic()
	if !isPublic {
		return nil
	}

	if err := s.repoAdapter.DeleteLike(space); err != nil {
		return err
	}
	return nil
}

func (s *spaceAppService) Recommend(ctx context.Context, user primitive.Account) []repository.SpaceSummary {
	var spacesSummary []repository.SpaceSummary

	if len(config.RecommendSpaces) == 0 {
		logrus.Errorf("missing recommend spaces config")
		return spacesSummary
	}

	indexs := make([]domain.SpaceIndex, 0, len(config.RecommendSpaces))
	for _, v := range config.RecommendSpaces {
		idx := domain.SpaceIndex{
			Name:  primitive.CreateMSDName(v.Reponame),
			Owner: primitive.CreateAccount(v.Owner),
		}
		indexs = append(indexs, idx)
	}

	for _, index := range indexs {
		idx := index
		dto, err := s.GetByName(ctx, user, &idx)
		if err != nil {
			logrus.Errorf("failed to get recommend space by name:%s err:%s", idx.Name.MSDName(), err)
			continue
		}

		spaceSummary := toSpaceSummary(&dto)
		spacesSummary = append(spacesSummary, spaceSummary)
	}

	s.setSpacesStatus(ctx, spacesSummary)

	return spacesSummary
}

func (s *spaceAppService) Boutique(ctx context.Context, user primitive.Account) []repository.SpaceSummary {
	var spacesSummary []repository.SpaceSummary

	if len(config.BoutiqueSpaces) == 0 {
		logrus.Errorf("missing boutique spaces config")
		return spacesSummary
	}

	indexs := make([]domain.SpaceIndex, 0, len(config.BoutiqueSpaces))
	for _, v := range config.BoutiqueSpaces {
		idx := domain.SpaceIndex{
			Name:  primitive.CreateMSDName(v.Reponame),
			Owner: primitive.CreateAccount(v.Owner),
		}
		indexs = append(indexs, idx)
	}

	for _, index := range indexs {
		idx := index
		dto, err := s.GetByName(ctx, user, &idx)
		if err != nil {
			logrus.Errorf("failed to get boutique space by name:%s err:%s", idx.Name.MSDName(), err)
			continue
		}

		spaceSummary := toSpaceSummary(&dto)
		spacesSummary = append(spacesSummary, spaceSummary)
	}

	s.setSpacesStatus(ctx, spacesSummary)

	return spacesSummary
}

func (s *spaceAppService) setSpacesStatus(
	ctx context.Context, spacesSummary []repository.SpaceSummary) []repository.SpaceSummary {
	for key, spaceSummary := range spacesSummary {
		if spaceSummary.Exception != "" {
			spacesSummary[key].Status = spaceSummary.Exception
			continue
		}

		spaceId, err := primitive.NewIdentity(spaceSummary.Id)
		if err != nil {
			continue
		}
		app, err := s.spaceappRepository.FindBySpaceId(ctx, spaceId)
		if err == nil {
			spacesSummary[key].Status = app.Status.AppStatus()
			continue
		}
		if spaceSummary.IsNpu && !spaceSummary.CompPowerAllocated {
			spacesSummary[key].Status = primitive.NoCompQuotaException
		}
	}

	return spacesSummary
}

func (s *spaceAppService) SendSpaceReportEmail(user primitive.Account, cmd *CmdToReportDatasetEmail) error {
	pattern, _ := regexp.Compile(config.RegexpRule)
	html := pattern.FindStringSubmatch(cmd.Msg)
	if len(html) > 0 {
		e := fmt.Errorf("contains illegal characters")
		err := allerror.NewNoPermission(e.Error(), e)
		return err
	}
	index := domain.SpaceIndex{
		Owner: cmd.Owner,
		Name:  cmd.SpaceName,
	}
	data, err := s.repoAdapter.FindByName(&index)
	if err != nil {
		err := allerror.NewNoPermission(err.Error(), err)
		return err
	}
	if data.Visibility == primitive.CreateVisibility("private") && data.Owner != user {
		e := fmt.Errorf("no permission")
		err := allerror.NewNoPermission(e.Error(), e)
		return err
	}
	safeMsg := utils.XSSEscapeString(cmd.Msg)
	url := fmt.Sprintf("%s/spaces/%s/%s", s.email.GetRootUrl(), data.Owner.Account(), data.Name)
	if err := s.email.Send(cmd.SpaceName.MSDName(), safeMsg, user.Account(), url); err != nil {
		return err
	}
	return nil
}

// UpdateCover upload cover and  updates the cover url of space.
func (s *spaceAppService) UploadCover(cmd *CmdToUploadCover) (SpaceCoverDTO, error) {
	cover := domain.CoverInfo{
		User:        cmd.User,
		FileName:    cmd.FileName,
		Bucket:      config.ObsBucket,
		Path:        config.ObsPath,
		CdnEndpoint: config.CdnEndpoint,
	}

	err := s.obs.CreateObject(cmd.Image, config.ObsBucket, cover.GetObsPath())
	if err != nil {
		err = allerror.NewResourceDisabled(allerror.ErrorCodeFileUploadFailed, "",
			xerrors.Errorf("failed to upload space cover: %w", err))

		return SpaceCoverDTO{}, err
	}

	return toSpaceCoverDTO(cover.GetCoverURL()), nil
}
