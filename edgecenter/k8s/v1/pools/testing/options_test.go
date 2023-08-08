package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/k8s/v1/pools"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
)

func TestUpdateOpts(t *testing.T) {
	options := pools.UpdateOpts{
		MinNodeCount: 5,
		MaxNodeCount: 3,
	}

	_, err := options.ToClusterPoolUpdateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")

	options = pools.UpdateOpts{
		MinNodeCount: 5,
		MaxNodeCount: 3,
	}

	_, err = options.ToClusterPoolUpdateMap()

	require.Error(t, err)

	options = pools.UpdateOpts{}

	_, err = options.ToClusterPoolUpdateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "MinNodeCount")
	require.Contains(t, err.Error(), "Name")
}

func TestCreateOpts(t *testing.T) {
	nodeCount1 := 0
	dockerVolumeSize1 := 0
	maxNodeCount1 := 3

	options := pools.CreateOpts{
		Name:             "",
		FlavorID:         "",
		NodeCount:        &nodeCount1,
		DockerVolumeSize: &dockerVolumeSize1,
		MinNodeCount:     5,
		MaxNodeCount:     &maxNodeCount1,
	}

	_, err := options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "FlavorID")
	require.Contains(t, err.Error(), "NodeCount")

	nodeCount2 := 5
	dockerVolumeSize2 := 10
	maxNodeCount2 := 3

	options = pools.CreateOpts{
		Name:             "name",
		FlavorID:         "flavor",
		NodeCount:        &nodeCount2,
		DockerVolumeSize: &dockerVolumeSize2,
		MinNodeCount:     4,
		MaxNodeCount:     &maxNodeCount2,
	}

	_, err = options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MaxNodeCount")
	require.Contains(t, err.Error(), "MinNodeCount")

	nodeCount3 := 5
	dockerVolumeSize3 := 10
	maxNodeCount3 := 8

	options = pools.CreateOpts{
		Name:             "name",
		FlavorID:         "flavor",
		NodeCount:        &nodeCount3,
		DockerVolumeSize: &dockerVolumeSize3,
		MinNodeCount:     6,
		MaxNodeCount:     &maxNodeCount3,
	}

	_, err = options.ToClusterPoolCreateMap()

	require.Error(t, err)
	require.Contains(t, err.Error(), "MinNodeCount")
	require.Contains(t, err.Error(), "NodeCount")
}

func TestDecodePoolTask(t *testing.T) {
	taskID := "732851e1-f792-4194-b966-4cbfa5f30093"
	rs := map[string]interface{}{"k8s_pools": []string{taskID}}
	taskInfo := tasks.Task{
		CreatedResources: &rs,
	}
	var result pools.PoolTaskResult
	err := edgecloud.NativeMapToStruct(taskInfo.CreatedResources, &result)
	require.NoError(t, err)
	require.Equal(t, taskID, result.K8sPools[0])
}
