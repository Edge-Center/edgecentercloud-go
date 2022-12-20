package edgecenter

import (
	"os"
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var nilAuthOptions = edgecloud.AuthOptions{}
var nilTokenOptions = edgecloud.TokenOptions{}

/*
AuthOptionsFromEnv fills out an identity.AuthOptions structure with the
settings found on environment variables.

The following variables provide sources of truth: EC_CLOUD_USERNAME, EC_CLOUD_PASSWORD, EC_CLOUD_AUTH_URL

	opts, err := edgecenter.AuthOptionsFromEnv()
	provider, err := edgecenter.AuthenticatedClient(opts)
*/
func AuthOptionsFromEnv() (edgecloud.AuthOptions, error) {
	authURL := os.Getenv("EC_CLOUD_AUTH_URL")
	username := os.Getenv("EC_CLOUD_USERNAME")
	password := os.Getenv("EC_CLOUD_PASSWORD")

	if authURL == "" {
		err := edgecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "EC_CLOUD_AUTH_URL",
		}
		return nilAuthOptions, err
	}

	if username == "" {
		err := edgecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "EC_CLOUD_USERNAME",
		}
		return nilAuthOptions, err
	}

	if password == "" {
		err := edgecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "EC_CLOUD_PASSWORD",
		}
		return nilAuthOptions, err
	}

	ao := edgecloud.AuthOptions{
		APIURL:   authURL,
		Username: username,
		Password: password,
	}

	return ao, nil
}

func TokenOptionsFromEnv() (edgecloud.TokenOptions, error) {

	apiURL := os.Getenv("EC_CLOUD_API_URL")
	accessToken := os.Getenv("EC_CLOUD_ACCESS_TOKEN")
	refreshToken := os.Getenv("EC_CLOUD_REFRESH_TOKEN")

	if apiURL == "" {
		err := edgecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "EC_CLOUD_API_URL",
		}
		return nilTokenOptions, err
	}

	if accessToken == "" {
		err := edgecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "EC_CLOUD_ACCESS_TOKEN",
		}
		return nilTokenOptions, err
	}

	if refreshToken == "" {
		err := edgecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "EC_CLOUD_REFRESH_TOKEN",
		}
		return nilTokenOptions, err
	}

	to := edgecloud.TokenOptions{
		APIURL:       apiURL,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AllowReauth:  true,
	}

	return to, nil
}

func NewECCloudPlatformAPISettingsFromEnv() (*edgecloud.PasswordAPISettings, error) {
	authURL := os.Getenv("EC_CLOUD_AUTH_URL")
	apiURL := os.Getenv("EC_CLOUD_API_URL")
	username := os.Getenv("EC_CLOUD_USERNAME")
	password := os.Getenv("EC_CLOUD_PASSWORD")
	apiVersion := os.Getenv("EC_CLOUD_API_VERSION")
	region := os.Getenv("EC_CLOUD_REGION")
	project := os.Getenv("EC_CLOUD_PROJECT")
	debugEnv := os.Getenv("EC_CLOUD_DEBUG")

	var (
		projectInt, regionInt int
		err                   error
		version               = "v1"
		debug                 bool
	)

	if project != "" {
		projectInt, err = strconv.Atoi(project)
		if err != nil {
			return nil, err
		}
	}

	if region != "" {
		regionInt, err = strconv.Atoi(region)
		if err != nil {
			return nil, err
		}
	}

	if apiVersion != "" {
		version = apiVersion
	}

	debug, err = strconv.ParseBool(debugEnv)
	if err != nil {
		debug = false
	}

	return &edgecloud.PasswordAPISettings{
		Version:     version,
		APIURL:      apiURL,
		AuthURL:     authURL,
		Username:    username,
		Password:    password,
		Region:      regionInt,
		Project:     projectInt,
		AllowReauth: true,
		Debug:       debug,
	}, nil
}

func NewEdgeCloudTokenAPISettingsFromEnv() (*edgecloud.TokenAPISettings, error) {
	apiURL := os.Getenv("EC_CLOUD_API_URL")
	apiVersion := os.Getenv("EC_CLOUD_API_VERSION")
	accessToken := os.Getenv("EC_CLOUD_ACCESS_TOKEN")
	refreshToken := os.Getenv("EC_CLOUD_REFRESH_TOKEN")
	region := os.Getenv("EC_CLOUD_REGION")
	project := os.Getenv("EC_CLOUD_PROJECT")
	debugEnv := os.Getenv("EC_CLOUD_DEBUG")

	var (
		projectInt, regionInt int
		err                   error
		version               = "v1"
		debug                 bool
	)

	if project != "" {
		projectInt, err = strconv.Atoi(project)
		if err != nil {
			return nil, err
		}
	}

	if region != "" {
		regionInt, err = strconv.Atoi(region)
		if err != nil {
			return nil, err
		}
	}

	if apiVersion != "" {
		version = apiVersion
	}

	debug, err = strconv.ParseBool(debugEnv)
	if err != nil {
		debug = false
	}

	return &edgecloud.TokenAPISettings{
		Version:      version,
		APIURL:       apiURL,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Region:       regionInt,
		Project:      projectInt,
		AllowReauth:  true,
		Debug:        debug,
	}, nil
}

func NewECCloudAPITokenAPISettingsFromEnv() (*edgecloud.APITokenAPISettings, error) {
	apiURL := os.Getenv("EC_CLOUD_API_URL")
	apiVersion := os.Getenv("EC_CLOUD_API_VERSION")
	region := os.Getenv("EC_CLOUD_REGION")
	project := os.Getenv("EC_CLOUD_PROJECT")
	debugEnv := os.Getenv("EC_CLOUD_DEBUG")
	apiToken := os.Getenv("EC_CLOUD_API_TOKEN")

	var (
		projectInt, regionInt int
		err                   error
		version               = "v1"
		debug                 bool
	)

	if project != "" {
		projectInt, err = strconv.Atoi(project)
		if err != nil {
			return nil, err
		}
	}

	if region != "" {
		regionInt, err = strconv.Atoi(region)
		if err != nil {
			return nil, err
		}
	}

	if apiVersion != "" {
		version = apiVersion
	}

	debug, err = strconv.ParseBool(debugEnv)
	if err != nil {
		debug = false
	}

	return &edgecloud.APITokenAPISettings{
		Version:  version,
		APIURL:   apiURL,
		Region:   regionInt,
		Project:  projectInt,
		APIToken: apiToken,
		Debug:    debug,
	}, nil
}
