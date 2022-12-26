package testing

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/testhelper"
)

func createClient() *edgecloud.ServiceClient {
	return &edgecloud.ServiceClient{
		ProviderClient: &edgecloud.ProviderClient{AccessTokenID: "abc123"},
		Endpoint:       testhelper.Endpoint(),
	}
}
