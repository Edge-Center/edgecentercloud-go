package testing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractToken(t *testing.T) {
	result := getTokenResult(t)
	token, err := result.ExtractTokens()
	require.NoError(t, err)
	require.Equal(t, &expectedToken, token)
}
