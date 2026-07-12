package domains

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOAuthTokenError_Error(t *testing.T) {
	require.Equal(t, "invalid_grant", (&OAuthTokenError{Code: "invalid_grant"}).Error())
	require.Equal(t, "bad token", (&OAuthTokenError{Code: "invalid_grant", Description: "bad token"}).Error())
}

func TestOIDCUserClaims_SplitDisplayName(t *testing.T) {
	var nilClaims *OIDCUserClaims
	first, last := nilClaims.SplitDisplayName()
	require.Empty(t, first)
	require.Empty(t, last)

	claims := &OIDCUserClaims{GivenName: "Jane", FamilyName: "Doe"}
	first, last = claims.SplitDisplayName()
	require.Equal(t, "Jane", first)
	require.Equal(t, "Doe", last)

	claims = &OIDCUserClaims{Name: "John Smith"}
	first, last = claims.SplitDisplayName()
	require.Equal(t, "John", first)
	require.Equal(t, "Smith", last)

	claims = &OIDCUserClaims{Name: "Madonna"}
	first, last = claims.SplitDisplayName()
	require.Equal(t, "Madonna", first)
	require.Empty(t, last)
}
