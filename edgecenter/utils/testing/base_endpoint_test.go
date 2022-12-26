package testing

import (
	"testing"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/utils"
	th "github.com/Edge-Center/edgecentercloud-go/testhelper"
)

type urlTestCases struct {
	URL     string
	NormURL string
}

func TestNormalizePath(t *testing.T) {
	tests := []urlTestCases{
		{
			URL:     "http://example.com:5000/v3////",
			NormURL: "http://example.com:5000/v3/",
		},
		{
			URL:     "http://example.com:5000/////v3",
			NormURL: "http://example.com:5000/v3/",
		},
	}

	for _, test := range tests {
		actual, err := utils.NormalizeURLPath(test.URL)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, test.NormURL, actual)
	}
}
