package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []User{{ID: 123}}

	mux.HandleFunc(usersBasePathV1, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Users.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestUsers_ListRoles(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []UserRole{{Role: "test-role"}}
	URL := path.Join(usersBasePathV1, usersRoles)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, err = w.Write(resp)
		if err != nil {
			t.Fatal(err)
		}
	})

	respActual, resp, err := client.Users.ListRoles(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestUsers_ListAssignment(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []RoleAssignment{{Role: "test-role"}}
	URL := path.Join(usersBasePathV1, usersAssignments)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Users.ListAssignment(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestUsers_DeleteAssignment(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(usersBasePathV1, usersAssignments, "123")

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

	request := &UpdateAssignmentRequest{}
	URL := path.Join(usersBasePathV1, usersAssignments, "123")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(UpdateAssignmentRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
	})

	resp, err := client.Users.UpdateAssignment(ctx, 123, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestUsers_AssignRole(t *testing.T) {
	setup()
	defer teardown()

	request := &UpdateAssignmentRequest{}
	expectedResp := &UserRole{Role: "test-role"}
	URL := path.Join(usersBasePathV1, usersAssignments)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(UpdateAssignmentRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Users.AssignRole(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
