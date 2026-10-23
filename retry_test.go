package edgecloud

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type conflictServer struct {
	conflicts int32
	message   string
	requests  atomic.Int32
	bodies    chan string
}

func (s *conflictServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	s.bodies <- string(body)

	if s.requests.Add(1) <= s.conflicts {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message": "` + s.message + `"}`))

		return
	}

	_, _ = w.Write([]byte(`{}`))
}

func newRetryTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := New(nil, WithRetryAndBackoffs(RetryConfig{
		RetryMax:     3,
		RetryWaitMin: PtrTo(0.001),
		RetryWaitMax: PtrTo(0.001),
	}))
	require.NoError(t, err)

	client.BaseURL, err = url.Parse(server.URL)
	require.NoError(t, err)

	return client
}

func doRetryTestRequest(t *testing.T, client *Client, method string, body interface{}) (*Response, error) {
	t.Helper()

	req, err := client.NewRequest(context.Background(), method, "/v1/resource", body)
	require.NoError(t, err)

	return client.Do(context.Background(), req, nil)
}

func TestRetryPolicyTransientConflict(t *testing.T) {
	cases := []struct {
		name     string
		method   string
		message  string
		requests int32
		ok       bool
	}{
		{name: "locked resource", method: http.MethodPost, message: "Resource 7e0b is locked", requests: 3, ok: true},
		{name: "pending load balancer", method: http.MethodDelete, message: "Invalid state PENDING_UPDATE of loadbalancer resource 7e0b", requests: 3, ok: true},
		{name: "immutable load balancer on write", method: http.MethodPost, message: "Load Balancer 7e0b is immutable and cannot be updated", requests: 3, ok: true},
		{name: "immutable load balancer on read", method: http.MethodGet, message: "Load Balancer 7e0b is immutable and cannot be updated", requests: 1},
		{name: "quota conflict", method: http.MethodPost, message: "Quota exceeded for resource", requests: 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := &conflictServer{conflicts: 2, message: tc.message, bodies: make(chan string, 8)}
			client := newRetryTestClient(t, server)

			var body interface{}
			if tc.method != http.MethodGet {
				body = map[string]string{"name": "test"}
			}

			resp, err := doRetryTestRequest(t, client, tc.method, body)
			assert.Equal(t, tc.requests, server.requests.Load())

			if tc.ok {
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				close(server.bodies)
				for replayed := range server.bodies {
					assert.JSONEq(t, `{"name": "test"}`, replayed)
				}

				return
			}

			var responseError *ResponseError
			require.True(t, errors.As(err, &responseError))
			assert.Equal(t, http.StatusConflict, resp.StatusCode)
			assert.Equal(t, tc.message, responseError.Message)
		})
	}
}

func TestRetryPolicyTransientConflictGivesUp(t *testing.T) {
	server := &conflictServer{conflicts: 100, message: "Resource 7e0b is locked", bodies: make(chan string, 8)}
	client := newRetryTestClient(t, server)

	resp, err := doRetryTestRequest(t, client, http.MethodPost, map[string]string{"name": "test"})

	assert.Equal(t, int32(4), server.requests.Load())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Resource 7e0b is locked")
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}
