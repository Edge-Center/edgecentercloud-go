package edgecloud

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func isValidUUID(id, name string) (*Response, error) {
	if _, err := uuid.Parse(id); err != nil {
		return &Response{
			Response: &http.Response{
				Status:     http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
			},
		}, NewArgError(name, fmt.Sprintf("should be the correct UUID. current value is: %s", id))
	}

	return nil, nil // nolint
}
