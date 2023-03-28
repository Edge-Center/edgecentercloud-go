package testing

import (
	"fmt"
	"testing"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/floatingips/v1/client"
	etest "github.com/Edge-Center/edgecentercloud-go/client/testing"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/floatingip/v1/floatingips"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/utils/metadata"
)

const (
	FloatingIPCreateTimeout = 1200
)

func createTestFloatingIP(client *edgecloud.ServiceClient, opts floatingips.CreateOpts) (string, error) {
	res, err := floatingips.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	floatingIPID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, FloatingIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		floatingIPID, err := floatingips.ExtractFloatingIPIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve FloatingIP ID from task info: %w", err)
		}
		return floatingIPID, nil
	})
	if err != nil {
		return "", err
	}
	return floatingIPID.(string), nil
}

func TestFipsMetadata(t *testing.T) {
	resourceName := "floatingip"
	args := []string{"edgeclient", resourceName}
	a, ctx := etest.InitTestApp(args)

	clientFip, err := client.NewFloatingIPClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := floatingips.CreateOpts{}

	resourceID, err := createTestFloatingIP(clientFip, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer floatingips.Delete(clientFip, resourceID)

	err = etest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := floatingips.Get(clientFip, resourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, resourceID)

	if err != nil {
		t.Fatal(err)
	}
}
