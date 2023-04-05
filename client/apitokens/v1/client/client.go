package client

import (
	"github.com/urfave/cli/v2"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/common"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter"
)

func NewAPITokenClient(_ *cli.Context) (*edgecloud.ServiceClient, error) {
	// todo refactor it, now apitokens could be generated only with platform client type
	settings, err := edgecenter.NewECCloudPlatformAPISettingsFromEnv()
	if err != nil {
		return nil, err
	}

	ao, err := edgecenter.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	ao.APIURL = settings.AuthURL

	return common.BuildAPITokenClient(ao)
}
