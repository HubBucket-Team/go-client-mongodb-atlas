package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const projectAPIKeysPath = "groups/%s/apiKeys"

//ProjectAPIKeysService is an interface for interfacing with the APIKeys
// endpoints of the MongoDB Atlas API.
//See more: https://docs.atlas.mongodb.com/reference/api/apiKeys/#organization-api-keys-on-projects-endpoints
type ProjectAPIKeysService interface {
	List(context.Context, string, *ListOptions) ([]APIKey, *Response, error)
	Create(context.Context, string, *APIKeyInput) (*APIKey, *Response, error)
	Assign(context.Context, string, string) (*Response, error)
	Unassign(context.Context, string, string) (*Response, error)
}

//ProjectAPIKeysOp handles communication with the APIKey related methods
// of the MongoDB Atlas API
type ProjectAPIKeysOp struct {
	client *Client
}

var _ ProjectAPIKeysService = &ProjectAPIKeysOp{}

//List all API-KEY in the organization associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/projectApiKeys/get-all-apiKeys-in-one-project/
func (s *ProjectAPIKeysOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]APIKey, *Response, error) {
	path := fmt.Sprintf(projectAPIKeysPath, groupID)

	//Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(apiKeysResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//Create an API Key by the {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/projectApiKeys/create-one-apiKey-in-one-project/
func (s *ProjectAPIKeysOp) Create(ctx context.Context, groupID string, createRequest *APIKeyInput) (*APIKey, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(projectAPIKeysPath, groupID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(APIKey)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Assign an API-KEY related to {GROUP-ID} to a the project with {API-KEY-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/projectApiKeys/assign-one-org-apiKey-to-one-project/
func (s *ProjectAPIKeysOp) Assign(ctx context.Context, groupID string, keyID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("apiKeyID", "must be set")
	}

	if keyID == "" {
		return nil, NewArgError("keyID", "must be set")
	}

	basePath := fmt.Sprintf(projectAPIKeysPath, groupID)

	path := fmt.Sprintf("%s/%s", basePath, keyID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

//Unassign an API-KEY related to {GROUP-ID} to a the project with {API-KEY-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/projectApiKeys/delete-one-apiKey-in-one-project/
func (s *ProjectAPIKeysOp) Unassign(ctx context.Context, groupID string, keyID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("apiKeyID", "must be set")
	}

	if keyID == "" {
		return nil, NewArgError("keyID", "must be set")
	}

	basePath := fmt.Sprintf(projectAPIKeysPath, groupID)

	path := fmt.Sprintf("%s/%s", basePath, keyID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
