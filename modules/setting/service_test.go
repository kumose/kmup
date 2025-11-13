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

package setting

import (
	"testing"

	"github.com/kumose/kmup/modules/glob"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestLoadServices(t *testing.T) {
	defer test.MockVariableValue(&Service)()

	cfg, err := NewConfigProviderFromData(`
[service]
EMAIL_DOMAIN_WHITELIST = d1, *.w
EMAIL_DOMAIN_ALLOWLIST = d2, *.a
EMAIL_DOMAIN_BLOCKLIST = d3, *.b
`)
	assert.NoError(t, err)
	loadServiceFrom(cfg)

	match := func(globs []glob.Glob, s string) bool {
		for _, g := range globs {
			if g.Match(s) {
				return true
			}
		}
		return false
	}

	assert.True(t, match(Service.EmailDomainAllowList, "d1"))
	assert.True(t, match(Service.EmailDomainAllowList, "foo.w"))
	assert.True(t, match(Service.EmailDomainAllowList, "d2"))
	assert.True(t, match(Service.EmailDomainAllowList, "foo.a"))
	assert.False(t, match(Service.EmailDomainAllowList, "d3"))

	assert.True(t, match(Service.EmailDomainBlockList, "d3"))
	assert.True(t, match(Service.EmailDomainBlockList, "foo.b"))
	assert.False(t, match(Service.EmailDomainBlockList, "d1"))
}

func TestLoadServiceVisibilityModes(t *testing.T) {
	defer test.MockVariableValue(&Service)()

	kases := map[string]func(){
		`
[service]
DEFAULT_USER_VISIBILITY = public
ALLOWED_USER_VISIBILITY_MODES = public,limited,private
`: func() {
			assert.Equal(t, "public", Service.DefaultUserVisibility)
			assert.Equal(t, structs.VisibleTypePublic, Service.DefaultUserVisibilityMode)
			assert.Equal(t, []string{"public", "limited", "private"}, Service.AllowedUserVisibilityModes)
		},
		`
		[service]
		DEFAULT_USER_VISIBILITY = public
		`: func() {
			assert.Equal(t, "public", Service.DefaultUserVisibility)
			assert.Equal(t, structs.VisibleTypePublic, Service.DefaultUserVisibilityMode)
			assert.Equal(t, []string{"public", "limited", "private"}, Service.AllowedUserVisibilityModes)
		},
		`
		[service]
		DEFAULT_USER_VISIBILITY = limited
		`: func() {
			assert.Equal(t, "limited", Service.DefaultUserVisibility)
			assert.Equal(t, structs.VisibleTypeLimited, Service.DefaultUserVisibilityMode)
			assert.Equal(t, []string{"public", "limited", "private"}, Service.AllowedUserVisibilityModes)
		},
		`
[service]
ALLOWED_USER_VISIBILITY_MODES = public,limited,private
`: func() {
			assert.Equal(t, "public", Service.DefaultUserVisibility)
			assert.Equal(t, structs.VisibleTypePublic, Service.DefaultUserVisibilityMode)
			assert.Equal(t, []string{"public", "limited", "private"}, Service.AllowedUserVisibilityModes)
		},
		`
[service]
DEFAULT_USER_VISIBILITY = public
ALLOWED_USER_VISIBILITY_MODES = limited,private
`: func() {
			assert.Equal(t, "limited", Service.DefaultUserVisibility)
			assert.Equal(t, structs.VisibleTypeLimited, Service.DefaultUserVisibilityMode)
			assert.Equal(t, []string{"limited", "private"}, Service.AllowedUserVisibilityModes)
		},
		`
[service]
DEFAULT_USER_VISIBILITY = my_type
ALLOWED_USER_VISIBILITY_MODES = limited,private
`: func() {
			assert.Equal(t, "limited", Service.DefaultUserVisibility)
			assert.Equal(t, structs.VisibleTypeLimited, Service.DefaultUserVisibilityMode)
			assert.Equal(t, []string{"limited", "private"}, Service.AllowedUserVisibilityModes)
		},
		`
[service]
DEFAULT_USER_VISIBILITY = public
ALLOWED_USER_VISIBILITY_MODES = public, limit, privated
`: func() {
			assert.Equal(t, "public", Service.DefaultUserVisibility)
			assert.Equal(t, structs.VisibleTypePublic, Service.DefaultUserVisibilityMode)
			assert.Equal(t, []string{"public"}, Service.AllowedUserVisibilityModes)
		},
	}

	for kase, fun := range kases {
		t.Run(kase, func(t *testing.T) {
			cfg, err := NewConfigProviderFromData(kase)
			assert.NoError(t, err)
			loadServiceFrom(cfg)
			fun()
			// reset
			Service.AllowedUserVisibilityModesSlice = []bool{true, true, true}
			Service.AllowedUserVisibilityModes = []string{}
			Service.DefaultUserVisibility = ""
			Service.DefaultUserVisibilityMode = structs.VisibleTypePublic
		})
	}
}

func TestLoadServiceRequireSignInView(t *testing.T) {
	defer test.MockVariableValue(&Service)()

	cfg, err := NewConfigProviderFromData(`
[service]
`)
	assert.NoError(t, err)
	loadServiceFrom(cfg)
	assert.False(t, Service.RequireSignInViewStrict)
	assert.False(t, Service.BlockAnonymousAccessExpensive)

	cfg, err = NewConfigProviderFromData(`
[service]
REQUIRE_SIGNIN_VIEW = true
`)
	assert.NoError(t, err)
	loadServiceFrom(cfg)
	assert.True(t, Service.RequireSignInViewStrict)
	assert.False(t, Service.BlockAnonymousAccessExpensive)

	cfg, err = NewConfigProviderFromData(`
[service]
REQUIRE_SIGNIN_VIEW = expensive
`)
	assert.NoError(t, err)
	loadServiceFrom(cfg)
	assert.False(t, Service.RequireSignInViewStrict)
	assert.True(t, Service.BlockAnonymousAccessExpensive)
}
