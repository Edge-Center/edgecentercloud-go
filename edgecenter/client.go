package edgecenter

import (
	"fmt"
	"reflect"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/identity/tokens"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/utils"
)

/*
NewECClient prepares an unauthenticated ProviderClient instance.
Most users will probably prefer using the AuthenticatedClient function instead.

This is useful if you wish to explicitly control the version of the identity
service that's used for authentication explicitly, for example.
*/
func NewECClient(endpoint string) (*edgecloud.ProviderClient, error) {
	base, err := utils.BaseRootEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	endpoint = edgecloud.NormalizeURL(endpoint)
	base = edgecloud.NormalizeURL(base)

	p := edgecloud.NewProviderClient()
	p.IdentityBase = base
	p.APIBase = endpoint
	p.IdentityEndpoint = endpoint
	p.UseTokenLock()

	return p, nil
}

/*
AuthenticatedClient logs in to an Edgecenter cloud found at the identity endpoint
specified by the options, acquires a token, and returns a Provider Client
instance that's ready to operate.

Example:

	ao, err := edgecenter.AuthOptionsFromEnv()
	provider, err := edgecenter.AuthenticatedClient(ao)
	client, err := edgecenter.NewMagnumV1(provider, edgecloud.EndpointOpts{})
*/
func AuthenticatedClient(options edgecloud.AuthOptions) (*edgecloud.ProviderClient, error) {
	client, err := NewECClient(options.APIURL)
	if err != nil {
		return nil, err
	}
	err = Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func AuthenticatedClientWithDebug(options edgecloud.AuthOptions, debug bool) (*edgecloud.ProviderClient, error) {
	client, err := AuthenticatedClient(options)
	if err != nil {
		return nil, err
	}
	client.SetDebug(debug)
	return client, err
}

func APITokenClient(options edgecloud.APITokenOptions) (*edgecloud.ProviderClient, error) {
	client, err := NewECClient(options.APIURL)
	if err != nil {
		return nil, err
	}

	if err := client.SetAPIToken(options); err != nil {
		return nil, err
	}
	return client, nil
}

func APITokenClientWithDebug(options edgecloud.APITokenOptions, debug bool) (*edgecloud.ProviderClient, error) {
	client, err := APITokenClient(options)
	if err != nil {
		return nil, err
	}
	client.SetDebug(debug)
	return client, err
}

func TokenClient(options edgecloud.TokenOptions) (*edgecloud.ProviderClient, error) {
	client, err := NewECClient(options.APIURL)
	if err != nil {
		return nil, err
	}
	err = client.SetTokensAndAuthResult(options)
	if err != nil {
		return nil, err
	}
	setECCloudReauth(client, "", options, edgecloud.EndpointOpts{})
	return client, nil
}

func TokenClientWithDebug(options edgecloud.TokenOptions, debug bool) (*edgecloud.ProviderClient, error) {
	client, err := TokenClient(options)
	if err != nil {
		return nil, err
	}
	client.SetDebug(debug)
	return client, err
}

// Authenticate or re-authenticate against the most recent identity service supported at the provided endpoint.
func Authenticate(client *edgecloud.ProviderClient, options edgecloud.AuthOptions) error {
	return auth(client, options.AuthURL, options, edgecloud.EndpointOpts{})
}

func auth(client *edgecloud.ProviderClient, endpoint string, options edgecloud.AuthOptions, eo edgecloud.EndpointOpts) error {

	identityClient, err := NewIdentity(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		identityClient.Endpoint = edgecloud.NormalizeURL(endpoint)
	}

	result := tokens.Create(identityClient, options)

	err = client.SetTokensAndAuthResult(result)
	if err != nil {
		return err
	}

	if options.ClientID != "" {
		newToken := tokens.SelectAccount(identityClient, options.ClientID)
		if err := client.SetTokensAndAuthResult(newToken); err != nil {
			return err
		}
	}

	if options.AllowReauth {
		// here we're creating a throw-away client (tac). it's a copy of the user's provider client, but
		// with the token and reauth func zeroed out. combined with setting `AllowReauth` to `false`,
		// this should retry authentication only once
		tac := *client
		tac.SetThrowaway(true)
		tac.ReauthFunc = nil
		err = tac.SetTokensAndAuthResult(nil)
		if err != nil {
			return err
		}
		tro := client.ToTokenOptions()
		tao := options
		tao.AllowReauth = false
		client.ReauthFunc = func() error {
			err := refreshPlatform(&tac, endpoint, tro, tao, eo)
			if err != nil {
				errAuth := auth(&tac, endpoint, tao, eo)
				if errAuth != nil {
					return errAuth
				}
			}
			client.CopyTokensFrom(&tac)
			return nil
		}
	}

	return nil
}

func refreshPlatform(client *edgecloud.ProviderClient, endpoint string, tokenOptions edgecloud.TokenOptions, authOptions edgecloud.AuthOptions, eo edgecloud.EndpointOpts) error {

	identityClient, err := NewIdentity(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		identityClient.Endpoint = edgecloud.NormalizeURL(endpoint)
	}

	result := tokens.RefreshPlatform(identityClient, tokenOptions)

	err = client.SetTokensAndAuthResult(result)
	if err != nil {
		return err
	}

	if tokenOptions.AllowReauth {
		// here we're creating a throw-away client (tac). it's a copy of the user's provider client, but
		// with the token and reauth func zeroed out. combined with setting `AllowReauth` to `false`,
		// this should retry authentication only once
		tac := *client
		tac.SetThrowaway(true)
		tac.ReauthFunc = nil
		_ = tac.SetTokensAndAuthResult(nil)
		tro := tokenOptions
		tro.AllowReauth = false
		tao := authOptions
		tao.AllowReauth = false
		client.ReauthFunc = func() error {
			err := refreshPlatform(&tac, endpoint, tro, tao, eo)
			if err != nil {
				errAuth := auth(&tac, endpoint, tao, eo)
				if errAuth != nil {
					return errAuth
				}
			}
			client.CopyTokensFrom(&tac)
			return nil
		}
	}

	return nil
}

func refreshECCloud(client *edgecloud.ProviderClient, endpoint string, options edgecloud.TokenOptions, eo edgecloud.EndpointOpts) error {

	identityClient, err := NewIdentity(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		base, err := utils.BaseRootEndpoint(endpoint)
		if err != nil {
			return err
		}
		identityClient.Endpoint = edgecloud.NormalizeURL(base)
	}

	result := tokens.RefreshECCloud(identityClient, options)

	err = client.SetTokensAndAuthResult(result)
	if err != nil {
		return err
	}

	if options.AllowReauth {
		// here we're creating a throw-away client (tac). it's a copy of the user's provider client, but
		// with the token and reauth func zeroed out. combined with setting `AllowReauth` to `false`,
		// this should retry authentication only once
		tac := *client
		tac.SetThrowaway(true)
		tac.ReauthFunc = nil
		_ = tac.SetTokensAndAuthResult(nil)
		tao := options
		tao.AllowReauth = false
		client.ReauthFunc = func() error {
			err := refreshECCloud(&tac, endpoint, tao, eo)
			if err != nil {
				return err
			}
			client.CopyTokensFrom(&tac)
			return nil
		}
	}

	return nil
}

func setECCloudReauth(client *edgecloud.ProviderClient, endpoint string, options edgecloud.TokenOptions, eo edgecloud.EndpointOpts) {

	if options.AllowReauth {
		// here we're creating a throw-away client (tac). it's a copy of the user's provider client, but
		// with the token and reauth func zeroed out. combined with setting `AllowReauth` to `false`,
		// this should retry authentication only once
		tac := *client
		tac.SetThrowaway(true)
		tac.ReauthFunc = nil
		_ = tac.SetTokensAndAuthResult(nil)
		tao := options
		tao.AllowReauth = false
		client.ReauthFunc = func() error {
			err := refreshECCloud(&tac, endpoint, tao, eo)
			if err != nil {
				return err
			}
			client.CopyTokensFrom(&tac)
			return nil
		}
	}
}

// NewIdentity creates a ServiceClient that may be used to interact with the edgecenter identity auth service.
func NewIdentity(client *edgecloud.ProviderClient, eo edgecloud.EndpointOpts) (*edgecloud.ServiceClient, error) {
	endpoint := client.IdentityBase
	clientType := "auth"
	var err error
	if !reflect.DeepEqual(eo, edgecloud.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	return &edgecloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
		RegionID:       eo.Region,
	}, nil
}

func initClientOpts(client *edgecloud.ProviderClient, eo edgecloud.EndpointOpts, clientType string) (*edgecloud.ServiceClient, error) {
	sc := new(edgecloud.ServiceClient)
	eo.ApplyDefaults(clientType)
	url, err := edgecloud.DefaultEndpointLocator(client.APIBase)(eo)
	if err != nil {
		return sc, err
	}
	url, err = utils.NormalizeURLPath(url)
	if err != nil {
		return sc, err
	}
	sc.ProviderClient = client
	sc.Endpoint = fmt.Sprintf("%s%s/", client.APIBase, eo.Version)
	sc.ResourceBase = url
	sc.Type = clientType
	sc.RegionID = eo.Region
	return sc, nil
}

func APITokenClientServiceWithDebug(options edgecloud.APITokenOptions, eo edgecloud.EndpointOpts, debug bool) (*edgecloud.ServiceClient, error) {
	provider, err := APITokenClientWithDebug(options, debug)
	if err != nil {
		return nil, err
	}
	return ClientServiceFromProvider(provider, eo)
}

func TokenClientService(options edgecloud.TokenOptions, eo edgecloud.EndpointOpts) (*edgecloud.ServiceClient, error) {
	provider, err := TokenClient(options)
	if err != nil {
		return nil, err
	}
	return ClientServiceFromProvider(provider, eo)
}

func TokenClientServiceWithDebug(options edgecloud.TokenOptions, eo edgecloud.EndpointOpts, debug bool) (*edgecloud.ServiceClient, error) {
	provider, err := TokenClientWithDebug(options, debug)
	if err != nil {
		return nil, err
	}
	return ClientServiceFromProvider(provider, eo)
}

func AuthClientService(options edgecloud.AuthOptions, eo edgecloud.EndpointOpts) (*edgecloud.ServiceClient, error) {
	provider, err := AuthenticatedClient(options)
	if err != nil {
		return nil, err
	}
	return ClientServiceFromProvider(provider, eo)
}

func AuthClientServiceWithDebug(options edgecloud.AuthOptions, eo edgecloud.EndpointOpts, debug bool) (*edgecloud.ServiceClient, error) {
	provider, err := AuthenticatedClientWithDebug(options, debug)
	if err != nil {
		return nil, err
	}
	return ClientServiceFromProvider(provider, eo)
}

func ClientServiceFromProvider(provider *edgecloud.ProviderClient, eo edgecloud.EndpointOpts) (*edgecloud.ServiceClient, error) {
	client, err := initClientOpts(provider, eo, eo.Type)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newEndpointOpts(region int, project int, name string, version string) edgecloud.EndpointOpts {
	return edgecloud.EndpointOpts{
		Type:    "",
		Name:    name,
		Region:  region,
		Project: project,
		Version: version,
	}
}

func NewK8sV1(provider *edgecloud.ProviderClient, region int, project int) (*edgecloud.ServiceClient, error) {
	return ClientServiceFromProvider(provider, newEndpointOpts(region, project, "k8s", "v1"))
}
