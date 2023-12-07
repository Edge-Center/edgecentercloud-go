package edgecloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

const (
	projectID                  = 27520
	regionID                   = 8
	clientID                   = 88888
	testResourceID             = "f0d19cec-5c3f-4853-886e-304915960ff6"
	taskID                     = "f0d19cec-5c3f-4853-886e-304915960ff6"
	testResourceIDNotValidUUID = "123"
)

var NotCorrectUUID = fmt.Sprintf("should be the correct UUID. current value is: %s", testResourceIDNotValidUUID)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Region = regionID
	client.Project = projectID
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	t.Helper()

	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func testClientDefaultBaseURL(t *testing.T, c *Client) {
	t.Helper()

	if c.BaseURL == nil || c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, defaultBaseURL)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	t.Helper()

	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", c.UserAgent, userAgent)
	}
}

func testClientDefaults(t *testing.T, c *Client) {
	t.Helper()

	testClientDefaultBaseURL(t, c)
	testClientDefaultUserAgent(t, c)
}

func testURLParseError(t *testing.T, urlErr error) {
	t.Helper()

	if urlErr == nil {
		t.Errorf("Expected error to be returned")
	}

	var err *url.Error
	if !errors.As(urlErr, &err) || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestNewWithRetries(t *testing.T) {
	c, err := NewWithRetries(nil)
	if err != nil {
		t.Fatalf("NewWithRetries(): %v", err)
	}
	testClientDefaults(t, c)
}

func TestNew(t *testing.T) {
	c, err := New(nil)
	if err != nil {
		t.Fatalf("New(): %v", err)
	}
	testClientDefaults(t, c)
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	testClientDefaults(t, c)
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "inst", defaultBaseURL+"inst"
	inBody, outBody := &InstanceCreateRequest{Names: []string{"inst"}},
		`{"names":["inst"],"flavor":"","interfaces":null,"volumes":null}`+"\n"
	req, _ := c.NewRequest(ctx, http.MethodPost, inURL, inBody)

	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	body, _ := io.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v)Body = %v, expected %v", inBody, string(body), outBody)
	}

	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.UserAgent)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest(ctx, http.MethodGet, ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_withCustomUserAgent(t *testing.T) {
	ua := "testing/0.0.1"
	c, err := New(nil, SetUserAgent(ua))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req, _ := c.NewRequest(ctx, http.MethodGet, "/foo", nil)

	expected := fmt.Sprintf("%s %s", ua, userAgent)
	if got := req.Header.Get("User-Agent"); got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

func TestNewRequest_withCustomHeaders(t *testing.T) {
	expectedIdentity := "identity"
	expectedCustom := "x_test_header"

	c, err := New(nil, SetRequestHeaders(map[string]string{
		"Accept-Encoding": expectedIdentity,
		"X-Test-Header":   expectedCustom,
	}))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req, _ := c.NewRequest(ctx, http.MethodGet, "/foo", nil)

	if got := req.Header.Get("Accept-Encoding"); got != expectedIdentity {
		t.Errorf("New() Custom Accept Encoding Header = %s; expected %s", got, expectedIdentity)
	}
	if got := req.Header.Get("X-Test-Header"); got != expectedCustom {
		t.Errorf("New() Custom Accept Encoding Header = %s; expected %s", got, expectedCustom)
	}
}

func TestNewRequest_withApiKeyToken(t *testing.T) {
	apiKey := "4010$252a09"
	c, err := New(nil, SetAPIKey(apiKey))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req, _ := c.NewRequest(ctx, http.MethodGet, "/foo", nil)

	expected := fmt.Sprintf("APIKey %s", apiKey)
	if got := req.Header.Get("Authorization"); got != expected {
		t.Errorf("New() APIKey = %s; expected %s", got, apiKey)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		_, _ = fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

func TestCustomUserAgent(t *testing.T) {
	ua := "testing/0.0.1"
	c, err := New(nil, SetUserAgent(ua))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := fmt.Sprintf("%s %s", ua, userAgent)
	if got := c.UserAgent; got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

func TestCustomBaseURL(t *testing.T) {
	baseURL := "http://localhost/foo"
	c, err := New(nil, SetBaseURL(baseURL))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := baseURL
	if got := c.BaseURL.String(); got != expected {
		t.Errorf("New() BaseURL = %s; expected %s", got, expected)
	}
}

func TestSetProject(t *testing.T) {
	c, err := New(nil, SetProject(projectID))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := projectID
	if got := c.Project; got != expected {
		t.Errorf("New() Project = %d; expected %d", got, expected)
	}
}

func TestSetRegion(t *testing.T) {
	c, err := New(nil, SetRegion(regionID))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := regionID
	if got := c.Region; got != expected {
		t.Errorf("New() Region = %d; expected %d", got, expected)
	}
}

func TestCustomBaseURL_badURL(t *testing.T) {
	baseURL := ":"
	_, err := New(nil, SetBaseURL(baseURL))

	testURLParseError(t, err)
}

func TestWithRetryAndBackoffs(t *testing.T) {
	setup()
	defer teardown()

	baseURL, _ := url.Parse(server.URL)
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message": "broken"}`))
	})

	waitMax := PtrTo(5.0)
	waitMin := PtrTo(2.0)

	retryConfig := RetryConfig{
		RetryMax:     2,
		RetryWaitMin: waitMin,
		RetryWaitMax: waitMax,
	}

	client, err := New(nil, WithRetryAndBackoffs(retryConfig))
	client.BaseURL = baseURL
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req, err := client.NewRequest(ctx, http.MethodGet, "/foo", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expectingErr := fmt.Sprintf("GET %s/foo: 500 broken; giving up after 3 attempt(s)", baseURL)
	_, err = client.Do(context.Background(), req, nil)
	if err == nil || (err.Error() != expectingErr) {
		t.Fatalf("expected giving up error, got: %#v", err)
	}
}

func TestWithRetryAndBackoffsForResourceMethods(t *testing.T) {
	// Mock server which always responds 500.
	setup()
	defer teardown()

	const (
		networkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	baseURL, _ := url.Parse(server.URL)
	getNetworkURL := fmt.Sprintf("/v1/networks/%d/%d/%s", projectID, regionID, networkID)
	mux.HandleFunc(getNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"message": "broken"}`))
	})

	waitMax := PtrTo(5.0)
	waitMin := PtrTo(2.0)

	retryConfig := RetryConfig{
		RetryMax:     2,
		RetryWaitMin: waitMin,
		RetryWaitMax: waitMax,
	}

	client, err := New(nil, WithRetryAndBackoffs(retryConfig))
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expectingErr := fmt.Sprintf("GET %s%s: 502 broken; giving up after 3 attempt(s)", baseURL, getNetworkURL)
	_, resp, err := client.Networks.Get(context.Background(), networkID)
	if err == nil || (err.Error() != expectingErr) {
		t.Fatalf("expected giving up error, got: %s", err.Error())
	}

	var responseError *ResponseError
	if !errors.As(err, &responseError) {
		t.Fatalf("expected error to be *edgecloud.ResponseError, got: %#v", err)
	}

	// Ensure that the *Response is properly populated
	if resp == nil {
		t.Fatal("expected non-nil *edgecloud.Response")
	}
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected %d, got: %d", http.StatusBadGateway, resp.StatusCode)
	}
}
