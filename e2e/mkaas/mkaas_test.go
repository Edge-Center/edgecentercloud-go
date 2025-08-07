//go:build e2e
// +build e2e

package cloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
	"github.com/Edge-Center/edgecentercloud-go/v2/util"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
)

const (
	clusterCreateTimeout = 15 * time.Minute
	poolCreateTimeout    = 10 * time.Minute
	poolUpdateTimeout    = 10 * time.Minute
	poolDeleteTimeout    = 3 * time.Minute
	mkaaSE2EProjectID    = 1350231
	mkaaSE2ENetworkID    = "62179813-5c87-4637-ad0e-7aba92cd3f70"
	mkaaSE2ESubnetID     = "b0a55566-aaff-46df-abad-e3506bf879e5"
	mkaaSE2EKeypairName  = "129662_test_cloud_common_edgecenter_ru_1753109908"
	mkaaSE2EFlavor       = "g1-standard-2-4"
)

var (
	randomString, _      = util.GenerateRandomString(6, true, true, false)
	clusterCreateRequest = edgecloud.MkaaSClusterCreateRequest{
		Name:           fmt.Sprintf("Cluster_from_go_client_%s", randomString),
		SSHKeyPairName: mkaaSE2EKeypairName,
		NetworkID:      mkaaSE2ENetworkID,
		SubnetID:       mkaaSE2ESubnetID,
		ControlPlane: edgecloud.ControlPlaneCreateRequest{
			Flavor:     mkaaSE2EFlavor,
			NodeCount:  1,
			VolumeSize: 10,
			VolumeType: "",
			Version:    "v1.31.0",
		},
		Pools: []edgecloud.MkaaSPoolCreateRequest{
			{
				Name:         fmt.Sprintf("pool_from_go_client_%s", randomString),
				Flavor:       mkaaSE2EFlavor,
				MaxNodeCount: edgecloud.PtrTo(5),
				MinNodeCount: edgecloud.PtrTo(1),
				NodeCount:    2,
				VolumeSize:   10,
				VolumeType:   "",
				Labels:       nil,
				Taints:       nil,
			},
		}}
	poolCreateRequest = edgecloud.MkaaSPoolCreateRequest{
		Name:            "pool-create",
		Flavor:          mkaaSE2EFlavor,
		NodeCount:       3,
		VolumeSize:      20,
		SecurityGroupID: nil,
		VolumeType:      "",
		Labels:          nil,
		Taints:          nil,
	}

	increaseNodeCountRequest = edgecloud.MkaaSPoolUpdateRequest{
		NodeCount:    edgecloud.PtrTo(4),
		MinNodeCount: edgecloud.PtrTo(1),
		MaxNodeCount: edgecloud.PtrTo(5),
	}

	decreaseNodeCountRequest = edgecloud.MkaaSPoolUpdateRequest{
		NodeCount:    edgecloud.PtrTo(1),
		MinNodeCount: edgecloud.PtrTo(1),
		MaxNodeCount: edgecloud.PtrTo(5),
	}
)

type MkaasSuite struct {
	suite.Suite
	cloudAPIUrl   string
	cloudAPIToken string
	clusterID     int
	client        *edgecloud.Client
}

func Test_MkaaS(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(MkaasSuite))
}

func (s *MkaasSuite) Test_CreateClusterSimple() {
	ctx := context.Background()
	client := s.client
	s.createClusterAndWait(ctx, clusterCreateRequest)

	cluster, resp, err := client.MkaaS.ClusterGet(ctx, s.clusterID)
	s.Require().NotNil(cluster)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	// s.Assert().Equal(clusterCreateRequest.Name, cluster.Name)
	s.Assert().Equal(clusterCreateRequest.NetworkID, cluster.NetworkID)
	s.Assert().Equal(clusterCreateRequest.SubnetID, cluster.SubnetID)
	// s.Assert().Equal(clusterCreateRequest.SSHKeyPairName, cluster.SSHKeypairName)
	s.Assert().Equal(clusterCreateRequest.ControlPlane.NodeCount, cluster.ControlPlane.NodeCount)
	s.Assert().Equal(clusterCreateRequest.ControlPlane.Flavor, cluster.ControlPlane.Flavor)
	// s.Assert().Equal(clusterCreateRequest.ControlPlane.VolumeSize, cluster.ControlPlane.VolumeSize)
	s.Assert().Equal(edgecloud.VolumeTypeSsdHiIops, cluster.ControlPlane.VolumeType)
	s.Assert().Equal(clusterCreateRequest.ControlPlane.Version, cluster.ControlPlane.Version)
	s.Assert().Equal(len(clusterCreateRequest.Pools), len(clusterCreateRequest.Pools))
	s.Assert().Equal(clusterCreateRequest.Pools[0].Name, cluster.Pools[0].Name)
	s.Assert().Equal(clusterCreateRequest.Pools[0].NodeCount, cluster.Pools[0].NodeCount)
	// s.Assert().Equal(clusterCreateRequest.Pools[0].VolumeSize, cluster.Pools[0].VolumeSize)
	s.Assert().Equal(edgecloud.VolumeTypeSsdHiIops, cluster.Pools[0].VolumeType)
	s.Assert().Equal(clusterCreateRequest.Pools[0].Flavor, cluster.Pools[0].Flavor)
}

func (s *MkaasSuite) Test_IncreasePoolsNodeCountSimple() {
	ctx := context.Background()
	if s.clusterID == 0 {
		s.createClusterAndWait(ctx, clusterCreateRequest)
	}
	s.changePoolsNodeCountTest(ctx, increaseNodeCountRequest)
}

func (s *MkaasSuite) Test_DecreasePoolsNodeCountSimple() {
	ctx := context.Background()
	if s.clusterID == 0 {
		s.createClusterAndWait(ctx, clusterCreateRequest)
	}
	s.changePoolsNodeCountTest(ctx, decreaseNodeCountRequest)
}

func (s *MkaasSuite) Test_CreateAndDeletePoolSimple() {
	ctx := context.Background()
	if s.clusterID == 0 {
		s.createClusterAndWait(ctx, clusterCreateRequest)
	}
	// s.clusterID = 87
	client := s.client

	s.T().Log("Pool creating started")
	taskResponse, resp, err := client.MkaaS.PoolCreate(ctx, s.clusterID, poolCreateRequest)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.waitTaskWithLoggingAndSaveClusterID(ctx, taskResponse, poolCreateTimeout, clusterCreateRequest.Name)
	cluster, resp, err := client.MkaaS.ClusterGet(ctx, s.clusterID)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.Assert().Equal(len(clusterCreateRequest.Pools)+1, len(cluster.Pools))
	createdPool, found := lo.Find(cluster.Pools, func(pool edgecloud.MkaaSPool) bool {
		return pool.Name == poolCreateRequest.Name
	})
	s.Require().Equal(true, found)
	s.Assert().Equal(poolCreateRequest.NodeCount, createdPool.NodeCount)
	s.Assert().Equal(poolCreateRequest.Flavor, createdPool.Flavor)
	s.Assert().Equal(poolCreateRequest.VolumeSize, createdPool.VolumeSize)
	s.deletePoolAndWait(ctx, createdPool.ID)
}

func (s *MkaasSuite) SetupSuite() {
	cloudApiToken := os.Getenv("EC_E2E_TEST_APIKEY")
	if cloudApiToken == "" {
		err := godotenv.Load("../.env")
		s.Require().NoError(err)
	}
	cloudApiToken = os.Getenv("EC_E2E_TEST_APIKEY")
	cloudApiURL := os.Getenv("EC_E2E_TEST_BASE_URL")
	if cloudApiToken == "" || cloudApiURL == "" {
		s.Fail("EC_E2E_TEST_APIKEY or EC_E2E_TEST_BASE_URL  env is not found")
	}
	client, err := edgecloud.NewWithRetries(nil,
		edgecloud.SetAPIKey(cloudApiToken),
		edgecloud.SetBaseURL(cloudApiURL),
		edgecloud.SetRegion(8),
		edgecloud.SetProject(mkaaSE2EProjectID))
	s.Require().NoError(err)
	s.client = client
	s.cloudAPIUrl = cloudApiURL
	s.cloudAPIToken = cloudApiToken
}

func (s *MkaasSuite) TearDownSuite() {
	if s.clusterID != 0 {
		s.T().Helper()
		s.T().Logf("Deleting cluster %d", s.clusterID)
		s.deleteCluster(s.clusterID)
	}
}

func (s *MkaasSuite) deleteCluster(clusterID int) {
	client := s.client
	_, resp, err := client.MkaaS.ClusterDelete(context.Background(), clusterID)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
}

func (s *MkaasSuite) createClusterAndWait(ctx context.Context, createRequest edgecloud.MkaaSClusterCreateRequest) {
	client := s.client
	s.T().Log("Creating cluster started")
	taskResponse, resp, err := client.MkaaS.ClusterCreate(ctx, createRequest)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.waitTaskWithLoggingAndSaveClusterID(ctx, taskResponse, clusterCreateTimeout, clusterCreateRequest.Name)
	s.T().Logf("created cluster id: %d", s.clusterID)
}

func (s *MkaasSuite) deletePoolAndWait(ctx context.Context, poolID int) {
	client := s.client
	s.T().Log("Deleting pool started")
	taskResponse, resp, err := client.MkaaS.PoolDelete(ctx, s.clusterID, poolID)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.waitTaskWithLoggingAndSaveClusterID(ctx, taskResponse, poolDeleteTimeout, clusterCreateRequest.Name)
}

func (s *MkaasSuite) updatePoolAndWait(ctx context.Context, poolID int, updatedPoolRequest edgecloud.MkaaSPoolUpdateRequest) {
	client := s.client
	s.T().Log("Updating pool started")
	taskResponse, resp, err := client.MkaaS.PoolUpdate(ctx, s.clusterID, poolID, updatedPoolRequest)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.waitTaskWithLoggingAndSaveClusterID(ctx, taskResponse, poolUpdateTimeout, clusterCreateRequest.Name)
}

func (s *MkaasSuite) waitTaskWithLoggingAndSaveClusterID(ctx context.Context, taskResponse *edgecloud.TaskResponse, timeout time.Duration, clusterName string) {
	var taskProcessed bool
	wg := sync.WaitGroup{}
	timeNow := time.Now().UTC()
	wg.Add(1)
	go func() {
		defer wg.Done()
		task, err := util.WaitAndGetTaskInfo(ctx, s.client, taskResponse.Tasks[0], timeout)
		if err != nil {
			taskProcessed = true
			cluster, _, err := s.client.MkaaS.ClusterGet(ctx, s.clusterID)
			s.Assert().NoError(err)
			s.T().Logf("cluster id: %d, cluster status: %s, cluster state: %s, task elapsed time: %.2fmin", cluster.ID, cluster.Status, cluster.State, time.Since(timeNow).Minutes())
		}
		s.Assert().NotNil(task)
		taskProcessed = true
	}()
	for !taskProcessed {
		time.Sleep(5 * time.Second)
		if s.clusterID == 0 {
			clusters, resp, err := s.client.MkaaS.ClustersList(ctx, &edgecloud.MkaaSClusterListOptions{
				Name: clusterName,
			})
			s.Assert().NoError(err)
			s.Assert().Equal(http.StatusOK, resp.StatusCode)
			s.Require().Equal(1, len(clusters))
			s.clusterID = clusters[0].ID
		}
		cluster, _, err := s.client.MkaaS.ClusterGet(ctx, s.clusterID)
		s.Assert().NoError(err)
		s.T().Logf("cluster id: %d, cluster status: %s, cluster state: %s, task elapsed time: %.2fmin", cluster.ID, cluster.Status, cluster.State, time.Since(timeNow).Minutes())
	}
	wg.Wait()
	s.T().Logf("task processed")
}

func (s *MkaasSuite) changePoolsNodeCountTest(ctx context.Context, changeNodeCountRequest edgecloud.MkaaSPoolUpdateRequest) {
	client := s.client
	cluster, resp, err := client.MkaaS.ClusterGet(ctx, s.clusterID)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.Require().LessOrEqual(1, len(cluster.Pools))

	s.T().Log("Pool updating started")
	taskResponse, resp, err := client.MkaaS.PoolUpdate(ctx, s.clusterID, cluster.Pools[0].ID, changeNodeCountRequest)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.waitTaskWithLoggingAndSaveClusterID(ctx, taskResponse, poolUpdateTimeout, clusterCreateRequest.Name)
	cluster, resp, err = client.MkaaS.ClusterGet(ctx, s.clusterID)
	s.Require().NoError(err)
	s.Assert().Equal(http.StatusOK, resp.StatusCode)
	s.Assert().Equal(*changeNodeCountRequest.NodeCount, cluster.Pools[0].NodeCount)
}
