package edgecenter

import (
	"fmt"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

// EndpointNotFoundError is the error when no suitable endpoint can be found in the user's catalog.
type EndpointNotFoundError struct{ edgecloud.BaseError }

func (e EndpointNotFoundError) Error() string {
	return "No suitable endpoint could be found in the service catalog."
}

// InvalidAvailabilityProvidedError is the error when an invalid endpoint availability is provided.
type InvalidAvailabilityProvidedError struct{ edgecloud.InvalidInputError }

func (e InvalidAvailabilityProvidedError) Error() string {
	return fmt.Sprintf("Unexpected availability in endpoint query: %s", e.Value)
}

// NoAuthURLError is the error when the OS_AUTH_URL environment variable is not found.
type NoAuthURLError struct{ edgecloud.InvalidInputError }

func (e NoAuthURLError) Error() string {
	return "Environment variable OS_AUTH_URL needs to be set."
}

// NoUsernameError is the error when the OS_USERNAME environment variable is not found.
type NoUsernameError struct{ edgecloud.InvalidInputError }

func (e NoUsernameError) Error() string {
	return "Environment variable OS_USERNAME needs to be set."
}

// NoPasswordError is the error when the OS_PASSWORD environment variable is not found.
type NoPasswordError struct{ edgecloud.InvalidInputError }

func (e NoPasswordError) Error() string {
	return "Environment variable OS_PASSWORD needs to be set."
}
