# Edgecloud

Edgecloud is a Go client library for accessing the Edgecenter Cloud API.

You can view Edgecenter Cloud API docs here: [https://apidocs.edgecenter.ru/cloud](https://apidocs.edgecenter.ru/cloud)

# Versions
| Version | Supported? | Support expiration | Edgecloud CLI (ec_client) | How to use                                                                                                                                               | Notes                                                   |
|---------|------------|--------------------|---------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------|
| master  | ✅          | _                  | :x:                       | `import "github.com/Edge-Center/edgecentercloud-go/v2"`  and choose version at go.mod file as commit sha like `b193c6019f9a196442db420ac20644772c064c65` | New features and bug fixes arrive here first            |
| v2      | ✅          | _                  | :x:                       | `import "github.com/Edge-Center/edgecentercloud-go/v2"`  and choose version at go.mod file as release version like `v2.X.Y`                              | Used for stable releases                                |                         
| v1      | ❌          | 05.03.2024         | ✅                         | `import "github.com/Edge-Center/edgecentercloud-go"`     and choose version at go.mod file as release version like `v1.X.Y`                              | Not recommended. Use only for edgecloud ec_client usage |   

## Install
```sh
go get github.com/Edge-Center/edgecentercloud-go/v2@vX.Y.Z
```

where X.Y.Z is the [version](https://github.com/Edge-Center/edgecentercloud-go/releases) you need.

or
```sh
go get github.com/Edge-Center/edgecentercloud-go/v2
```
for non Go modules usage or latest version.

## Usage

```go
import edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
```

Create a new EdgeCloud client, then use the exposed services to
access different parts of the Edgecenter Cloud API.

### Authentication

Currently, permanent api-key is the only method of authenticating with the API.
You can find more information about api-key in the [knowledge base](https://support.edgecenter.ru/knowledge_base/item/257788).

You can then use your api-key to create a new client. 
Additionally, you need to set the base URL for the API, Region and Project 

```go
package main

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

func main() {
	cloud, err := edgecloud.NewWithRetries(nil,
		edgecloud.SetAPIKey("<api-key>"),
		edgecloud.SetBaseURL("<base-url>"),  // e.g. "https://api.edgecenter.online/cloud" (string)
		edgecloud.SetRegion(10),             // e.g. 10 (int)
		edgecloud.SetProject(12345),         // e.g. 12345 (int)
	)
	if err != nil {
		// error processing 
    }
}
```

## Examples

To create a new Security group:

```go
securityGroupCreateRequest := &edgecloud.SecurityGroupCreateRequest{
    SecurityGroup: edgecloud.SecurityGroupCreateRequestInner{
        Name: "secgroup",
        SecurityGroupRules: []edgecloud.SecurityGroupRuleCreateRequest{
            {
                Direction:    edgecloud.SGRuleDirectionIngress,
                EtherType:    edgecloud.EtherTypeIPv4,
                Protocol:     edgecloud.SGRuleProtocolTCP,
                PortRangeMin: 10250,
                PortRangeMax: 10259,
            },
        },
    },
}

ctx := context.TODO()

securityGroup, _, err := cloud.SecurityGroups.Create(ctx, securityGroupCreateRequest)
if err != nil {
    // error processing 
}
```

### Create with task response

The creation of some resources does not occur immediately; 
first, a task is launched that needs to be processed.

example 1, when you only need to wait for the task to complete
```go
import "github.com/Edge-Center/edgecentercloud-go/util"

task, _, err := cloud.Floatingips.Create(ctx, &edgecloud.FloatingIPCreateRequest{})
if err != nil {
    // error processing 
}

if err = util.WaitForTaskComplete(ctx, cloud, task.Tasks[0]); err != nil {
    // error processing 
}
```

example 2, when you need to get the id of the created resource
```go
import "github.com/Edge-Center/edgecentercloud-go/util"

task, _, err := cloud.Floatingips.Create(ctx, &edgecloud.FloatingIPCreateRequest{})
if err != nil {
    // error processing 
}

taskInfo, err := util.WaitAndGetTaskInfo(ctx, cloud, task.Tasks[0])
if err != nil {
    // error processing 
}

taskResult, err := util.ExtractTaskResultFromTask(taskInfo)
if err != nil {
    // error processing 
}

fipID := taskResult.FloatingIPs[0]
```

example 3, when you need to get the id of the created resource but using helper method (for some resources)
```go
import "github.com/Edge-Center/edgecentercloud-go/util"

opts := &edgecloud.FloatingIPCreateRequest{}
taskResult, err := util.ExecuteAndExtractTaskResult(ctx, client.Floatingips.Create, opts, cloud)
if err != nil {
    // error processing 
}

fipID := taskResult.FloatingIPs[0]
```

### Helpers
You can find other helpers that extend the api using `util` package

for example, wait for Loadbalancer provisioning status is ACTIVE
```go
loadBalancerID := "..."
attempts := uint(8)

if err := util.WaitLoadbalancerProvisioningStatusActive(ctx, cloud, loadBalancerID, &attempts); err != nil {
    // error processing 
}
```

or, get volumes list by name
```go
volumeName := "my-awesome-volume"

listVolumes, _ := util.VolumesListByName(ctx, cloud, volumeName)
```

or, check that the resource has been deleted
```go
loadBalancerID := "..."
if err := util.ResourceIsDeleted(ctx, cloud.Loadbalancers.Get, loadBalancerID); err != nil {
    // error processing 
}
```
and others helpers

### How to run tests 
```
make test
```
