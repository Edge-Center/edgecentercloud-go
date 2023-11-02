package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	secretsBasePathV1 = "/v1/secrets"
	secretsBasePathV2 = "/v2/secrets"
)

// SecretsService is an interface for creating and managing Secrets with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/secrets
type SecretsService interface {
	List(context.Context) ([]Secret, *Response, error)
	Create(context.Context, *SecretCreateRequest) (*TaskResponse, *Response, error)
	CreateV2(context.Context, *SecretCreateRequestV2) (*TaskResponse, *Response, error)
	Get(context.Context, string) (*Secret, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
}

// SecretsServiceOp handles communication with Secrets methods of the EdgecenterCloud API.
type SecretsServiceOp struct {
	client *Client
}

var _ SecretsService = &SecretsServiceOp{}

// Secret represents an EdgecenterCloud Secret.
type Secret struct {
	Expiration   string            `json:"expiration"`
	Algorithm    string            `json:"algorithm"`
	Name         string            `json:"name"`
	Mode         string            `json:"mode"`
	ID           string            `json:"id"`
	BitLength    int               `json:"bit_length"`
	Created      string            `json:"created"`
	Status       string            `json:"status"`
	SecretType   string            `json:"secret_type"`
	ContentTypes map[string]string `json:"content_types"`
}

type SecretCreateRequest struct {
	Expiration             string     `json:"expiration,omitempty"`
	PayloadContentType     string     `json:"payload_content_type" required:"true" validate:"required"`
	Algorithm              string     `json:"algorithm,omitempty"`
	Name                   string     `json:"name" required:"true" validate:"required"`
	Mode                   string     `json:"mode,omitempty"`
	BitLength              int        `json:"bit_length,omitempty"`
	PayloadContentEncoding string     `json:"payload_content_encoding" required:"true" validate:"required"`
	Payload                string     `json:"payload" required:"true" validate:"required"`
	SecretType             SecretType `json:"secret_type" required:"true" validate:"required"`
}

type SecretCreateRequestV2 struct {
	Expiration string `json:"expiration,omitempty"`
	Name       string `json:"name" required:"true" validate:"required"`
	Payload    string `json:"payload" required:"true" validate:"required"`
}

type Payload struct {
	CertificateChain string `json:"certificate_chain" required:"true" validate:"required"`
	PrivateKey       string `json:"private_key" required:"true" validate:"required"`
	Certificate      string `json:"certificate" required:"true" validate:"required"`
}

type SecretType string

const (
	SecretTypeSymmetric   SecretType = "symmetric"
	SecretTypePublic      SecretType = "public"
	SecretTypePrivate     SecretType = "private"
	SecretTypePassphrase  SecretType = "passphrase"
	SecretTypeCertificate SecretType = "certificate"
	SecretTypeOpaque      SecretType = "opaque"
)

// secretsRoot represents a Secrets root.
type secretsRoot struct {
	Count   int
	Secrets []Secret `json:"results"`
}

// List get secrets.
func (s *SecretsServiceOp) List(ctx context.Context) ([]Secret, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(secretsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(secretsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Secrets, resp, err
}

// Create a Secret.
func (s *SecretsServiceOp) Create(ctx context.Context, reqBody *SecretCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(secretsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
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

// CreateV2 a Secret V2.
func (s *SecretsServiceOp) CreateV2(ctx context.Context, reqBody *SecretCreateRequestV2) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(secretsBasePathV2)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
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

// Get a Secret.
func (s *SecretsServiceOp) Get(ctx context.Context, secretID string) (*Secret, *Response, error) {
	if resp, err := isValidUUID(secretID, "secretID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(secretsBasePathV1), secretID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	secret := new(Secret)
	resp, err := s.client.Do(ctx, req, secret)
	if err != nil {
		return nil, resp, err
	}

	return secret, resp, err
}

// Delete a Secret.
func (s *SecretsServiceOp) Delete(ctx context.Context, secretID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(secretID, "secretID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(secretsBasePathV1), secretID)

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
