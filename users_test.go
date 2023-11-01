package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers_List(t *testing.T) {
	setup()
	defer teardown()

	users := []User{{ID: 123}}
	URL := "/v1/users"

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(users)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Users.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, users) {
		t.Errorf("Users.List\n returned %+v,\n expected %+v", resp, users)
	}
}

func TestUsers_ListRoles(t *testing.T) {
	setup()
	defer teardown()

	roles := []UserRole{{Role: "test-role"}}
	URL := "/v1/users/roles"

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(roles)
		_, err := w.Write(resp)
		if err != nil {
			t.Fatal(err)
		}
	})

	resp, _, err := client.Users.ListRoles(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, roles) {
		t.Errorf("Users.List\n returned %+v,\n expected %+v", resp, roles)
	}
}

func TestUsers_ListAssignment(t *testing.T) {
	setup()
	defer teardown()

	assignments := []RoleAssignment{{Role: "test-role"}}
	URL := "/v1/users/assignments"

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(assignments)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Users.ListAssignment(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, assignments) {
		t.Errorf("Users.ListAssignment\n returned %+v,\n expected %+v", resp, assignments)
	}
}

func TestUsers_DeleteAssignment(t *testing.T) {
	setup()
	defer teardown()

	URL := fmt.Sprintf("/v1/users/assignments/%d", 123)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Users.DeleteAssignment(ctx, 123)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestUsers_UpdateAssignment(t *testing.T) {
	setup()
	defer teardown()

	updateAssignmentRequest := &UpdateAssignmentRequest{}
	URL := fmt.Sprintf("/v1/users/assignments/%d", 123)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(UpdateAssignmentRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, updateAssignmentRequest, reqBody)
	})

	resp, err := client.Users.UpdateAssignment(ctx, 123, updateAssignmentRequest)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestUsers_AssignRole(t *testing.T) {
	setup()
	defer teardown()

	updateAssignmentRequest := &UpdateAssignmentRequest{}
	userRoleResponse := &UserRole{Role: "test-role"}
	URL := "/v1/users/assignments"

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(UpdateAssignmentRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, updateAssignmentRequest, reqBody)
		resp, _ := json.Marshal(userRoleResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Users.AssignRole(ctx, updateAssignmentRequest)
	require.NoError(t, err)

	assert.Equal(t, userRoleResponse, resp)
}
