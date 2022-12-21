package tokens

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func processToken(c *edgecloud.ServiceClient, opts edgecloud.AuthOptionsBuilder, url string) (r TokenResult) {
	b := opts.ToMap()
	resp, err := c.Post(url, b, &r.Body, &edgecloud.RequestOpts{})
	r.Err = err
	if resp != nil {
		r.Header = resp.Header
	}
	return
}

// Create authenticates and either generates a new token
func Create(c *edgecloud.ServiceClient, opts edgecloud.AuthOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, tokenURL(c))
}

// RefreshPlatform token with EdgeCenter platform API
func RefreshPlatform(c *edgecloud.ServiceClient, opts edgecloud.TokenOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, refreshURL(c))
}

// RefreshPlatform token with gcloud API
func RefreshECCloud(c *edgecloud.ServiceClient, opts edgecloud.TokenOptionsBuilder) (r TokenResult) {
	return processToken(c, opts, refreshECCloudURL(c))
}

// SelectAccount select an account which you want to get access to
func SelectAccount(c *edgecloud.ServiceClient, clientID string) (r TokenResult) {
	url := selectAccountURL(c, clientID)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
