package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/server-common-lib/utils"

	commonctl "github.com/openmerlin/merlin-server/common/controller"
	"github.com/openmerlin/merlin-server/other"
)

const (
	refererHeader     = "Referer"
	secondLevelDomain = "osinfra.cn"
)

func AddRouterForOtherController(rg *gin.RouterGroup, c *other.Analyse) {
	ctl := OtherController{
		cfg: c,
		cli: utils.NewHttpClient(3),
	}

	rg.GET("/v1/analytics/key", ctl.Get)
}

type OtherController struct {
	cfg *other.Analyse
	cli utils.HttpClient
}

// @Summary  Ayalyse key
// @Description  get analyse key
// @Tags     Other
// @Success  200   {object}  commonctl.ResponseData{data=tokenResponse,msg=string,code=string}
// @Router   /v1/analytics/key [get]
func (ctl *OtherController) Get(ctx *gin.Context) {
	u, err := url.Parse(ctx.GetHeader(refererHeader))
	if !strings.HasSuffix(u.Host, secondLevelDomain) || err != nil {
		commonctl.SendError(ctx, errors.New("illegal request"))

		return
	}

	if token, err := ctl.getToken(); err != nil {
		commonctl.SendError(ctx, err)
	} else {
		commonctl.SendRespOfGet(ctx, token)
	}
}

type tokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	UseJwt       int    `json:"useJwt"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ClientID    string `json:"client_id"`
	ExpiresIn   int    `json:"expires_in"`
}

func (ctl *OtherController) getToken() (t tokenResponse, err error) {
	r := tokenRequest{
		GrantType:    "client_credentials",
		ClientID:     ctl.cfg.ClientID,
		ClientSecret: ctl.cfg.ClientSecret,
		UseJwt:       1,
	}

	body, err := utils.JsonMarshal(&r)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, ctl.cfg.GetTokenUrl, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	_, err = ctl.cli.ForwardTo(req, &t)

	t.ClientID = ctl.cfg.ClientID

	return
}
