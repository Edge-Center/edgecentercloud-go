package testing

import (
	"errors"
	"fmt"
	"net"
	"testing"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	netclient "github.com/Edge-Center/edgecentercloud-go/client/networks/v1/client"
	"github.com/Edge-Center/edgecentercloud-go/client/subnets/v1/client"
	etest "github.com/Edge-Center/edgecentercloud-go/client/testing"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/network/v1/networks"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/subnet/v1/subnets"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/utils/metadata"
)

const (
	SubnetDeleting        int = 1200
	SubnetCreatingTimeout int = 1200
)

func createTestSubnet(client *edgecloud.ServiceClient, opts subnets.CreateOpts, subCidr string) (string, error) {
	var gccidr edgecloud.CIDR
	_, netIPNet, err := net.ParseCIDR(subCidr)
	if err != nil {
		return "", err
	}
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask
	opts.CIDR = gccidr

	res, err := subnets.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	subnetID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, SubnetCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		Subnet, err := subnets.ExtractSubnetIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Subnet ID from task info: %w", err)
		}
		return Subnet, nil
	},
	)

	return subnetID.(string), err
}

func deleteTestSubnet(client *edgecloud.ServiceClient, subnetID string) error {
	results, err := subnets.Delete(client, subnetID).Extract()
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, SubnetDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := subnets.Get(client, subnetID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete subnet with ID: %s", subnetID)
		}
		var e edgecloud.Default404Error
		if errors.As(err, &e) {
			return nil, nil
		}
		return nil, err
	})

	return err
}

func TestSubnetsMetadata(t *testing.T) {
	resourceName := "subnet"

	args := []string{"edgeclient", resourceName}
	a, ctx := etest.InitTestApp(args)

	netClient, err := netclient.NewNetworkClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	resourceClient, err := client.NewSubnetClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := networks.CreateOpts{
		Name: "test-network1",
	}

	networkID, err := etest.CreateTestNetwork(netClient, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer func(client *edgecloud.ServiceClient, networkID string) {
		err := etest.DeleteTestNetwork(client, networkID)
		if err != nil {
			t.Errorf("error while network delete: %s", err.Error())
		}
	}(netClient, networkID)

	optsSubnet := subnets.CreateOpts{
		Name:      "test-subnet",
		NetworkID: networkID,
	}

	resourceID, err := createTestSubnet(resourceClient, optsSubnet, "192.168.42.0/24")
	if err != nil {
		t.Fatal(err)
	}

	defer func(client *edgecloud.ServiceClient, subnetID string) {
		err := deleteTestSubnet(client, subnetID)
		if err != nil {
			t.Errorf("error while subnet delete: %s", err.Error())
		}
	}(resourceClient, resourceID)

	err = etest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := subnets.Get(resourceClient, resourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, resourceID)

	if err != nil {
		t.Fatal(err)
	}
}
