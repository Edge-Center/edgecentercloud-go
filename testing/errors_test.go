package testing

import (
	"testing"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	th "github.com/Edge-Center/edgecentercloud-go/testhelper"
)

func TestGetResponseCode(t *testing.T) {
	respErr := edgecloud.ErrUnexpectedResponseCode{
		URL:      "http://example.com",
		Method:   "GET",
		Expected: []int{200},
		Actual:   404,
		Body:     nil,
	}

	var err404 error = edgecloud.ErrDefault404{ErrUnexpectedResponseCode: respErr}

	err, ok := err404.(edgecloud.StatusCodeError)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, err.GetStatusCode(), 404)
}
