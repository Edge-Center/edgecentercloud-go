package ports

import (
	"net/http"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/reservedfixedip/v1/reservedfixedips"
)

// AllowAddressPairsOptsBuilder allows extensions to add additional parameters to the AllowAddressPairs request.
type AllowAddressPairsOptsBuilder interface {
	ToAllowAddressPairsMap() (map[string]interface{}, error)
}

// AllowAddressPairsOpts represents options used to allow address pairs.
type AllowAddressPairsOpts struct {
	AllowedAddressPairs []reservedfixedips.AllowedAddressPairs `json:"allowed_address_pairs"`
}

// ToAllowAddressPairsMap builds a request body from AllowAddressPairsOpts.
func (opts AllowAddressPairsOpts) ToAllowAddressPairsMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

// Validate AllowAddressPairsOpts.
func (opts AllowAddressPairsOpts) Validate() error {
	return edgecloud.TranslateValidationError(edgecloud.Validate.Struct(opts))
}

// EnablePortSecurity by portID.
func EnablePortSecurity(c *edgecloud.ServiceClient, portID string) (r UpdateResult) {
	_, r.Err = c.Post(enablePortSecurityURL(c, portID), nil, &r.Body, nil)
	return
}

// DisablePortSecurity by portID.
func DisablePortSecurity(c *edgecloud.ServiceClient, portID string) (r UpdateResult) {
	_, r.Err = c.Post(disablePortSecurityURL(c, portID), nil, &r.Body, nil)
	return
}

// AllowAddressPairs assign allowed address pairs for instance port.
func AllowAddressPairs(c *edgecloud.ServiceClient, portID string, opts AllowAddressPairsOptsBuilder) (r AssignResult) {
	b, err := opts.ToAllowAddressPairsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(assignAllowedAddressPairsURL(c, portID), b, &r.Body, &edgecloud.RequestOpts{OkCodes: []int{http.StatusOK}})
	return
}
