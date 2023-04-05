package testing

import (
	"testing"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/networks/v1/client"
	etest "github.com/Edge-Center/edgecentercloud-go/client/testing"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/network/v1/networks"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/utils/metadata"
)

func TestNetworksMetadata(t *testing.T) {
	resourceName := "network"
	args := []string{"edgeclient", resourceName}
	a, ctx := etest.InitTestApp(args)

	resourceClient, err := client.NewNetworkClientV1(ctx)
	if err != nil {
		t.Fatal(err)
	}

	opts := networks.CreateOpts{
		Name: "test-network1",
	}

	resourceID, err := etest.CreateTestNetwork(resourceClient, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer func(client *edgecloud.ServiceClient, networkID string) {
		err := etest.DeleteTestNetwork(client, networkID)
		if err != nil {
			t.Errorf("error while network delete: %s", err.Error())
		}
	}(resourceClient, resourceID)

	err = etest.MetadataTest(func() ([]metadata.Metadata, error) {
		res, err := networks.Get(resourceClient, resourceID).Extract()
		if err != nil {
			return nil, err
		}
		return res.Metadata, nil
	}, a, resourceName, resourceID)

	if err != nil {
		t.Fatal(err)
	}
}
