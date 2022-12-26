package edgecenter

import (
	"os"
	"strconv"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var nilEndpointOptions = edgecloud.EndpointOpts{}

func EndpointOptionsFromEnv() (edgecloud.EndpointOpts, error) {
	region := os.Getenv("EC_CLOUD_REGION")
	project := os.Getenv("EC_CLOUD_PROJECT")

	if region == "" {
		err := edgecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "EC_CLOUD_REGION",
		}
		return nilEndpointOptions, err
	}

	regionInt, err := strconv.Atoi(region)
	if err != nil {
		return nilEndpointOptions, err
	}
	if project == "" {
		err := edgecloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "EC_CLOUD_PROJECT",
		}
		return nilEndpointOptions, err
	}

	projectInt, err := strconv.Atoi(project)
	if err != nil {
		return nilEndpointOptions, err
	}

	eo := edgecloud.EndpointOpts{
		Region:  regionInt,
		Project: projectInt,
	}

	return eo, nil
}
