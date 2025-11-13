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

package v1_20

import (
	"sort"
	"strings"
	"testing"

	"github.com/kumose/kmup/models/migrations/base"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Old OldAccessTokenScope
	New AccessTokenScope
}

func createOldTokenScope(scopes ...OldAccessTokenScope) OldAccessTokenScope {
	s := make([]string, 0, len(scopes))
	for _, os := range scopes {
		s = append(s, string(os))
	}
	return OldAccessTokenScope(strings.Join(s, ","))
}

func createNewTokenScope(scopes ...AccessTokenScope) AccessTokenScope {
	s := make([]string, 0, len(scopes))
	for _, os := range scopes {
		s = append(s, string(os))
	}
	return AccessTokenScope(strings.Join(s, ","))
}

func Test_ConvertScopedAccessTokens(t *testing.T) {
	tests := []testCase{
		{
			createOldTokenScope(OldAccessTokenScopeRepo, OldAccessTokenScopeUserFollow),
			createNewTokenScope(AccessTokenScopeWriteRepository, AccessTokenScopeWriteUser),
		},
		{
			createOldTokenScope(OldAccessTokenScopeUser, OldAccessTokenScopeWritePackage, OldAccessTokenScopeSudo),
			createNewTokenScope(AccessTokenScopeWriteAdmin, AccessTokenScopeWritePackage, AccessTokenScopeWriteUser),
		},
		{
			createOldTokenScope(),
			createNewTokenScope(),
		},
		{
			createOldTokenScope(OldAccessTokenScopeReadGPGKey, OldAccessTokenScopeReadOrg, OldAccessTokenScopeAll),
			createNewTokenScope(AccessTokenScopeAll),
		},
		{
			createOldTokenScope(OldAccessTokenScopeReadGPGKey, "invalid"),
			createNewTokenScope("invalid", AccessTokenScopeReadUser),
		},
	}

	// add a test for each individual mapping
	for oldScope, newScope := range accessTokenScopeMap {
		tests = append(tests, testCase{
			oldScope,
			createNewTokenScope(newScope...),
		})
	}

	x, deferable := base.PrepareTestEnv(t, 0, new(AccessToken))
	defer deferable()
	if x == nil || t.Failed() {
		t.Skip()
		return
	}

	// verify that no fixtures were loaded
	count, err := x.Count(&AccessToken{})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	for _, tc := range tests {
		_, err = x.Insert(&AccessToken{
			Scope: string(tc.Old),
		})
		assert.NoError(t, err)
	}

	// migrate the scopes
	err = ConvertScopedAccessTokens(x)
	assert.NoError(t, err)

	// migrate the scopes again (migration should be idempotent)
	err = ConvertScopedAccessTokens(x)
	assert.NoError(t, err)

	tokens := make([]AccessToken, 0)
	err = x.Find(&tokens)
	assert.NoError(t, err)
	assert.Len(t, tokens, len(tests))

	// sort the tokens (insertion order by auto-incrementing primary key)
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].ID < tokens[j].ID
	})

	// verify that the converted scopes are equal to the expected test result
	for idx, newToken := range tokens {
		assert.Equal(t, string(tests[idx].New), newToken.Scope)
	}
}
