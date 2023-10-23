package util

import (
	"context"
	"errors"
	"net/http"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var (
	errResourceNotDeleted = errors.New("could not delete the resource")
	errGetResourceInfo    = errors.New("error when retrieving resource information")
)

type RetrieveResourceFunc[T any] func(ctx context.Context, id string) (*T, *edgecloud.Response, error)

func ResourceIsDeleted[T any](ctx context.Context, retrieveResourceFunc RetrieveResourceFunc[T], id string) error {
	_, resp, err := retrieveResourceFunc(ctx, id)
	if err == nil {
		return errResourceNotDeleted
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return errGetResourceInfo
}
