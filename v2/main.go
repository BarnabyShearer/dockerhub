// API client for hub.docker.com
package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// APIs default base URL.
const BaseURLV2 = "https://hub.docker.com/v2"

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type Client struct {
	BaseURL    string
	auth       Auth
	HTTPClient *http.Client
}

// Create the API client, providing the authentication.
func NewClient(username string, password string) *Client {
	return &Client{
		BaseURL: BaseURLV2,
		auth: Auth{
			Username: username,
			Password: password,
		},
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type Repository struct {
	User            string `json:"user,omitempty"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	Description     string `json:"description"`
	Private         bool   `json:"is_private"`
	FullDescription string `json:"full_description,omitempty"`
}

type CreatePersonalAccessToken struct {
	TokenLabel string   `json:"token_label"`
	Scopes     []string `json:"scopes"`
}

type CreatePersonalAccessTokenResponse struct {
	UUID       string   `json:"uuid"`
	Token      string   `json:"token"`
	TokenLabel string   `json:"token_label"`
	Scopes     []string `json:"scopes"`
}

type Group struct {
    Id          int    `json:"id,omitempty"`
    Name        string `json:"name"`
	Description string `json:"description"`
}

type GroupMember struct {
    Username    string `json:"username"`
    FullName    string `json:"full_name"`
}

type RepositoryGroup struct {
    GroupId     int `json:"group_id"`
    GroupId2    int `json:"groupid"`
    GroupName   string `json:"group_name"`
    GroupName2  string `json:"groupname"`
    Permission  string `json:"permission"`
}

// Associate members to group:
// GET /v2/orgs/organisation_name/groups/NAME/members/
// RESPONSE: [{"username": "username", "full_name": "John Smith"}]

// POST /v2/orgs/organisation_name/groups/NAME/members/
// BODY: {"member":"username"}

// DELETE /v2/orgs/organisation_name/groups/test/members/username/

// Group
//------

// Create a group / team in organization
// POST https://hub.docker.com/v2/orgs/organisation_name/groups/
// BODY: {"name":"test"}
// RESPONSE: {"id":123,"name":"test","description":""}
func (c *Client) CreateGroup(ctx context.Context, organisation string, createGroup Group) (Group, error) {
	group := Group{}
	createGroupJson, err := json.Marshal(createGroup)
	if err != nil {
		return group, err
	}
	err = c.sendRequest(ctx, "POST", fmt.Sprintf("/orgs/%s/groups", organisation), createGroupJson, &group)
	if err != nil {
		return group, err
	}
	return group, err
}

// Edit a group / team in organization
// PATCH /v2/orgs/organisation_name/groups/group_name/
// BODY: {"name":"test","description":"x"}
// RESPONSE: {"id": 123, "name": "test", "description": "x"}
func (c *Client) UpdateGroup(ctx context.Context, organisation string, id string, updateGroup Group) (Group, error) {
	group := Group{}
	updateGroupJson, err := json.Marshal(updateGroup)
	if err != nil {
		return group, err
	}
	err = c.sendRequest(ctx, "PATCH", fmt.Sprintf("/orgs/%s/groups/%s/", organisation, id), updateGroupJson, &group)
	if err != nil {
		return group, err
	}
	return group, err
}

// Read a group / team in organization
// GET /v2/orgs/organization_name/groups/group_name/
// RESPONSE: {"id": 123, "name": "test", "description": "x"}
func (c *Client) GetGroup(ctx context.Context, organisation string, id string) (Group, error) {
	group := Group{}
    err := c.sendRequest(ctx, "GET", fmt.Sprintf("/orgs/%s/groups/%s/", organisation, id), nil, &group)
	return group, err
}

// Delete a group in organization
// DELETE /v2/orgs/organization_name/groups/group_name/
func (c *Client) DeleteGroup(ctx context.Context, organisation string, id string) error {
	return c.sendRequest(ctx, "DELETE", fmt.Sprintf("/orgs/%s/groups/%s/", organisation, id), nil, nil)
}

// Group-Repository Association
//-----------------------------

// Create a repository --> group association in organization
// POST https://hub.docker.com/v2/repositories/organisation_name/example-fixture-loader/groups/
// BODY: {"group_id":123,"groupid":123,"group_name":"example","groupname":"example","permission":"write"}
// permission: "read" / "write" / "admin"
func (c *Client) CreateRepositoryGroup(ctx context.Context, repository string, createRepositoryGroup RepositoryGroup) (RepositoryGroup, error) {
	repository_group := RepositoryGroup{}
	createRepositoryGroupJson, err := json.Marshal(createRepositoryGroup)
	if err != nil {
		return repository_group, err
	}
	err = c.sendRequest(ctx, "POST", fmt.Sprintf("/repositories/%s/groups/", repository), createRepositoryGroupJson, &repository_group)
	if err != nil {
		return repository_group, err
	}
	return repository_group, err
}

// Edit a repository --> group association in organization
// PATCH /v2/repositories/organisation_name/example-fixture-loader/groups/123/
// BODY: {"group_id":123,"group_name":"example","groupname":"example","permission":"admin"}
func (c *Client) UpdateRepositoryGroup(ctx context.Context, repository string, id string, updateRepositoryGroup RepositoryGroup) (RepositoryGroup, error) {
	repository_group := RepositoryGroup{}
	updateRepositoryGroupJson, err := json.Marshal(updateRepositoryGroup)
	if err != nil {
		return repository_group, err
	}
	err = c.sendRequest(ctx, "PATCH", fmt.Sprintf("/repositories/%s/groups/%s/", repository, id), updateRepositoryGroupJson, &repository_group)
	if err != nil {
		return repository_group, err
	}
	return repository_group, err
}

// Get a repository --> group association in organization
// GET /v2/repositories/organisation_name/example-fixture-loader/groups/123/
func (c *Client) GetRepositoryGroup(ctx context.Context, repository string, id string) (RepositoryGroup, error) {
	repository_group := RepositoryGroup{}
    err := c.sendRequest(ctx, "GET", fmt.Sprintf("/repositories/%s/groups/%s/", repository, id), nil, &repository_group)
	return repository_group, err
}

// Delete a repository --> group association in organization
// DELETE /v2/repositories/organisation_name/example-fixture-loader/groups/123/
func (c *Client) DeleteRepositoryGroup(ctx context.Context, repository string, id string) error {
	return c.sendRequest(ctx, "DELETE", fmt.Sprintf("/repositories/%s/groups/%s/", repository, id), nil, nil)
}

// Repository
//-----------
func (c *Client) CreateRepository(ctx context.Context, createRepository Repository) (Repository, error) {
	repository := Repository{}
	createRepositoryJson, err := json.Marshal(createRepository)
	if err != nil {
		return repository, err
	}
	err = c.sendRequest(ctx, "POST", "/repositories/", createRepositoryJson, &repository)
	if err != nil {
		return repository, err
	}
	return repository, err
}

func (c *Client) UpdateRepository(ctx context.Context, id string, updateRepository Repository) error {
	updateRepositoryJson, err := json.Marshal(updateRepository)
	if err != nil {
		return err
	}
	return c.sendRequest(ctx, "PATCH", fmt.Sprintf("/repositories/%s/", id), updateRepositoryJson, nil)
}

func (c *Client) GetRepository(ctx context.Context, id string) (Repository, error) {
	repository := Repository{}
	err := c.sendRequest(ctx, "GET", fmt.Sprintf("/repositories/%s/", id), nil, &repository)
	return repository, err
}

func (c *Client) DeleteRepository(ctx context.Context, id string) error {
	return c.sendRequest(ctx, "DELETE", fmt.Sprintf("/repositories/%s/", id), nil, nil)
}

// Personal Access Token
//----------------------
func (c *Client) CreatePersonalAccessToken(ctx context.Context, createPersonalAccessToken CreatePersonalAccessToken) (CreatePersonalAccessTokenResponse, error) {
	personalAccessToken := CreatePersonalAccessTokenResponse{}
	createRepositoryJson, err := json.Marshal(createPersonalAccessToken)
	if err != nil {
		return personalAccessToken, err
	}
	err = c.sendRequest(ctx, "POST", "/access-tokens", createRepositoryJson, &personalAccessToken)
	if err != nil {
		return personalAccessToken, err
	}
	return personalAccessToken, err
}

// Note: Returned token will always be blank.
func (c *Client) GetPersonalAccessToken(ctx context.Context, uuid string) (CreatePersonalAccessTokenResponse, error) {
	personalAccessToken := CreatePersonalAccessTokenResponse{}
	err := c.sendRequest(ctx, "GET", fmt.Sprintf("/access-tokens/%s", uuid), nil, &personalAccessToken)
	return personalAccessToken, err
}

func (c *Client) DeletePersonalAccessToken(ctx context.Context, uuid string) error {
	return c.sendRequest(ctx, "DELETE", fmt.Sprintf("/access-tokens/%s", uuid), nil, nil)
}

// Helpers
//--------
func (c *Client) sendRequest(ctx context.Context, method string, url string, body []byte, result interface{}) error {

	authJson, err := json.Marshal(c.auth)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/login/", c.BaseURL), bytes.NewBuffer(authJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	req = req.WithContext(ctx)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(bodyBytes))
	}
	token := Token{}
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return err
	}

	req, err = http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", token.Token))

	req = req.WithContext(ctx)

	res, err = c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(bodyBytes))
	}

	if result != nil {
		if err = json.NewDecoder(res.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}
