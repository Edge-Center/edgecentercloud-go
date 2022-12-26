package securitygrouprules

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/securitygroup/v1/securitygroups"
)

// Replace accepts a CreateOpts struct and creates a new security group rule using the values provided.
func Replace(c *edgecloud.ServiceClient, ruleID string, opts securitygroups.CreateRuleOptsBuilder) (r securitygroups.CreateRuleResult) {
	b, err := opts.ToRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, ruleID), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the security group associated with it.
func Delete(c *edgecloud.ServiceClient, securityGroupID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, securityGroupID), nil)
	return
}
