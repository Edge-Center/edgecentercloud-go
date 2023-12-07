package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/k8s/v1/clusters"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/k8s/v1/pools"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
)

func TestResizeOpts(t *testing.T) {
	nodeCount1 := 0
	options := clusters.ResizeOpts{
		NodeCount:     &nodeCount1,
		NodesToRemove: nil,
	}

	_, err := options.ToClusterResizeMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NodeCount")

	nodeCount2 := 1
	options = clusters.ResizeOpts{
		NodeCount:     &nodeCount2,
		NodesToRemove: []string{"1"},
	}

	_, err = options.ToClusterResizeMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "NodesToRemove")

	options = clusters.ResizeOpts{
		NodeCount: &nodeCount2,
	}

	_, err = options.ToClusterResizeMap()
	require.NoError(t, err)
}

func TestCreateOptions(t *testing.T) {
	options := clusters.CreateOpts{
		Name:         Cluster1.Name,
		KeyPair:      "",
		FixedNetwork: "",
		FixedSubnet:  "",
		Version:      "",
		Pools:        []pools.CreateOpts{},
	}

	_, err := options.ToClusterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "FixedNetwork")
	require.Contains(t, err.Error(), "FixedSubnet")
	require.Contains(t, err.Error(), "Pools")

	options = clusters.CreateOpts{
		Name:         Cluster1.Name,
		KeyPair:      "",
		FixedNetwork: fixedNetwork,
		FixedSubnet:  fixedNetwork,
		Version:      "111",
	}

	_, err = options.ToClusterCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Version")
	require.Contains(t, err.Error(), "Pools")
}

func TestDecodeClusterTask(t *testing.T) {
	taskID := "732851e1-f792-4194-b966-4cbfa5f30093"
	rs := map[string]interface{}{"k8s_clusters": []string{taskID}}
	taskInfo := tasks.Task{
		CreatedResources: &rs,
	}
	var result clusters.ClusterTaskResult
	err := edgecloud.NativeMapToStruct(taskInfo.CreatedResources, &result)
	require.NoError(t, err)
	require.Equal(t, taskID, result.K8sClusters[0])
}