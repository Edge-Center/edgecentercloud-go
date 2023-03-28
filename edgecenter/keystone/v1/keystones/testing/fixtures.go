package testing

import (
	"time"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/keystones"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/types"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "state": "NEW",
      "id": 1,
      "keystone_federated_domain_id": "5ac0a17e556d4a9c8f946928a7953990",
      "admin_password": "******",
      "url": "https://ed-10.cloud.core.pw:5000/v3",
      "created_on": "2020-04-10T11:37:58"
    }
  ]
}
`

const GetResponse = `
{
  "state": "NEW",
  "id": 1,
  "keystone_federated_domain_id": "5ac0a17e556d4a9c8f946928a7953990",
  "admin_password": "******",
  "url": "https://ed-10.cloud.core.pw:5000/v3",
  "created_on": "2020-04-10T11:37:58"
}
`

const CreateRequest = `
{
  "state": "NEW",
  "keystone_federated_domain_id": "5ac0a17e556d4a9c8f946928a7953990",
  "url": "https://ed-10.cloud.core.pw:5000/v3"
}
`

const UpdateRequest = `
{
  "state": "DELETED",
  "url": "https://ed-10.cloud.core.pw:5000/v3"
}	
`

const (
	CreateResponse = GetResponse
	UpdateResponse = GetResponse
)

var (
	createdTimeString    = "2020-04-10T11:37:58"
	createdTimeParsed, _ = time.Parse(edgecloud.RFC3339NoZ, createdTimeString)
	createdTime          = edgecloud.JSONRFC3339NoZ{Time: createdTimeParsed}
	keystoneURL, _       = edgecloud.ParseURL("https://ed-10.cloud.core.pw:5000/v3")

	Keystone1 = keystones.Keystone{
		ID:                        1,
		URL:                       *keystoneURL,
		State:                     types.KeystoneStateNew,
		KeystoneFederatedDomainID: "5ac0a17e556d4a9c8f946928a7953990",
		CreatedOn:                 createdTime,
		AdminPassword:             "******",
	}

	ExpectedKeystoneSlice = []keystones.Keystone{Keystone1}
)
