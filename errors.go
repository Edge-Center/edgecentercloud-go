package edgecloud

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BaseError is an error type that all other error types embed.
type BaseError struct {
	DefaultErrString string
	Info             string
}

func (e BaseError) Error() string {
	e.DefaultErrString = "An error occurred while executing a EdgeCenter cloud request."
	return e.choseErrString()
}

func (e BaseError) choseErrString() string {
	if e.Info != "" {
		return e.Info
	}
	return e.DefaultErrString
}

// MissingInputError is the error when input is required in a particular situation but not provided by the user.
type MissingInputError struct {
	BaseError
	Argument string
}

func (e MissingInputError) Error() string {
	e.DefaultErrString = fmt.Sprintf("Missing input for argument [%s]", e.Argument)
	return e.choseErrString()
}

// InvalidInputError is an error type used for most non-HTTP Gophercloud errors.
type InvalidInputError struct {
	MissingInputError
	Value interface{}
}

func (e InvalidInputError) Error() string {
	e.DefaultErrString = fmt.Sprintf("Invalid input provided for argument [%s]: [%+v]", e.Argument, e.Value)
	return e.choseErrString()
}

// MissingEnvironmentVariableError is the error when environment variable is required
// in a particular situation but not provided by the user.
type MissingEnvironmentVariableError struct {
	BaseError
	EnvironmentVariable string
}

func (e MissingEnvironmentVariableError) Error() string {
	e.DefaultErrString = fmt.Sprintf("Missing environment variable [%s]", e.EnvironmentVariable)
	return e.choseErrString()
}

// MissingAnyoneOfEnvironmentVariablesError is the error when anyone of the environment variables
// is required in a particular situation but not provided by the user.
type MissingAnyoneOfEnvironmentVariablesError struct {
	BaseError
	EnvironmentVariables []string
}

func (e MissingAnyoneOfEnvironmentVariablesError) Error() string {
	e.DefaultErrString = fmt.Sprintf(
		"Missing one of the following environment variables [%s]",
		strings.Join(e.EnvironmentVariables, ", "),
	)
	return e.choseErrString()
}

// UnexpectedResponseCodeError is returned by the Request method when a response code other than
// those listed in OkCodes is encountered.
type UnexpectedResponseCodeError struct {
	BaseError
	URL      string
	Method   string
	Expected []int
	Actual   int
	Body     []byte
}

func (e UnexpectedResponseCodeError) Error() string {
	e.DefaultErrString = fmt.Sprintf(
		"Expected HTTP response code %v when accessing [%s %s], but got %d instead\n%s",
		e.Expected, e.Method, e.URL, e.Actual, e.Body,
	)
	ecErr := ECErrorType{}
	err := json.Unmarshal(e.Body, &ecErr)
	if err != nil {
		e.Info = ecErr.Message
	}
	return e.choseErrString()
}

func (e *UnexpectedResponseCodeError) ReadEcError() {
	ecErr := ECErrorType{}
	err := json.Unmarshal(e.Body, &ecErr)
	if err == nil {
		e.Info = ecErr.Message
	}
}

// GetStatusCode returns the actual status code of the error.
func (e UnexpectedResponseCodeError) GetStatusCode() int {
	return e.Actual
}

// StatusCodeError is a convenience interface to easily allow access to the
// status code field of the various ErrDefault* types.
//
// By using this interface, you only have to make a single type cast of
// the returned error to err.(StatusCodeError) and then call GetStatusCode()
// instead of having a large switch statement checking for each of the
// ErrDefault* types.
type StatusCodeError interface {
	Error() string
	GetStatusCode() int
}

// Default400Error is the default error type returned on a 400 HTTP response code.
type Default400Error struct {
	UnexpectedResponseCodeError
}

// Default401Error is the default error type returned on a 401 HTTP response code.
type Default401Error struct {
	UnexpectedResponseCodeError
}

// Default403Error is the default error type returned on a 403 HTTP response code.
type Default403Error struct {
	UnexpectedResponseCodeError
}

// Default404Error is the default error type returned on a 404 HTTP response code.
type Default404Error struct {
	UnexpectedResponseCodeError
}

// Default405Error is the default error type returned on a 405 HTTP response code.
type Default405Error struct {
	UnexpectedResponseCodeError
}

// Default408Error is the default error type returned on a 408 HTTP response code.
type Default408Error struct {
	UnexpectedResponseCodeError
}

// Default409Error is the default error type returned on a 409 HTTP response code.
type Default409Error struct {
	UnexpectedResponseCodeError
}

// Default429Error is the default error type returned on a 429 HTTP response code.
type Default429Error struct {
	UnexpectedResponseCodeError
}

// Default500Error is the default error type returned on a 500 HTTP response code.
type Default500Error struct {
	UnexpectedResponseCodeError
}

// Default503Error is the default error type returned on a 503 HTTP response code.
type Default503Error struct {
	UnexpectedResponseCodeError
}

func (e Default400Error) Error() string {
	e.DefaultErrString = fmt.Sprintf(
		"Bad request with: [%s %s], error message: %s",
		e.Method, e.URL, e.Body,
	)
	e.ReadEcError()
	if e.Info != "" {
		return e.choseErrString()
	}
	return e.choseErrString()
}

func (e Default401Error) Error() string {
	e.ReadEcError()
	if e.Info != "" {
		return e.choseErrString()
	}
	return "Authentication failed"
}

func (e Default403Error) Error() string {
	e.DefaultErrString = fmt.Sprintf(
		"Request forbidden: [%s %s], error message: %s",
		e.Method, e.URL, e.Body,
	)
	e.ReadEcError()
	if e.Info != "" {
		return e.choseErrString()
	}
	return e.choseErrString()
}

func (e Default404Error) Error() string {
	e.ReadEcError()
	if e.Info != "" {
		return e.choseErrString()
	}
	return "Resource not found"
}

func (e Default405Error) Error() string {
	return "Method not allowed"
}

func (e Default408Error) Error() string {
	return "The server timed out waiting for the request"
}

func (e Default409Error) Error() string {
	e.ReadEcError()
	if e.Info != "" {
		return e.choseErrString()
	}
	return "Conflict"
}

func (e Default429Error) Error() string {
	return "Too many requests have been sent in a given amount of time. Pause" +
		" requests, wait up to one minute, and try again."
}

func (e Default500Error) Error() string {
	e.ReadEcError()
	if e.Info != "" {
		return e.choseErrString()
	}
	return "Internal Server Error"
}

func (e Default503Error) Error() string {
	return "The service is currently unable to handle the request due to a temporary" +
		" overloading or maintenance. This is a temporary condition. Try again later."
}

// Err400er is the interface resource error types implement to override the error message
// from a 400 error.
type Err400er interface {
	Error400(UnexpectedResponseCodeError) error
}

// Err401er is the interface resource error types implement to override the error message
// from a 401 error.
type Err401er interface {
	Error401(UnexpectedResponseCodeError) error
}

// Err403er is the interface resource error types implement to override the error message
// from a 403 error.
type Err403er interface {
	Error403(UnexpectedResponseCodeError) error
}

// Err404er is the interface resource error types implement to override the error message
// from a 404 error.
type Err404er interface {
	Error404(UnexpectedResponseCodeError) error
}

// Err405er is the interface resource error types implement to override the error message
// from a 405 error.
type Err405er interface {
	Error405(UnexpectedResponseCodeError) error
}

// Err408er is the interface resource error types implement to override the error message
// from a 408 error.
type Err408er interface {
	Error408(UnexpectedResponseCodeError) error
}

// Err409er is the interface resource error types implement to override the error message
// from a 409 error.
type Err409er interface {
	Error409(UnexpectedResponseCodeError) error
}

// Err429er is the interface resource error types implement to override the error message
// from a 429 error.
type Err429er interface {
	Error429(UnexpectedResponseCodeError) error
}

// Err500er is the interface resource error types implement to override the error message
// from a 500 error.
type Err500er interface {
	Error500(UnexpectedResponseCodeError) error
}

// Err503er is the interface resource error types implement to override the error message
// from a 503 error.
type Err503er interface {
	Error503(UnexpectedResponseCodeError) error
}

// TimeOutError is the error type returned when an operations time out.
type TimeOutError struct {
	BaseError
}

func (e TimeOutError) Error() string {
	e.DefaultErrString = "A time out occurred"
	return e.choseErrString()
}

// UnableToReauthenticateError is the error type returned when reauthentication fails.
type UnableToReauthenticateError struct {
	BaseError
	ErrOriginal error
}

func (e UnableToReauthenticateError) Error() string {
	e.DefaultErrString = fmt.Sprintf("Unable to re-authenticate: %s", e.ErrOriginal)
	return e.choseErrString()
}

// AfterReauthenticationError is the error type returned when reauthentication
// succeeds, but an error occurs afterword (usually an HTTP error).
type AfterReauthenticationError struct {
	BaseError
	ErrOriginal error
}

func (e AfterReauthenticationError) Error() string {
	e.DefaultErrString = fmt.Sprintf("Successfully re-authenticated, but got error executing request: %s", e.ErrOriginal)
	return e.choseErrString()
}

// ServiceNotFoundError is returned when no service in a service catalog matches
// the provided EndpointOpts. This is generally returned by provider service
// factory methods like "NewComputeV2()" and can mean that a service is not
// enabled for your account.
type ServiceNotFoundError struct {
	BaseError
}

func (e ServiceNotFoundError) Error() string {
	e.DefaultErrString = "No suitable service could be found in the service catalog."
	return e.choseErrString()
}

// EndpointNotFoundError is returned when no available endpoints match the
// provided EndpointOpts. This is also generally returned by provider service
// factory methods, and usually indicates that a region was specified
// incorrectly.
type EndpointNotFoundError struct {
	BaseError
}

func (e EndpointNotFoundError) Error() string {
	e.DefaultErrString = "No suitable endpoint could be found in the service catalog."
	return e.choseErrString()
}

// ResourceNotFoundError is the error when trying to retrieve a resource's
// ID by name and the resource doesn't exist.
type ResourceNotFoundError struct {
	BaseError
	Name         string
	ResourceType string
}

func (e ResourceNotFoundError) Error() string {
	e.DefaultErrString = fmt.Sprintf("Unable to find %s with name %s", e.ResourceType, e.Name)
	return e.choseErrString()
}

// MultipleResourcesFoundError is the error when trying to retrieve a resource's
// ID by name and multiple resources have the user-provided name.
type MultipleResourcesFoundError struct {
	BaseError
	Name         string
	Count        int
	ResourceType string
}

func (e MultipleResourcesFoundError) Error() string {
	e.DefaultErrString = fmt.Sprintf("Found %d %ss matching %s", e.Count, e.ResourceType, e.Name)
	return e.choseErrString()
}

// UnexpectedTypeError is the error when an unexpected type is encountered.
type UnexpectedTypeError struct {
	BaseError
	Expected string
	Actual   string
}

func (e UnexpectedTypeError) Error() string {
	e.DefaultErrString = fmt.Sprintf("Expected %s but got %s", e.Expected, e.Actual)
	return e.choseErrString()
}

func unacceptedAttributeErr(attribute string) string {
	return fmt.Sprintf("The base Identity V3 API does not accept authentication by %s", attribute)
}

func redundantWithTokenErr(attribute string) string {
	return fmt.Sprintf("%s may not be provided when authenticating with a AccessTokenID", attribute)
}

func redundantWithUserID(attribute string) string {
	return fmt.Sprintf("%s may not be provided when authenticating with a UserID", attribute)
}

// APIKeyProvidedError indicates that an APIKey was provided but can't be used.
type APIKeyProvidedError struct{ BaseError }

func (e APIKeyProvidedError) Error() string {
	return unacceptedAttributeErr("APIKey")
}

// TenantIDProvidedError indicates that a TenantID was provided but can't be used.
type TenantIDProvidedError struct{ BaseError }

func (e TenantIDProvidedError) Error() string {
	return unacceptedAttributeErr("TenantID")
}

// TenantNameProvidedError indicates that a TenantName was provided but can't be used.
type TenantNameProvidedError struct{ BaseError }

func (e TenantNameProvidedError) Error() string {
	return unacceptedAttributeErr("TenantName")
}

// UsernameWithTokenError indicates that a Username was provided, but token authentication is being used instead.
type UsernameWithTokenError struct{ BaseError }

func (e UsernameWithTokenError) Error() string {
	return redundantWithTokenErr("Username")
}

// UserIDWithTokenError indicates that a UserID was provided, but token authentication is being used instead.
type UserIDWithTokenError struct{ BaseError }

func (e UserIDWithTokenError) Error() string {
	return redundantWithTokenErr("UserID")
}

// DomainIDWithTokenError indicates that a DomainID was provided, but token authentication is being used instead.
type DomainIDWithTokenError struct{ BaseError }

func (e DomainIDWithTokenError) Error() string {
	return redundantWithTokenErr("DomainID")
}

// DomainNameWithTokenError indicates that a DomainName was provided, but token authentication is being used instead.
type DomainNameWithTokenError struct{ BaseError }

func (e DomainNameWithTokenError) Error() string {
	return redundantWithTokenErr("DomainName")
}

// UsernameOrUserIDError indicates that neither username nor userID are specified, or both are at once.
type UsernameOrUserIDError struct{ BaseError }

func (e UsernameOrUserIDError) Error() string {
	return "Exactly one of Username and UserID must be provided for password authentication"
}

// DomainIDWithUserIDError indicates that a DomainID was provided, but unnecessary because a UserID is being used.
type DomainIDWithUserIDError struct{ BaseError }

func (e DomainIDWithUserIDError) Error() string {
	return redundantWithUserID("DomainID")
}

// DomainNameWithUserIDError indicates that a DomainName was provided, but unnecessary because a UserID is being used.
type DomainNameWithUserIDError struct{ BaseError }

func (e DomainNameWithUserIDError) Error() string {
	return redundantWithUserID("DomainName")
}

// DomainIDOrDomainNameError indicates that a username was provided, but no domain to scope it.
// It may also indicate that both a DomainID and a DomainName were provided at once.
type DomainIDOrDomainNameError struct{ BaseError }

func (e DomainIDOrDomainNameError) Error() string {
	return "You must provide exactly one of DomainID or DomainName to authenticate by Username"
}

// MissingPasswordError indicates that no password was provided and no token is available.
type MissingPasswordError struct{ BaseError }

func (e MissingPasswordError) Error() string {
	return "You must provide a password to authenticate"
}

// ScopeDomainIDOrDomainNameError indicates that a domain ID or Name was required in a Scope, but not present.
type ScopeDomainIDOrDomainNameError struct{ BaseError }

func (e ScopeDomainIDOrDomainNameError) Error() string {
	return "You must provide exactly one of DomainID or DomainName in a Scope with ProjectName"
}

// ScopeProjectIDOrProjectNameError indicates that both a ProjectID and a ProjectName were provided in a Scope.
type ScopeProjectIDOrProjectNameError struct{ BaseError }

func (e ScopeProjectIDOrProjectNameError) Error() string {
	return "You must provide at most one of ProjectID or ProjectName in a Scope"
}

// ScopeProjectIDAloneError indicates that a ProjectID was provided with other constraints in a Scope.
type ScopeProjectIDAloneError struct{ BaseError }

func (e ScopeProjectIDAloneError) Error() string {
	return "ProjectID must be supplied alone in a Scope"
}

// ScopeEmptyError indicates that no credentials were provided in a Scope.
type ScopeEmptyError struct{ BaseError }

func (e ScopeEmptyError) Error() string {
	return "You must provide either a Project or Domain in a Scope"
}

// AppCredMissingSecretError indicates that no Application Credential Secret was provided with Application Credential ID or Name.
type AppCredMissingSecretError struct{ BaseError }

func (e AppCredMissingSecretError) Error() string {
	return "You must provide an Application Credential Secret"
}
