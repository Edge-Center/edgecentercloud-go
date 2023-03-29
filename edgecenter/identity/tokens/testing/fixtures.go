package testing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/identity/tokens"
	"github.com/Edge-Center/edgecentercloud-go/testhelper"
	fake "github.com/Edge-Center/edgecentercloud-go/testhelper/client"
)

// TokenOutput is a sample response to a AccessToken call.
var TokenOutput = fmt.Sprintf(`
{
   "access": "%s",
   "refresh": "%s"
}`, fake.AccessToken,
	fake.RefreshToken,
)

var expectedToken = tokens.Token{
	Access:  fake.AccessToken,
	Refresh: fake.RefreshToken,
}

func getTokenResult(t *testing.T) tokens.TokenResult {
	t.Helper()
	result := tokens.TokenResult{}
	result.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", fake.AccessToken)},
	}
	err := json.Unmarshal([]byte(TokenOutput), &result.Body)
	testhelper.AssertNoErr(t, err)
	return result
}
