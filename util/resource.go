package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return errGetResourceInfo
}

func ResourceIsExist[T any](ctx context.Context, getResourceFunc GetResourceFunc[T], id string) (bool, error) {
	_, resp, err := getResourceFunc(ctx, id)

	if resp == nil {
		return false, fmt.Errorf("%w, details: %w", errGetResourceInfo, err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound, http.StatusForbidden:
		return false, nil
	default:
		return false, fmt.Errorf("%w, status code: %d, details: %w", errGetResourceInfo, resp.StatusCode, err)
	}
}

func DeleteResourceIfExist(ctx context.Context, client *edgecloud.Client, resource interface{}, resourceID string, timeouts ...time.Duration) error {
	deleteAndWait := func(
		deleter func(ctx context.Context, resourceID string) (*edgecloud.TaskResponse, *edgecloud.Response, error),
	) error {
		task, resp, err := deleter(ctx, resourceID)
		if err != nil {
			if IsNotFoundErr(resp) {
				return nil
			}

			if IsLockedErr(resp) {
				return RetryDeleteLocked(ctx, client, deleter, resourceID, timeouts...)
			}

			return err
		}

		if task == nil || len(task.Tasks) == 0 {
			return nil
		}

		return WaitForTaskComplete(ctx, client, task.Tasks[0], timeouts...)
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

func RetryDeleteLocked(ctx context.Context, client *edgecloud.Client,
	deleter func(ctx context.Context, resourceID string) (*edgecloud.TaskResponse, *edgecloud.Response, error),
	resourceID string,
	timeouts ...time.Duration,
) error {
	timeout := defaultTimeout

	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			task, resp, err := deleter(ctx, resourceID)
			if err != nil {
				if IsNotFoundErr(resp) {
					return nil
				}

				if IsLockedErr(resp) {
					continue
				}

				return err
			}

			if task == nil || len(task.Tasks) == 0 {
				return nil
			}

			return WaitForTaskComplete(ctx, client, task.Tasks[0], timeouts...)
		}
	}
}

func IsNotFoundErr(resp *edgecloud.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

func IsLockedErr(resp *edgecloud.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusConflict
}
