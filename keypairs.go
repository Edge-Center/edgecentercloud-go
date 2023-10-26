package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	keypairsBasePathV1 = "/v1/keypairs"
)

// KeyPairsService is an interface for creating and managing SSH keys with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/keypairs
type KeyPairsService interface {
	List(context.Context) ([]KeyPair, *Response, error)
	Get(context.Context, string) (*KeyPair, *Response, error)
	Create(context.Context, *KeyPairCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
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

// KeyPairCreateRequest represents a request to create a Key Pair.
type KeyPairCreateRequest struct {
	SSHKeyName      string `json:"sshkey_name" required:"true"`
	PublicKey       string `json:"public_key"`
	SharedInProject bool   `json:"shared_in_project"`
}

// keyPairsRoot represents a KeyPair root.
type keyPairsRoot struct {
	Count   int
	KeyPair []KeyPair `json:"results"`
}

// List get KeyPairs.
func (s *KeyPairsServiceOp) List(ctx context.Context) ([]KeyPair, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addServicePath(keypairsBasePathV1)

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

// Get individual Key Pair.
func (s *KeyPairsServiceOp) Get(ctx context.Context, keypairID string) (*KeyPair, *Response, error) {
	if resp, err := isValidUUID(keypairID, "keypairID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(keypairsBasePathV1), keypairID)

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

// Create a Key Pair.
func (s *KeyPairsServiceOp) Create(ctx context.Context, createRequest *KeyPairCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addServicePath(keypairsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
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

// Delete the Key Pair.
func (s *KeyPairsServiceOp) Delete(ctx context.Context, keypairID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(keypairID, "keypairID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(keypairsBasePathV1), keypairID)

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
