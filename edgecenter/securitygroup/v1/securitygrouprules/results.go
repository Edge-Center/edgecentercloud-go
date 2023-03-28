package securitygrouprules

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	edgecloud.ErrResult
}
