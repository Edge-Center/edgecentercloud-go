package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	th "github.com/Edge-Center/edgecentercloud-go/testhelper"
)

func TestServiceURL(t *testing.T) {
	c := &edgecloud.ServiceClient{Endpoint: "http://123.45.67.8/"}
	expected := "http://123.45.67.8/more/parts/here"
	actual := c.ServiceURL("more", "parts", "here")
	th.CheckEquals(t, expected, actual)
}

func TestMoreHeaders(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	c := new(edgecloud.ServiceClient)
	c.MoreHeaders = map[string]string{
		"custom": "header",
	}
	c.ProviderClient = new(edgecloud.ProviderClient)
	resp, err := c.Get(fmt.Sprintf("%s/route", th.Endpoint()), nil, nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, resp.Request.Header.Get("custom"), "header")
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()
}
