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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	libraryVersion = "1.0.0"
	defaultBaseURL = "https://api.edgecenter.ru/cloud/"
	userAgent      = "edgecloud/" + libraryVersion
	mediaType      = "application/json"

	internalHeaderRetryAttempts = "X-Edgecloud-Retry-Attempts"

	defaultRetryMax     = 3
	defaultRetryWaitMax = 30
	defaultRetryWaitMin = 1
)

// Client manages communication with EdgecenterCloud API.
type Client struct {
	// HTTP client used to communicate with the EdgecenterCloud API.
	HTTPClient *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// APIKey token for client
	APIKey string

	// RegionID for client
	Region int

	// ProjectID for client
	Project int

	Flavors        FlavorsService
	Floatingips    FloatingIPsService
	Images         ImagesService
	Instances      InstancesService
	KeyPairs       KeyPairsService
	Loadbalancers  LoadbalancersService
	Networks       NetworksService
	Ports          PortsService
	Projects       ProjectsService
	Quotas         QuotasService
	Regions        RegionsService
	SecurityGroups SecurityGroupsService
	Secrets        SecretsService
	ServerGroups   ServerGroupsService
	Subnetworks    SubnetworksService
	Tasks          TasksService
	Volumes        VolumesService
	Users          UsersService

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback

	// Optional extra HTTP headers to set on every request to the API.
	headers map[string]string

	// Optional retry values. Setting the RetryConfig.RetryMax value enables automatically retrying requests
	// that fail with 429 or 500-level response codes
	RetryConfig RetryConfig
}

// RetryConfig sets the values used for enabling retries and backoffs for
// requests that fail with 429 or 500-level response codes using the go-retryablehttp client.
// RetryConfig.RetryMax must be configured to enable this behavior. RetryConfig.RetryWaitMin and
// RetryConfig.RetryWaitMax are optional, with the default values being 1.0 and 30.0, respectively.
//
// Note: Opting to use the go-retryablehttp client will overwrite any custom HTTP client passed into New().
type RetryConfig struct {
	RetryMax     int
	RetryWaitMin *float64    // Minimum time to wait
	RetryWaitMax *float64    // Maximum time to wait
	Logger       interface{} // Customer logger instance. Must implement either go-retryablehttp.Logger or go-retryablehttp.LeveledLogger
}

// RequestCompletionCallback defines the type of the request callback function.
type RequestCompletionCallback func(*http.Request, *http.Response)

func (c *Client) addProjectRegionPath(s string) string {
	projectStr := strconv.Itoa(c.Project)
	regionStr := strconv.Itoa(c.Region)

	return path.Join(s, projectStr, regionStr)
}

func (c *Client) addRegionPath(s string) string {
	regionStr := strconv.Itoa(c.Region)

	return path.Join(s, regionStr)
}

func (c *Client) Validate() (*Response, error) {
	badResponse := &Response{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusBadRequest),
			StatusCode: http.StatusBadRequest,
		},
	}
	if c.Project == 0 {
		return badResponse, NewArgError("Client.Project", "is not set")
	}
	if c.Region == 0 {
		return badResponse, NewArgError("Client.Region", "is not set")
	}

	return nil, nil // nolint
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

	// Attempts is the number of times the request was attempted when retries are enabled.
	Attempts int
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()

	return origURL.String(), nil
}

// NewWithRetries returns a new EdgecenterCloud API client with default retries config.
func NewWithRetries(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	opts = append(opts, WithRetryAndBackoffs(
		RetryConfig{
			RetryMax:     defaultRetryMax,
			RetryWaitMin: PtrTo(float64(defaultRetryWaitMin)),
			RetryWaitMax: PtrTo(float64(defaultRetryWaitMax)),
		},
	))

	return New(httpClient, opts...)
}

// NewClient returns a new EdgecenterCloud API, using the given
// http.Client to perform all requests.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{HTTPClient: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	c.Flavors = &FlavorsServiceOp{client: c}
	c.Floatingips = &FloatingipsServiceOp{client: c}
	c.Images = &ImagesServiceOp{client: c}
	c.Instances = &InstancesServiceOp{client: c}
	c.KeyPairs = &KeyPairsServiceOp{client: c}
	c.Loadbalancers = &LoadbalancersServiceOp{client: c}
	c.Networks = &NetworksServiceOp{client: c}
	c.Ports = &PortsServiceOp{client: c}
	c.Projects = &ProjectsServiceOp{client: c}
	c.Quotas = &QuotasServiceOp{client: c}
	c.Regions = &RegionsServiceOp{client: c}
	c.SecurityGroups = &SecurityGroupsServiceOp{client: c}
	c.Secrets = &SecretsServiceOp{client: c}
	c.ServerGroups = &ServerGroupsServiceOp{client: c}
	c.Subnetworks = &SubnetworksServiceOp{client: c}
	c.Tasks = &TasksServiceOp{client: c}
	c.Volumes = &VolumesServiceOp{client: c}
	c.Users = &UsersServiceOp{client: c}

	c.headers = make(map[string]string)

	return c
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new EdgecenterCloud API client instance.
func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	c := NewClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// if retryMax is set it will use the retryablehttp client.
	if c.RetryConfig.RetryMax > 0 {
		retryableClient := retryablehttp.NewClient()
		retryableClient.RetryMax = c.RetryConfig.RetryMax

		if c.RetryConfig.RetryWaitMin != nil {
			retryableClient.RetryWaitMin = time.Duration(*c.RetryConfig.RetryWaitMin * float64(time.Second))
		}
		if c.RetryConfig.RetryWaitMax != nil {
			retryableClient.RetryWaitMax = time.Duration(*c.RetryConfig.RetryWaitMax * float64(time.Second))
		}

		// By default, this is nil and does not log.
		retryableClient.Logger = c.RetryConfig.Logger

		// if timeout is set, it is maintained before overwriting client with StandardClient()
		retryableClient.HTTPClient.Timeout = c.HTTPClient.Timeout

		// This custom ErrorHandler is required to provide errors that are consistent
		// with a *edgecloud.ErrorResponse and a non-nil *edgecloud.Response while providing
		// insight into retries using an internal header.
		retryableClient.ErrorHandler = func(resp *http.Response, err error, numTries int) (*http.Response, error) {
			if resp != nil {
				resp.Header.Add(internalHeaderRetryAttempts, strconv.Itoa(numTries))

				return resp, err
			}

			return resp, err
		}

		c.HTTPClient = retryableClient.StandardClient()
	}

	return c, nil
}

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u

		return nil
	}
}

// SetAPIKey is a client option for setting the APIKey token.
func SetAPIKey(apiKey string) ClientOpt {
	return func(c *Client) error {
		tokenPartsCount := 2
		parts := strings.SplitN(apiKey, " ", tokenPartsCount)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "apikey" {
			apiKey = parts[1]
		}
		c.APIKey = apiKey
		c.headers["Authorization"] = fmt.Sprintf("APIKey %s", c.APIKey)

		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}

// SetRequestHeaders sets optional HTTP headers on the client that are
// sent on each HTTP request.
func SetRequestHeaders(headers map[string]string) ClientOpt {
	return func(c *Client) error {
		for k, v := range headers {
			c.headers[k] = v
		}
		return nil
	}
}

// WithRetryAndBackoffs sets retry values. Setting the RetryConfig.RetryMax value enables automatically retrying requests
// that fail with 429 or 500-level response codes using the go-retryablehttp client.
func WithRetryAndBackoffs(retryConfig RetryConfig) ClientOpt {
	return func(c *Client) error {
		c.RetryConfig.RetryMax = retryConfig.RetryMax
		c.RetryConfig.RetryWaitMax = retryConfig.RetryWaitMax
		c.RetryConfig.RetryWaitMin = retryConfig.RetryWaitMin
		c.RetryConfig.Logger = retryConfig.Logger
		return nil
	}
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
		return &Response{
			Response: &http.Response{
				Status:     http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
			},
		}, err
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

		if rErr := resp.Body.Close(); err == nil {
			err = rErr
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
			return &Response{
				Response: &http.Response{
					Status:     http.StatusText(http.StatusInternalServerError),
					StatusCode: http.StatusInternalServerError,
				},
			}, err
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
	var attempted string
	if r.Attempts > 0 {
		attempted = fmt.Sprintf("; giving up after %d attempt(s)", r.Attempts)
	}

	return fmt.Sprintf("%v %v: %d %v%s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, attempted)
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

	attempts, strconvErr := strconv.Atoi(r.Header.Get(internalHeaderRetryAttempts))
	if strconvErr == nil {
		errorResponse.Attempts = attempts
	}

	return errorResponse
}

// PtrTo returns a pointer to the provided input.
func PtrTo[T any](v T) *T {
	return &v
}
