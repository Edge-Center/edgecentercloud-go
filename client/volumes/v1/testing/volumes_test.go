package testing

import (
	"fmt"
	"testing"

	etest "github.com/Edge-Center/edgecentercloud-go/client/testing"
	"github.com/Edge-Center/edgecentercloud-go/client/volumes/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/utils/metadata"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/volume/v1/volumes"
)

func TestVolumeMetadata(t *testing.T) {
	resourceName := "volume"
	args := []string{"edgeclient", resourceName}
	a, ctx := etest.InitTestApp(args)

	clientVolume, err := client.NewVolumeClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := volumes.CreateOpts{
		Name:     "test-volume-1",
		Size:     1,
		Source:   volumes.NewVolume,
		TypeName: volumes.Standard,
	}

	res, err := volumes.Create(clientVolume, opts).Extract()
	if err != nil {
		t.Fatal(err)
	}

	taskID := res.Tasks[0]
	resourceID, err := tasks.WaitTaskAndReturnResult(clientVolume, taskID, true, 1200, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(clientVolume, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		resourceID, err := volumes.ExtractVolumeIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve volume ID from task info: %w", err)
		}
		return resourceID, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	typedResourceID := resourceID.(string)
	defer volumes.Delete(clientVolume, typedResourceID, volumes.DeleteOpts{})

	err = etest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := volumes.Get(clientVolume, typedResourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, typedResourceID)

	if err != nil {
		t.Fatal(err)
	}
}
