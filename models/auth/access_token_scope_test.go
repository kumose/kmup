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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type scopeTestNormalize struct {
	in  AccessTokenScope
	out AccessTokenScope
	err error
}

func TestAccessTokenScope_Normalize(t *testing.T) {
	assert.Equal(t, []string{"activitypub", "admin", "issue", "misc", "notification", "organization", "package", "repository", "user"}, GetAccessTokenCategories())
	tests := []scopeTestNormalize{
		{"", "", nil},
		{"write:misc,write:notification,read:package,write:notification,public-only", "public-only,write:misc,write:notification,read:package", nil},
		{"all", "all", nil},
		{"write:activitypub,write:admin,write:misc,write:notification,write:organization,write:package,write:issue,write:repository,write:user", "all", nil},
		{"write:activitypub,write:admin,write:misc,write:notification,write:organization,write:package,write:issue,write:repository,write:user,public-only", "public-only,all", nil},
	}

	for _, scope := range GetAccessTokenCategories() {
		tests = append(tests,
			scopeTestNormalize{AccessTokenScope("read:" + scope), AccessTokenScope("read:" + scope), nil},
			scopeTestNormalize{AccessTokenScope("write:" + scope), AccessTokenScope("write:" + scope), nil},
			scopeTestNormalize{AccessTokenScope(fmt.Sprintf("write:%[1]s,read:%[1]s", scope)), AccessTokenScope("write:" + scope), nil},
			scopeTestNormalize{AccessTokenScope(fmt.Sprintf("read:%[1]s,write:%[1]s", scope)), AccessTokenScope("write:" + scope), nil},
			scopeTestNormalize{AccessTokenScope(fmt.Sprintf("read:%[1]s,write:%[1]s,write:%[1]s", scope)), AccessTokenScope("write:" + scope), nil},
		)
	}

	for _, test := range tests {
		t.Run(string(test.in), func(t *testing.T) {
			scope, err := test.in.Normalize()
			assert.Equal(t, test.out, scope)
			assert.Equal(t, test.err, err)
		})
	}
}

type scopeTestHasScope struct {
	in    AccessTokenScope
	scope AccessTokenScope
	out   bool
	err   error
}

func TestAccessTokenScope_HasScope(t *testing.T) {
	tests := []scopeTestHasScope{
		{"read:admin", "write:package", false, nil},
		{"all", "write:package", true, nil},
		{"write:package", "all", false, nil},
		{"public-only", "read:issue", false, nil},
	}

	for _, scope := range GetAccessTokenCategories() {
		tests = append(tests,
			scopeTestHasScope{
				AccessTokenScope("read:" + scope),
				AccessTokenScope("read:" + scope), true, nil,
			},
			scopeTestHasScope{
				AccessTokenScope("write:" + scope),
				AccessTokenScope("write:" + scope), true, nil,
			},
			scopeTestHasScope{
				AccessTokenScope("write:" + scope),
				AccessTokenScope("read:" + scope), true, nil,
			},
			scopeTestHasScope{
				AccessTokenScope("read:" + scope),
				AccessTokenScope("write:" + scope), false, nil,
			},
		)
	}

	for _, test := range tests {
		t.Run(string(test.in), func(t *testing.T) {
			hasScope, err := test.in.HasScope(test.scope)
			assert.Equal(t, test.out, hasScope)
			assert.Equal(t, test.err, err)
		})
	}
}
