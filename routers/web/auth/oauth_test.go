// Copyright (C) Kumo inc. and its affiliates.
// Author: Jeff.li lijippy@163.com
// All rights reserved.
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package auth

import (
	"testing"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/services/oauth2_provider"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func createAndParseToken(t *testing.T, grant *auth.OAuth2Grant) *oauth2_provider.OIDCToken {
	signingKey, err := oauth2_provider.CreateJWTSigningKey("HS256", make([]byte, 32))
	assert.NoError(t, err)
	assert.NotNil(t, signingKey)

	response, terr := oauth2_provider.NewAccessTokenResponse(t.Context(), grant, signingKey, signingKey)
	assert.Nil(t, terr)
	assert.NotNil(t, response)

	parsedToken, err := jwt.ParseWithClaims(response.IDToken, &oauth2_provider.OIDCToken{}, func(token *jwt.Token) (any, error) {
		assert.NotNil(t, token.Method)
		assert.Equal(t, signingKey.SigningMethod().Alg(), token.Method.Alg())
		return signingKey.VerifyKey(), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	oidcToken, ok := parsedToken.Claims.(*oauth2_provider.OIDCToken)
	assert.True(t, ok)
	assert.NotNil(t, oidcToken)

	return oidcToken
}

func TestNewAccessTokenResponse_OIDCToken(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	grants, err := auth.GetOAuth2GrantsByUserID(t.Context(), 3)
	assert.NoError(t, err)
	assert.Len(t, grants, 1)

	// Scopes: openid
	oidcToken := createAndParseToken(t, grants[0])
	assert.Empty(t, oidcToken.Name)
	assert.Empty(t, oidcToken.PreferredUsername)
	assert.Empty(t, oidcToken.Profile)
	assert.Empty(t, oidcToken.Picture)
	assert.Empty(t, oidcToken.Website)
	assert.Empty(t, oidcToken.UpdatedAt)
	assert.Empty(t, oidcToken.Email)
	assert.False(t, oidcToken.EmailVerified)

	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 5})
	grants, err = auth.GetOAuth2GrantsByUserID(t.Context(), user.ID)
	assert.NoError(t, err)
	assert.Len(t, grants, 1)

	// Scopes: openid profile email
	oidcToken = createAndParseToken(t, grants[0])
	assert.Equal(t, user.DisplayName(), oidcToken.Name)
	assert.Equal(t, user.Name, oidcToken.PreferredUsername)
	assert.Equal(t, user.HTMLURL(t.Context()), oidcToken.Profile)
	assert.Equal(t, user.AvatarLink(t.Context()), oidcToken.Picture)
	assert.Equal(t, user.Website, oidcToken.Website)
	assert.Equal(t, user.UpdatedUnix, oidcToken.UpdatedAt)
	assert.Equal(t, user.Email, oidcToken.Email)
	assert.Equal(t, user.IsActive, oidcToken.EmailVerified)
}
