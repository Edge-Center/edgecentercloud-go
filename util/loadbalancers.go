package util

import (
	"github.com/mitchellh/mapstructure"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

type TaskResult struct {
	DdosProfiles   []int    `json:"ddos_profiles"`
	FloatingIPs    []string `json:"floatingips"`
	HealthMonitors []string `json:"healthmonitors"`
	Images         []string `json:"images"`
	Instances      []string `json:"instances"`
	Listeners      []string `json:"listeners"`
	LoadBalancers  []string `json:"loadbalancers"`
	Members        []string `json:"members"`
	Networks       []string `json:"networks"`
	Pools          []string `json:"pools"`
	Ports          []string `json:"ports"`
	Routers        []string `json:"routers"`
	Secrets        []string `json:"secrets"`
	Snapshots      []string `json:"snapshots"`
	Subnets        []string `json:"subnets"`
	Volumes        []string `json:"volumes"`
}

func ExtractTaskResultFromTask(task *edgecloud.Task) (*TaskResult, error) {
	var result TaskResult
	if err := mapstructure.Decode(task.CreatedResources, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
