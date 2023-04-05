package testing

import (
	"time"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/project/v1/projects"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/project/v1/types"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "id": 1,
      "state": "ACTIVE",
      "created_at": "2020-04-10T11:37:57",
      "description": "",
      "client_id": 1,
      "name": "default"
    }
  ]
}
`

const GetResponse = `
{
  "id": 1,
  "state": "ACTIVE",
  "created_at": "2020-04-10T11:37:57",
  "description": "",
  "client_id": 1,
  "name": "default"
}
`

const CreateRequest = `
{
	"client_id": 1,
	"state": "ACTIVE",
	"name": "default"
}
`

const UpdateRequest = `
{
	"name": "default",
	"description": "description"
}	
`

const (
	CreateResponse = GetResponse
	UpdateResponse = GetResponse
)

var (
	createdTimeString    = "2020-04-10T11:37:57"
	createdTimeParsed, _ = time.Parse(edgecloud.RFC3339NoZ, createdTimeString)
	createdTime          = edgecloud.JSONRFC3339NoZ{Time: createdTimeParsed}

	Project1 = projects.Project{
		ID:          1,
		ClientID:    1,
		Name:        "default",
		Description: "",
		State:       types.ProjectStateActive,
		TaskID:      nil,
		CreatedAt:   createdTime,
	}

	ExpectedProjectSlice = []projects.Project{Project1}
)
