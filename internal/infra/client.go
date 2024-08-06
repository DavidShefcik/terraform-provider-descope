package infra

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/descope/go-sdk/descope/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const PrincipalProjectID = "<principal>"

type Response struct {
	Entity string         `json:"entity"`
	ID     string         `json:"id"`
	Data   map[string]any `json:"data"`
}

type Client struct {
	projectID     string
	managementKey string
	baseURL       string

	apiClients map[string]*api.Client
	lock       sync.Mutex
}

func NewClient(projectID, managementKey, baseURL string) *Client {
	return &Client{
		projectID:     projectID,
		managementKey: managementKey,
		baseURL:       baseURL,
		apiClients:    map[string]*api.Client{},
	}
}

func (c *Client) Create(ctx context.Context, projectID, entity string, data map[string]any) (*Response, error) {
	httpBody := map[string]any{
		"entity": entity,
		"data":   data,
	}

	tflog.Info(ctx, "Starting CREATE request", map[string]any{"body": debugRequest(httpBody)})
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, "/v1/mgmt/infra", httpBody, nil, c.managementKey)
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished CREATE request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Read(ctx context.Context, projectID, entity, entityID string) (*Response, error) {
	if projectID == c.projectID || projectID == PrincipalProjectID {
		return nil, errors.New("principal project may not be read by the provider")
	}

	httpQuery := map[string]string{
		"entity": entity,
		"id":     entityID,
	}

	tflog.Info(ctx, "Starting READ request", map[string]any{"query": debugRequest(httpQuery)})
	httpRes, err := c.getAPIClient(projectID).DoGetRequest(ctx, "/v1/mgmt/infra", &api.HTTPRequest{QueryParams: httpQuery}, c.managementKey)
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished READ request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Update(ctx context.Context, projectID, entity, entityID string, data map[string]any) (*Response, error) {
	if projectID == c.projectID || projectID == PrincipalProjectID {
		return nil, errors.New("principal project may not be updated by the provider")
	}

	httpBody := map[string]any{
		"entity": entity,
		"id":     entityID,
		"data":   data,
	}

	tflog.Info(ctx, "Starting UPDATE request", map[string]any{"body": debugRequest(httpBody)})
	httpRes, err := c.getAPIClient(projectID).DoPutRequest(ctx, "/v1/mgmt/infra", httpBody, nil, c.managementKey)
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished UPDATE request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Delete(ctx context.Context, projectID, entity, entityID string) error {
	if projectID == c.projectID || projectID == PrincipalProjectID {
		return errors.New("principal project may not be deleted by the provider")
	}

	httpQuery := map[string]string{
		"entity": entity,
		"id":     entityID,
	}

	tflog.Info(ctx, "Starting DELETE request", map[string]any{"query": debugRequest(httpQuery)})
	httpRes, err := c.getAPIClient(projectID).DoDeleteRequest(ctx, "/v1/mgmt/infra", &api.HTTPRequest{QueryParams: httpQuery}, c.managementKey)
	if err != nil {
		return err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return err
	}

	tflog.Info(ctx, "Finished DELETE request")
	return nil
}

func (c *Client) getAPIClient(projectID string) *api.Client {
	if projectID == PrincipalProjectID {
		projectID = c.projectID
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	apiClient, ok := c.apiClients[projectID]
	if !ok {
		params := api.ClientParams{ProjectID: projectID, BaseURL: c.baseURL}
		apiClient = api.NewClient(params)
		c.apiClients[projectID] = apiClient
	}

	return apiClient
}
