package edgecloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	libraryVersion = "1.0.0"
	defaultBaseURL = "https://api.edgecenter.ru/cloud/"
	userAgent      = "edgecloud/" + libraryVersion
	mediaType      = "application/json"
)

// Client manages communication with EdgecenterCloud API.
type Client struct {
	// HTTP client used to communicate with the EdgecenterCloud API.
	HTTPClient *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// APIToken for client
	APIToken string

	Floatingips   FloatingipsService
	Instances     InstancesService
	Loadbalancers LoadbalancersService
	Networks      NetworksService
	Projects      ProjectsService
	Subnetworks   SubnetworksService
	Volumes       VolumesService

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback

	// Optional extra HTTP headers to set on every request to the API.
	headers map[string]string
}

// RequestCompletionCallback defines the type of the request callback function.
type RequestCompletionCallback func(*http.Request, *http.Response)

// ServicePath specifies additional paths for service endpoints.
type ServicePath struct {
	// Region is EdgeCenter region
	Region string

	// Project is EdgeCenter project
	Project string
}

func addServicePath(s string, p *ServicePath) string {
	return path.Join(s, p.Project, p.Region)
}

// TaskResponse is an EdgecenterCloud response with list of created tasks.
type TaskResponse struct {
	Tasks []string `json:"tasks"`
}

// Response is a EdgecenterCloud response. This wraps the standard http.Response returned from EdgecenterCloud.
type Response struct {
	*http.Response
}

// An ResponseError reports the error caused by an API request.
type ResponseError struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`
}

func NewClientWithAPIToken(httpClient *http.Client, apiToken string) *Client {
	c := NewClient(httpClient)

	c.APIToken = apiToken
	c.headers["Authorization"] = fmt.Sprintf("APIKey %s", apiToken)

	return c
}

// NewClient returns a new EdgecenterCloud API, using the given
// http.Client to perform all requests.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{HTTPClient: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	c.Floatingips = &FloatingipsServiceOp{client: c}
	c.Instances = &InstancesServiceOp{client: c}
	c.Loadbalancers = &LoadbalancersServiceOp{client: c}
	c.Networks = &NetworksServiceOp{client: c}
	c.Projects = &ProjectsServiceOp{client: c}
	c.Subnetworks = &SubnetworksServiceOp{client: c}
	c.Volumes = &VolumesServiceOp{client: c}

	c.headers = make(map[string]string)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(_ context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCPConnection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			_, _ = io.CopyN(io.Discard, resp.Body, maxBodySlurpSize)
		}

		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != http.StatusNoContent && v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
		if err != nil {
			return nil, err
		}
	}

	return response, err
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	return client.Do(req)
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ResponseError. Any other response body will be silently ignored.
// If the API error response does not include the request ID in its body, the one from its header will be used.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ResponseError{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}
