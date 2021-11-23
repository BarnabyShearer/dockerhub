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
	updateRepositoryJSON, err := json.Marshal(updateRepository)
	if err != nil {
		return err
	}
	return c.sendRequest(ctx, "PATCH", fmt.Sprintf("/repositories/%s/", id), updateRepositoryJSON, nil)
}

func (c *Client) GetRepository(ctx context.Context, id string) (Repository, error) {
	repository := Repository{}
	err := c.sendRequest(ctx, "GET", fmt.Sprintf("/repositories/%s/", id), nil, &repository)
	return repository, err
}

func (c *Client) DeleteRepository(ctx context.Context, id string) error {
	return c.sendRequest(ctx, "DELETE", fmt.Sprintf("/repositories/%s/", id), nil, nil)
}

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

// Returned token will always be blank.
func (c *Client) GetPersonalAccessToken(ctx context.Context, uuid string) (CreatePersonalAccessTokenResponse, error) {
	personalAccessToken := CreatePersonalAccessTokenResponse{}
	err := c.sendRequest(ctx, "GET", fmt.Sprintf("/access-tokens/%s", uuid), nil, &personalAccessToken)
	return personalAccessToken, err
}

func (c *Client) DeletePersonalAccessToken(ctx context.Context, uuid string) error {
	return c.sendRequest(ctx, "DELETE", fmt.Sprintf("/access-tokens/%s", uuid), nil, nil)
}

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
