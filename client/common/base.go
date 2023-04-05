package common

import (
	"github.com/urfave/cli/v2"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/client/flags"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter"
)

func buildTokenClient(c *cli.Context, endpointName, endpointType string, version string) (*edgecloud.ServiceClient, error) {
	settings, err := edgecenter.NewEdgeCloudTokenAPISettingsFromEnv()
	if err != nil {
		return nil, err
	}

	accessToken := c.String("access")
	if accessToken != "" {
		settings.AccessToken = accessToken
	}

	refreshToken := c.String("refresh")
	if refreshToken != "" {
		settings.RefreshToken = refreshToken
	}

	if version == "" {
		version = c.String("api-version")
	}
	if version != "" {
		settings.Version = version
	}

	url := c.String("api-url")
	if url != "" {
		settings.APIURL = url
	}

	region := c.Int("region")
	if region != 0 {
		settings.Region = region
	}

	project := c.Int("project")
	if project != 0 {
		settings.Project = project
	}

	debug := c.Bool("debug")
	if debug {
		settings.Debug = true
	}

	settings.Name = endpointName
	settings.Type = endpointType

	err = settings.Validate()
	if err != nil {
		return nil, err
	}

	options := settings.ToTokenOptions()
	eo := settings.ToEndpointOptions()

	return edgecenter.TokenClientServiceWithDebug(options, eo, settings.Debug)
}

func buildAPITokenClient(c *cli.Context, endpointName, endpointType string, version string) (*edgecloud.ServiceClient, error) {
	settings, err := edgecenter.NewECCloudAPITokenAPISettingsFromEnv()
	if err != nil {
		return nil, err
	}

	apiToken := c.String("api-token")
	if apiToken != "" {
		settings.APIToken = apiToken
	}

	if version == "" {
		version = c.String("api-version")
	}
	if version != "" {
		settings.Version = version
	}

	url := c.String("api-url")
	if url != "" {
		settings.APIURL = url
	}

	region := c.Int("region")
	if region != 0 {
		settings.Region = region
	}

	project := c.Int("project")
	if project != 0 {
		settings.Project = project
	}

	debug := c.Bool("debug")
	if debug {
		settings.Debug = true
	}

	settings.Name = endpointName
	settings.Type = endpointType

	err = settings.Validate()
	if err != nil {
		return nil, err
	}

	options := settings.ToAPITokenOptions()
	eo := settings.ToEndpointOptions()

	return edgecenter.APITokenClientServiceWithDebug(options, eo, settings.Debug)
}

func buildPlatformClient(c *cli.Context, endpointName, endpointType string, version string) (*edgecloud.ServiceClient, error) {
	settings, err := edgecenter.NewECCloudPlatformAPISettingsFromEnv()
	if err != nil {
		return nil, err
	}

	username := c.String("username")
	if username != "" {
		settings.Username = username
	}

	password := c.String("password")
	if password != "" {
		settings.Password = password
	}

	if version == "" {
		version = c.String("api-version")
	}
	if version != "" {
		settings.Version = version
	}

	url := c.String("api-url")
	if url != "" {
		settings.APIURL = url
	}

	region := c.Int("region")
	if region != 0 {
		settings.Region = region
	}

	project := c.Int("project")
	if project != 0 {
		settings.Project = project
	}

	debug := c.Bool("debug")

	if debug {
		settings.Debug = true
	}

	settings.Name = endpointName
	settings.Type = endpointType

	err = settings.Validate()
	if err != nil {
		return nil, err
	}

	options := settings.ToAuthOptions()
	eo := settings.ToEndpointOptions()

	return edgecenter.AuthClientServiceWithDebug(options, eo, settings.Debug)
}

func BuildClient(c *cli.Context, endpointName, version string) (*edgecloud.ServiceClient, error) {
	clientType := flags.ClientType
	if clientType == "" {
		clientType = c.String("client-type")
	}

	switch clientType {
	case "token":
		return buildTokenClient(c, endpointName, "", version)
	case "api-token":
		return buildAPITokenClient(c, endpointName, "", version)
	default:
		return buildPlatformClient(c, endpointName, "", version)
	}
}

func BuildAPITokenClient(ao edgecloud.AuthOptions) (*edgecloud.ServiceClient, error) {
	provider, err := edgecenter.AuthenticatedClient(ao)
	if err != nil {
		return nil, err
	}
	return edgecenter.NewIdentity(provider, edgecloud.EndpointOpts{})
}
