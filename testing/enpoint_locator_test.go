package testing

import (
	"testing"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/stretchr/testify/require"
)

func TestEndpointLocationWithoutRegionAndProject(t *testing.T) {
	baseEndpoint := "http://test.com"

	eo := edgecloud.EndpointOpts{
		Type:    "test",
		Name:    "test",
		Region:  0,
		Project: 0,
		Version: "v1",
	}

	el := edgecloud.DefaultEndpointLocator(baseEndpoint)

	url, err := el(eo)
	require.NoError(t, err)
	require.Equal(t, "http://test.com/v1/test///test", url)
}
