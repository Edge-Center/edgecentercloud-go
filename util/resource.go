package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

var (
	errResourceNotDeleted                  = errors.New("could not delete the resource")
	errGetResourceInfo                     = errors.New("error when retrieving resource information")
	errDeleteResourceIfExistIsNotSupported = errors.New("method DeleteResourceIfExist isn't supported")
)

type GetResourceFunc[T any] func(ctx context.Context, id string) (*T, *edgecloud.Response, error)

func ResourceIsDeleted[T any](ctx context.Context, getResourceFunc GetResourceFunc[T], id string) error {
	_, resp, err := getResourceFunc(ctx, id)
	if err == nil {
		return errResourceNotDeleted
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return errGetResourceInfo
}

func ResourceIsExist[T any](ctx context.Context, getResourceFunc GetResourceFunc[T], id string) (bool, error) {
	_, resp, err := getResourceFunc(ctx, id)

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound, http.StatusForbidden:
		return false, nil
	default:
		return false, fmt.Errorf("%w, status code: %d, details: %w", errGetResourceInfo, resp.StatusCode, err)
	}
}

func DeleteResourceIfExist(ctx context.Context, client *edgecloud.Client, resource interface{}, resourceID string) error {
	deleteAndWait := func(
		deleter func(ctx context.Context, resourceID string) (*edgecloud.TaskResponse, *edgecloud.Response, error),
	) error {
		task, _, err := deleter(ctx, resourceID)
		if err != nil {
			return err
		}
		return WaitForTaskComplete(ctx, client, task.Tasks[0])
	}

	switch v := resource.(type) {
	case edgecloud.LoadbalancersService:
		if err := deleteAndWait(v.Delete); err != nil {
			return err
		}
		return ResourceIsDeleted(ctx, v.Get, resourceID)
	case edgecloud.FloatingIPsService:
		if err := deleteAndWait(v.Delete); err != nil {
			return err
		}
		return ResourceIsDeleted(ctx, v.Get, resourceID)
	case edgecloud.VolumesService:
		if err := deleteAndWait(v.Delete); err != nil {
			return err
		}
		return ResourceIsDeleted(ctx, v.Get, resourceID)
	case edgecloud.L7PoliciesService:
		if err := deleteAndWait(v.Delete); err != nil {
			return err
		}
		return ResourceIsDeleted(ctx, v.Get, resourceID)
	case edgecloud.SnapshotsService:
		if err := deleteAndWait(v.Delete); err != nil {
			return err
		}
		return ResourceIsDeleted(ctx, v.Get, resourceID)
	default:
		return errDeleteResourceIfExistIsNotSupported
	}
}
