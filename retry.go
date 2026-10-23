package edgecloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

var transientConflictMarkers = []string{
	"is locked",
	"invalid state pending_",
}

const immutableConflictMarker = "is immutable and cannot be updated"

func retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	retry, policyErr := retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	if retry || policyErr != nil || resp == nil || resp.StatusCode != http.StatusConflict {
		return retry, policyErr
	}

	return isTransientConflict(resp), nil
}

func isTransientConflict(resp *http.Response) bool {
	if resp.Body == nil {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return false
	}

	message := strings.ToLower(string(body))
	for _, marker := range transientConflictMarkers {
		if strings.Contains(message, marker) {
			return true
		}
	}

	return resp.Request != nil && isMutatingMethod(resp.Request.Method) &&
		strings.Contains(message, immutableConflictMarker)
}

func isMutatingMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}
