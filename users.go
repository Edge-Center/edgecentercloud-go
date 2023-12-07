package edgecloud

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

const (
	usersBasePathV1 = "/v1/users"
)

const (
	usersRoles       = "roles"
	usersAssignments = "assignments"
)

// UsersService is an interface for creating and managing UsersService with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/users
type UsersService interface {
	List(context.Context, *UserListOptions) ([]User, *Response, error)
	ListRoles(context.Context, *UserRoleListOptions) ([]UserRole, *Response, error)
	ListAssignment(context.Context, *UserRoleListOptions) ([]RoleAssignment, *Response, error)
	DeleteAssignment(context.Context, int) (*Response, error)
	UpdateAssignment(context.Context, int, *UpdateAssignmentRequest) (*Response, error)
	AssignRole(context.Context, *UpdateAssignmentRequest) (*UserRole, *Response, error)
}

// UsersServiceOp handles communication with Users methods of the EdgecenterCloud API.
type UsersServiceOp struct {
	client *Client
}

var _ UsersService = &UsersServiceOp{}

// User represents a EdgecenterCloud User configuration.
type User struct {
	Activated bool   `json:"activated"`
	IsAdmin   bool   `json:"is_admin"`
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

// UserRole represents a EdgecenterCloud User Role configuration.
type UserRole struct {
	Scope string `json:"scope"`
	Role  string `json:"role"`
}

// RoleAssignment represents a EdgecenterCloud User Role Assignment configuration.
type RoleAssignment struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id,omitempty"`
	ClientID  int    `json:"client_id,omitempty"`
	Role      string `json:"role"`
	UserID    int    `json:"user_id"`
}

// UserListOptions specifies the optional query parameters to List method.
type UserListOptions struct {
	ClientID int `url:"client_id,omitempty" validate:"omitempty"`
}

// UserRoleListOptions specifies the optional query parameters to ListRoles method.
type UserRoleListOptions struct {
	ClientID  int `url:"client_id,omitempty" validate:"omitempty"`
	ProjectID int `url:"project_id,omitempty" validate:"omitempty"`
}

type UpdateAssignmentRequest struct {
	ProjectID int    `json:"project_id,omitempty" validate:"omitempty"`
	UserID    int    `json:"user_id"  required:"true"`
	Role      string `json:"role"  required:"true"`
	ClientID  int    `json:"client_id,omitempty"  required:"omitempty"`
}

// usersRoot represents Users root.
type usersRoot struct {
	Count int
	Users []User `json:"results"`
}

// roleAssignmentsRoot represents Users Role Assignments root.
type roleAssignmentsRoot struct {
	Count           int
	RoleAssignments []RoleAssignment `json:"results"`
}

// List get client’s users.
func (s *UsersServiceOp) List(ctx context.Context, opts *UserListOptions) ([]User, *Response, error) {
	userPath, err := addOptions(usersBasePathV1, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, userPath, nil)
	if err != nil {
		return nil, nil, err
	}

	usersRoot := new(usersRoot)
	resp, err := s.client.Do(ctx, req, usersRoot)
	if err != nil {
		return nil, resp, err
	}

	return usersRoot.Users, resp, err
}

// ListRoles get available roles.
func (s *UsersServiceOp) ListRoles(ctx context.Context, opts *UserRoleListOptions) ([]UserRole, *Response, error) {
	userPath, err := addOptions(path.Join(usersBasePathV1, usersRoles), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, userPath, nil)
	if err != nil {
		return nil, nil, err
	}

	var userRole []UserRole
	resp, err := s.client.Do(ctx, req, &userRole)
	if err != nil {
		return nil, resp, err
	}

	return userRole, resp, err
}

// ListAssignment get available assignment roles.
func (s *UsersServiceOp) ListAssignment(ctx context.Context, opts *UserRoleListOptions) ([]RoleAssignment, *Response, error) {
	userPath, err := addOptions(path.Join(usersBasePathV1, usersAssignments), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, userPath, nil)
	if err != nil {
		return nil, nil, err
	}

	roleAssignmentsRoot := new(roleAssignmentsRoot)
	resp, err := s.client.Do(ctx, req, roleAssignmentsRoot)
	if err != nil {
		return nil, resp, err
	}

	return roleAssignmentsRoot.RoleAssignments, resp, err
}

// DeleteAssignment deletes a role assignment.
func (s *UsersServiceOp) DeleteAssignment(ctx context.Context, assignmentID int) (*Response, error) {
	assignmentsPath := fmt.Sprintf("%s/%s/%d", usersBasePathV1, usersAssignments, assignmentID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, assignmentsPath, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UpdateAssignment updates a role assignment.
func (s *UsersServiceOp) UpdateAssignment(ctx context.Context, assignmentID int, reqBody *UpdateAssignmentRequest) (*Response, error) {
	if reqBody == nil {
		return nil, NewArgError("reqBody", "cannot be nil")
	}

	assignmentsPath := fmt.Sprintf("%s/%s/%d", usersBasePathV1, usersAssignments, assignmentID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, assignmentsPath, reqBody)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AssignRole to an existing user.
func (s *UsersServiceOp) AssignRole(ctx context.Context, reqBody *UpdateAssignmentRequest) (*UserRole, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	assignmentsPath := fmt.Sprintf("%s/%s", usersBasePathV1, usersAssignments)

	req, err := s.client.NewRequest(ctx, http.MethodPost, assignmentsPath, reqBody)
	if err != nil {
		return nil, nil, err
	}

	userRole := new(UserRole)
	resp, err := s.client.Do(ctx, req, userRole)
	if err != nil {
		return nil, resp, err
	}

	return userRole, resp, err
}
