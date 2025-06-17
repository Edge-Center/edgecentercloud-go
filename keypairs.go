package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	keypairsBasePathV1 = "/v1/keypairs"
	keypairsBasePathV2 = "/v2/keypairs"
)

const (
	keypairsShare = "share"
)

// KeyPairsService is an interface for creating and managing SSH keys with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/keypairs
type KeyPairsService interface {
	List(context.Context) ([]KeyPair, *Response, error)
	ListV2(ctx context.Context, opts *KeyPairsListOptionsV2) ([]KeyPairV2, *Response, error)
	Get(context.Context, string) (*KeyPair, *Response, error)
	GetV2(context.Context, string) (*KeyPairV2, *Response, error)
	Create(context.Context, *KeyPairCreateRequest) (*KeyPair, *Response, error)
	CreateV2(context.Context, *KeyPairCreateRequestV2) (*KeyPairV2, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	DeleteV2(context.Context, string) (*Response, error)
	Share(context.Context, string, *KeyPairShareRequest) (*KeyPair, *Response, error)
}

// KeyPairsServiceOp handles communication with Key Pairs (SSH keys) methods of the EdgecenterCloud API.
type KeyPairsServiceOp struct {
	client *Client
}

var _ KeyPairsService = &KeyPairsServiceOp{}

// KeyPair represents an EdgecenterCloud Key Pair.
type KeyPair struct {
	SSHKeyID        string `json:"sshkey_id"`
	PublicKey       string `json:"public_key"`
	PrivateKey      string `json:"private_key"`
	Fingerprint     string `json:"fingerprint"`
	SSHKeyName      string `json:"sshkey_name"`
	State           string `json:"state"`
	SharedInProject bool   `json:"shared_in_project"`
	CreatedAt       string `json:"created_at"`
	ProjectID       int    `json:"project_id"`
}

// KeyPairV2 represents an EdgecenterCloud Key Pair.
type KeyPairV2 struct {
	PublicKey       string `json:"public_key"`
	Fingerprint     string `json:"fingerprint"`
	CreatedAt       string `json:"created_at"`
	State           string `json:"state"`
	ProjectID       int    `json:"project_id"`
	SSHKeyID        string `json:"sshkey_id"`
	SharedInProject bool   `json:"shared_in_project"`
	PrivateKey      string `json:"private_key"`
	SSHKeyName      string `json:"sshkey_name"`
}

// KeyPairCreateRequest represents a request to create a Key Pair.
type KeyPairCreateRequest struct {
	SSHKeyName      string `json:"sshkey_name" required:"true"`
	PublicKey       string `json:"public_key,omitempty"`
	SharedInProject bool   `json:"shared_in_project"`
}

// KeyPairCreateRequestV2 represents a request to create a Key Pair.
type KeyPairCreateRequestV2 struct {
	ProjectID       int    `json:"project_id" required:"true"`
	SSHKeyName      string `json:"sshkey_name" required:"true"`
	PublicKey       string `json:"public_key,omitempty"`
	SharedInProject bool   `json:"shared_in_project,omitempty"`
}

// KeyPairShareRequest represents a request to share a Key Pair.
type KeyPairShareRequest struct {
	SharedInProject bool `json:"shared_in_project" required:"true"`
}

// keyPairsRoot represents a KeyPair root.
type keyPairsRoot struct {
	Count   int
	KeyPair []KeyPair `json:"results"`
}

// keyPairsRoot represents a KeyPair root.
type keyPairsRootV2 struct {
	Count   int
	KeyPair []KeyPairV2 `json:"results"`
}

// KeyPairsListOptionsV2 specifies the optional query parameters to List method.
type KeyPairsListOptionsV2 struct {
	ProjectID int `url:"project_id,omitempty"  validate:"omitempty"`
	UserID    int `url:"user_id,omitempty"  validate:"omitempty"`
}

// List get KeyPairs.
func (s *KeyPairsServiceOp) List(ctx context.Context) ([]KeyPair, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(keypairsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyPairsRoot)

	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.KeyPair, resp, err
}

// ListV2 get KeyPairs.
func (s *KeyPairsServiceOp) ListV2(ctx context.Context, opts *KeyPairsListOptionsV2) ([]KeyPairV2, *Response, error) {
	path, err := addOptions(keypairsBasePathV2, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyPairsRootV2)

	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.KeyPair, resp, err
}

// Get individual Key Pair.
func (s *KeyPairsServiceOp) Get(ctx context.Context, keypairID string) (*KeyPair, *Response, error) {
	if resp, err := isValidUUID(keypairID, "keypairID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(keypairsBasePathV1), keypairID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	keyPair := new(KeyPair)
	resp, err := s.client.Do(ctx, req, keyPair)
	if err != nil {
		return nil, resp, err
	}

	return keyPair, resp, err
}

// GetV2 individual Key Pair.
func (s *KeyPairsServiceOp) GetV2(ctx context.Context, keypairID string) (*KeyPairV2, *Response, error) {
	if resp, err := isValidUUID(keypairID, "keypairID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", keypairsBasePathV2, keypairID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	keyPair := new(KeyPairV2)

	resp, err := s.client.Do(ctx, req, keyPair)
	if err != nil {
		return nil, resp, err
	}

	return keyPair, resp, err
}

// Create a Key Pair.
func (s *KeyPairsServiceOp) Create(ctx context.Context, reqBody *KeyPairCreateRequest) (*KeyPair, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(keypairsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	keyPair := new(KeyPair)
	resp, err := s.client.Do(ctx, req, keyPair)
	if err != nil {
		return nil, resp, err
	}

	return keyPair, resp, err
}

func (s *KeyPairsServiceOp) CreateV2(ctx context.Context, reqBody *KeyPairCreateRequestV2) (*KeyPairV2, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, keypairsBasePathV2, reqBody)
	if err != nil {
		return nil, nil, err
	}

	keyPair := new(KeyPairV2)

	resp, err := s.client.Do(ctx, req, keyPair)
	if err != nil {
		return nil, resp, err
	}

	return keyPair, resp, err
}

// Delete the Key Pair.
func (s *KeyPairsServiceOp) Delete(ctx context.Context, keypairID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(keypairID, "keypairID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(keypairsBasePathV1), keypairID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

// DeleteV2 the Key Pair.
func (s *KeyPairsServiceOp) DeleteV2(ctx context.Context, keypairID string) (*Response, error) {
	if resp, err := isValidUUID(keypairID, "keypairID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s", keypairsBasePathV2, keypairID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Share a Key Pair to view for all users in project.
func (s *KeyPairsServiceOp) Share(ctx context.Context, keypairID string, reqBody *KeyPairShareRequest) (*KeyPair, *Response, error) {
	if resp, err := isValidUUID(keypairID, "keypairID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("shareRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(keypairsBasePathV1), keypairID, keypairsShare)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	keyPair := new(KeyPair)
	resp, err := s.client.Do(ctx, req, keyPair)
	if err != nil {
		return nil, resp, err
	}

	return keyPair, resp, err
}
