package testing

import (
	"fmt"
	"testing"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/loadbalancers/v1/client"
	etest "github.com/Edge-Center/edgecentercloud-go/client/testing"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/loadbalancer/v1/loadbalancers"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/loadbalancer/v1/types"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/utils/metadata"
)

const (
	LoadBalancerCreateTimeout = 2400
	lbTestName                = "test-lb"
	lbListenerTestName        = "test-listener"
)

func createTestLoadBalancerWithListener(client *edgecloud.ServiceClient, opts loadbalancers.CreateOpts) (string, error) {
	res, err := loadbalancers.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	lbID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, LoadBalancerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		lbID, err := loadbalancers.ExtractLoadBalancerIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve LoadBalancer ID from task info: %w", err)
		}
		return lbID, nil
	})
	if err != nil {
		return "", err
	}

	return lbID.(string), nil
}

func TestLBSMetadata(t *testing.T) {
	resourceName := "loadbalancer"
	args := []string{"edgeclient", resourceName}
	a, ctx := etest.InitTestApp(args)

	clientLb, err := client.NewLoadbalancerClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := loadbalancers.CreateOpts{
		Name: lbTestName,
		Listeners: []loadbalancers.CreateListenerOpts{{
			Name:         lbListenerTestName,
			ProtocolPort: 80,
			Protocol:     types.ProtocolTypeHTTP,
		}},
	}

	resourceID, err := createTestLoadBalancerWithListener(clientLb, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer loadbalancers.Delete(clientLb, resourceID)

	err = etest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := loadbalancers.Get(clientLb, resourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, resourceID)

	if err != nil {
		t.Fatal(err)
	}
}
