package edgecloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	keypairsBasePathV1 = "/v1/keypairs"
)

// KeyPairsService is an interface for creating and managing SSH keys with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/keypairs
type KeyPairsService interface {
	Get(context.Context, string, *ServicePath) (*KeyPair, *Response, error)
	Create(context.Context, *KeyPairCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	Delete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
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

// keyPairRoot represents a Key Pair root.
type keyPairRoot struct {
	KeyPair *KeyPair      `json:"keypair"`
	Tasks   *TaskResponse `json:"tasks"`
}

// Get individual Key Pair.
func (s *KeyPairsServiceOp) Get(ctx context.Context, keypairID string, p *ServicePath) (*KeyPair, *Response, error) {
	if _, err := uuid.Parse(keypairID); err != nil {
		return nil, nil, NewArgError("keypairID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(keypairsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, keypairID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyPairRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.KeyPair, resp, err
}

// Create a Key Pair.
func (s *KeyPairsServiceOp) Create(ctx context.Context, createRequest *KeyPairCreateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(keypairsBasePathV1, p)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyPairRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// Delete the Key Pair.
func (s *KeyPairsServiceOp) Delete(ctx context.Context, keypairID string, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(keypairID); err != nil {
		return nil, nil, NewArgError("keypairID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(keypairsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, keypairID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyPairRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}
