package client

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter"
	"github.com/Edge-Center/edgecentercloud-go/testhelper"
)

// Fake token to use.
const (
	TokenID      = "cbc36478b0bd8e67e89469c7749d4127"
	AccessToken  = "cbc36478b0bd8e67e89469c7749d4127"
	RefreshToken = "tbc36478b0bd8e67e89469c7749d4127"
	Username     = "username"
	Password     = "password"
	RegionID     = 1
	ProjectID    = 1
)

// ServiceClient returns a generic service client for use in tests.
func ServiceClient() *edgecloud.ServiceClient {
	return &edgecloud.ServiceClient{
		ProviderClient: &edgecloud.ProviderClient{
			AccessTokenID:  AccessToken,
			RefreshTokenID: RefreshToken,
		},
		Endpoint: testhelper.Endpoint(),
	}
}

func ServiceTokenClient(name string, version string) *edgecloud.ServiceClient {
	options := edgecloud.TokenOptions{
		APIURL:       testhelper.Endpoint(),
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		AllowReauth:  true,
	}
	endpointOpts := edgecloud.EndpointOpts{
		Name:    name,
		Region:  RegionID,
		Project: ProjectID,
		Version: version,
	}
	client, err := edgecenter.TokenClientService(options, endpointOpts)
	if err != nil {
		panic(err)
	}

	return client
}

func ServiceAuthClient(name string, version string) *edgecloud.ServiceClient {
	options := edgecloud.AuthOptions{
		APIURL:      testhelper.Endpoint(),
		AuthURL:     testhelper.ECRefreshTokenIdentifyEndpoint(),
		Username:    Username,
		Password:    Password,
		AllowReauth: true,
	}
	endpointOpts := edgecloud.EndpointOpts{
		Name:    name,
		Region:  RegionID,
		Project: ProjectID,
		Version: version,
	}
	client, err := edgecenter.AuthClientService(options, endpointOpts)
	if err != nil {
		panic(err)
	}

	return client
}

type AuthResultTest struct {
	accessToken  string
	refreshToken string
}

func (ar AuthResultTest) ExtractAccessToken() (string, error) {
	return ar.accessToken, nil
}

func (ar AuthResultTest) ExtractRefreshToken() (string, error) {
	return ar.accessToken, nil
}

func (ar AuthResultTest) ExtractTokensPair() (string, string, error) {
	return ar.accessToken, ar.refreshToken, nil
}

func NewAuthResultTest(accessToken string, refreshToken string) AuthResultTest {
	return AuthResultTest{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}
