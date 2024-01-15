package git

import (
	"fmt"

	"github.com/openmerlin/go-sdk/gitea"

	"github.com/openmerlin/merlin-server/common/domain/primitive"
	common "github.com/openmerlin/merlin-server/common/infrastructure/gitea"
	org "github.com/openmerlin/merlin-server/organization/domain"
	"github.com/openmerlin/merlin-server/user/domain"
	"github.com/openmerlin/merlin-server/utils"
)

type BaseAuthClient struct {
	username string
	client   *gitea.Client
}

func NewBaseAuthClient(username, password string) (*BaseAuthClient, error) {
	client, err := common.NewClient(username, password)
	if err != nil {
		return nil, err
	}

	return &BaseAuthClient{
		username: username,
		client:   client,
	}, nil
}

type TokenInfo struct {
	Token     string
	CreatedAt int64
}

type TokenCreatedCmd struct {
	Name string
}

func (c *BaseAuthClient) CreateToken(cmd *domain.TokenCreatedCmd) (token domain.PlatformToken, err error) {
	if cmd == nil {
		err = fmt.Errorf("create token param is empty")
		return
	}

	if cmd.Account.Account() != c.username {
		err = fmt.Errorf("username mismatched, requested user: %s, client user: %s", cmd.Account.Account(), c.username)
		return
	}

	var perms []gitea.AccessTokenScope
	// hack, gitea go sdk do not contains perm we need
	for _, p := range domain.ToPerms(cmd.Permission) {
		perms = append(perms, gitea.AccessTokenScope(p))
	}
	// create token first
	t, _, err := c.client.CreateAccessToken(gitea.CreateAccessTokenOption{
		Name:   cmd.Name.Account(),
		Scopes: perms,
	})
	if err != nil {
		return
	}

	token.CreatedAt = utils.Now()
	token.Token = t.Token
	token.Name = cmd.Name
	token.Account = cmd.Account
	token.Permission = cmd.Permission
	token.Expire = cmd.Expire
	token.LastEight = t.TokenLastEight

	return
}

func (c *BaseAuthClient) DeleteToken(cmd *domain.TokenDeletedCmd) (err error) {
	if cmd == nil {
		return fmt.Errorf("nil cmd")
	}

	if cmd.Account.Account() != c.username {
		return fmt.Errorf("username mismatched, requested user: %s, client user: %s", cmd.Account.Account(), c.username)
	}
	_, err = c.client.DeleteAccessToken(cmd.Name.Account())

	return
}

func (c *BaseAuthClient) CreateOrg(cmd *org.Organization) (err error) {
	if cmd == nil {
		err = fmt.Errorf("nil cmd")
		return
	}

	if cmd.Owner == nil {
		err = fmt.Errorf("nil owner")
		return
	}

	if cmd.Name == nil {
		err = fmt.Errorf("nil name")
		return
	}

	if cmd.Owner.Account() != c.username {
		err = fmt.Errorf("username mismatched, requested user: %s, client user: %s", cmd.Owner, c.username)
		return
	}

	tmp, _, err := c.client.CreateOrg(gitea.CreateOrgOption{
		Name:                      cmd.Name.Account(),
		FullName:                  cmd.Fullname.MSDFullname(),
		Description:               cmd.Description.MSDDesc(),
		Website:                   cmd.Website,
		Visibility:                gitea.VisibleTypePublic,
		RepoAdminChangeTeamAccess: false,
	})
	// we also create write & read team
	if err != nil {
		return
	}

	teams, _, err := c.client.ListOrgTeams(cmd.Name.Account(), gitea.ListTeamsOptions{})
	if err != nil {
		_, _ = c.client.DeleteOrg(cmd.Name.Account())
		err = fmt.Errorf("failed to list org teams: %w", err)
		return
	}
	// first team must be owner team
	if len(teams) != 1 {
		_, _ = c.client.DeleteOrg(cmd.Name.Account())
		err = fmt.Errorf("invalid org team count: %d", len(teams))
		return
	}
	cmd.OwnerTeamId = teams[0].ID

	team, _, err := c.client.CreateTeam(cmd.Name.Account(), gitea.CreateTeamOption{
		Name:                    "contributor",
		Description:             "contributor team",
		Permission:              gitea.AccessModeRead,
		CanCreateOrgRepo:        true,
		IncludesAllRepositories: true,
		Units:                   []gitea.RepoUnitType{gitea.RepoUnitCode},
	})
	if err != nil {
		_, _ = c.client.DeleteOrg(cmd.Name.Account())
		return
	}
	cmd.ContributorTeamId = team.ID

	team, _, err = c.client.CreateTeam(cmd.Name.Account(), gitea.CreateTeamOption{
		Name:                    "write",
		Description:             "write team",
		Permission:              gitea.AccessModeWrite,
		CanCreateOrgRepo:        true,
		IncludesAllRepositories: true,
		Units:                   []gitea.RepoUnitType{gitea.RepoUnitCode},
	})
	if err != nil {
		_, _ = c.client.DeleteOrg(cmd.Name.Account())
		return
	}

	cmd.WriteTeamId = team.ID

	team, _, err = c.client.CreateTeam(cmd.Name.Account(), gitea.CreateTeamOption{
		Name:                    "read",
		Description:             "read team",
		Permission:              gitea.AccessModeRead,
		CanCreateOrgRepo:        false,
		IncludesAllRepositories: true,
		Units:                   []gitea.RepoUnitType{gitea.RepoUnitCode},
	})
	if err != nil {
		_, _ = c.client.DeleteOrg(cmd.Name.Account())
		return
	}

	cmd.ReadTeamId = team.ID
	cmd.PlatformId = tmp.ID

	return
}

func (c *BaseAuthClient) DeleteOrg(name primitive.Account) (err error) {
	repos, _, err := c.client.ListOrgRepos(name.Account(), gitea.ListOrgReposOptions{})
	if err != nil {
		err = fmt.Errorf("failed to list org repos: %w", err)
		return
	}
	if len(repos) != 0 {
		err = fmt.Errorf("org %s has repos, cannot delete", name)
		return
	}

	_, err = c.client.DeleteOrg(name.Account())

	return
}

func (c *BaseAuthClient) AddMember(o *org.Organization, member *org.OrgMember) (err error) {
	if o == nil {
		return fmt.Errorf("nil cmd")
	}
	if member == nil {
		return fmt.Errorf("nil member")
	}

	switch member.Role {
	case org.OrgRoleContributor:
		_, err = c.client.AddTeamMember(o.ContributorTeamId, member.Username.Account())
	case org.OrgRoleReader:
		_, err = c.client.AddTeamMember(o.ReadTeamId, member.Username.Account())
	case org.OrgRoleWriter:
		_, err = c.client.AddTeamMember(o.WriteTeamId, member.Username.Account())
	case org.OrgRoleAdmin:
		_, err = c.client.AddTeamMember(o.OwnerTeamId, member.Username.Account())
	default:
		return fmt.Errorf("member role %s is not supported", member.Role)
	}

	return err
}

func (c *BaseAuthClient) RemoveMember(o *org.Organization, member *org.OrgMember) (err error) {
	if o == nil {
		return fmt.Errorf("nil cmd")
	}

	if member == nil {
		return fmt.Errorf("nil member")
	}

	switch member.Role {
	case org.OrgRoleContributor:
		_, err = c.client.RemoveTeamMember(o.ContributorTeamId, member.Username.Account())
	case org.OrgRoleReader:
		_, err = c.client.RemoveTeamMember(o.ReadTeamId, member.Username.Account())
	case org.OrgRoleWriter:
		_, err = c.client.RemoveTeamMember(o.WriteTeamId, member.Username.Account())
	case org.OrgRoleAdmin:
		_, err = c.client.RemoveTeamMember(o.OwnerTeamId, member.Username.Account())
	default:
		return fmt.Errorf("member role %s is not supported", member.Role)
	}

	return err
}

func (c *BaseAuthClient) CanDelete(name primitive.Account) (can bool, err error) {
	repos, _, err := c.client.ListOrgRepos(name.Account(), gitea.ListOrgReposOptions{})
	if err != nil {
		return
	}
	if len(repos) != 0 {
		can = false
		return
	}

	can = true
	return
}

func (c *BaseAuthClient) EditMemberRole(o *org.Organization, orig org.OrgRole, now *org.OrgMember) (err error) {
	switch orig {
	case org.OrgRoleContributor:
		_, err = c.client.RemoveTeamMember(o.ContributorTeamId, now.Username.Account())
	case org.OrgRoleReader:
		_, err = c.client.RemoveTeamMember(o.ReadTeamId, now.Username.Account())
	case org.OrgRoleWriter:
		_, err = c.client.RemoveTeamMember(o.WriteTeamId, now.Username.Account())
	case org.OrgRoleAdmin:
		_, err = c.client.RemoveTeamMember(o.OwnerTeamId, now.Username.Account())
	default:
		return fmt.Errorf("member role %s is not supported", now.Role)
	}

	if err != nil {
		return fmt.Errorf("failed to remove team member when editing member role: %w", err)
	}

	switch now.Role {
	case org.OrgRoleContributor:
		_, err = c.client.AddTeamMember(o.ContributorTeamId, now.Username.Account())
	case org.OrgRoleReader:
		_, err = c.client.AddTeamMember(o.ReadTeamId, now.Username.Account())
	case org.OrgRoleWriter:
		_, err = c.client.AddTeamMember(o.WriteTeamId, now.Username.Account())
	case org.OrgRoleAdmin:
		_, err = c.client.AddTeamMember(o.OwnerTeamId, now.Username.Account())
	default:
		return fmt.Errorf("member role %s is not supported", now.Role)
	}

	if err != nil {
		return fmt.Errorf("failed to add team member when editing member role: %w", err)
	}

	return
}
