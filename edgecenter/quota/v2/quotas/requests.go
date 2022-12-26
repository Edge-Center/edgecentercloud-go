package quotas

import edgecloud "github.com/Edge-Center/edgecentercloud-go"

// ListCombinedOptsBuilder allows extensions to add additional parameters to the List request.
type ListCombinedOptsBuilder interface {
	ToCombinedListQuery() (string, error)
}

// ListCombinedOpts allows the filtering and sorting List API response.
type ListCombinedOpts struct {
	ClientID int `q:"client_id"`
}

// ToCombinedListQuery formats a ListCombinedOpts into a query string.
func (opts ListCombinedOpts) ToCombinedListQuery() (string, error) {
	q, err := edgecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListCombined retrieves list of combined quota.
func ListCombined(c *edgecloud.ServiceClient, opts ListCombinedOptsBuilder) (r CombinedResult) {
	url := getCombinedURL(c)
	if opts != nil {
		query, err := opts.ToCombinedListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// ListGlobal retrieves list of global quota.
func ListGlobal(c *edgecloud.ServiceClient, clientID int) (r CommonResult) {
	url := getGlobalURL(c, clientID)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// ListRegional retrieves list of regional quota.
func ListRegional(c *edgecloud.ServiceClient, clientID, regionID int) (r CommonResult) {
	url := getRegionURL(c, clientID, regionID)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
