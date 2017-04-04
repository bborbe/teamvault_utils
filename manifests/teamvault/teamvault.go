package teamvault

import (
	"fmt"
	http_header "github.com/bborbe/http/header"
	"github.com/bborbe/http/rest"
	"github.com/bborbe/kubernetes_tools/manifests/model"
	"net/http"
)

type teamvaultPasswordProvider struct {
	url  model.TeamvaultUrl
	user model.TeamvaultUser
	pass model.TeamvaultPassword
	rest rest.Rest
}

func New(
	executeRequest func(req *http.Request) (resp *http.Response, err error),
	url model.TeamvaultUrl,
	user model.TeamvaultUser,
	pass model.TeamvaultPassword,
) *teamvaultPasswordProvider {
	t := new(teamvaultPasswordProvider)
	t.rest = rest.New(executeRequest)
	t.url = url
	t.user = user
	t.pass = pass
	return t
}

func (t *teamvaultPasswordProvider) Password(key model.TeamvaultKey) (model.TeamvaultPassword, error) {
	currentRevision, err := t.CurrentRevision(key)
	if err != nil {
		return "", err
	}
	var response struct {
		Password model.TeamvaultPassword `json:"password"`
	}
	if err := t.rest.Call(fmt.Sprintf("%sdata", currentRevision.String()), nil, http.MethodGet, nil, &response, t.createHeader()); err != nil {
		return "", err
	}
	return response.Password, nil
}

func (t *teamvaultPasswordProvider) CurrentRevision(key model.TeamvaultKey) (model.TeamvaultCurrentRevision, error) {
	var response struct {
		CurrentRevision model.TeamvaultCurrentRevision `json:"current_revision"`
	}
	if err := t.rest.Call(fmt.Sprintf("%s/api/secrets/%s/", t.url.String(), key.String()), nil, http.MethodGet, nil, &response, t.createHeader()); err != nil {
		return "", err
	}
	return response.CurrentRevision, nil
}

func (t *teamvaultPasswordProvider) createHeader() http.Header {
	header := make(http.Header)
	header.Add("Authorization", fmt.Sprintf("Basic %s", http_header.CreateAuthorizationToken(t.user.String(), t.pass.String())))
	header.Add("Content-Type", "application/json")
	return header
}