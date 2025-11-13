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

package git_test

import (
	"testing"

	git_model "github.com/kumose/kmup/models/git"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestIsUserAllowed(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	pt := &git_model.ProtectedTag{}
	allowed, err := git_model.IsUserAllowedModifyTag(t.Context(), pt, 1)
	assert.NoError(t, err)
	assert.False(t, allowed)

	pt = &git_model.ProtectedTag{
		AllowlistUserIDs: []int64{1},
	}
	allowed, err = git_model.IsUserAllowedModifyTag(t.Context(), pt, 1)
	assert.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = git_model.IsUserAllowedModifyTag(t.Context(), pt, 2)
	assert.NoError(t, err)
	assert.False(t, allowed)

	pt = &git_model.ProtectedTag{
		AllowlistTeamIDs: []int64{1},
	}
	allowed, err = git_model.IsUserAllowedModifyTag(t.Context(), pt, 1)
	assert.NoError(t, err)
	assert.False(t, allowed)

	allowed, err = git_model.IsUserAllowedModifyTag(t.Context(), pt, 2)
	assert.NoError(t, err)
	assert.True(t, allowed)

	pt = &git_model.ProtectedTag{
		AllowlistUserIDs: []int64{1},
		AllowlistTeamIDs: []int64{1},
	}
	allowed, err = git_model.IsUserAllowedModifyTag(t.Context(), pt, 1)
	assert.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = git_model.IsUserAllowedModifyTag(t.Context(), pt, 2)
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestIsUserAllowedToControlTag(t *testing.T) {
	cases := []struct {
		name    string
		userid  int64
		allowed bool
	}{
		{
			name:    "test",
			userid:  1,
			allowed: true,
		},
		{
			name:    "test",
			userid:  3,
			allowed: true,
		},
		{
			name:    "kmup",
			userid:  1,
			allowed: true,
		},
		{
			name:    "kmup",
			userid:  3,
			allowed: false,
		},
		{
			name:    "test-kmup",
			userid:  1,
			allowed: true,
		},
		{
			name:    "test-kmup",
			userid:  3,
			allowed: false,
		},
		{
			name:    "kmup-test",
			userid:  1,
			allowed: true,
		},
		{
			name:    "kmup-test",
			userid:  3,
			allowed: true,
		},
		{
			name:    "v-1",
			userid:  1,
			allowed: false,
		},
		{
			name:    "v-1",
			userid:  2,
			allowed: true,
		},
		{
			name:    "release",
			userid:  1,
			allowed: false,
		},
	}

	t.Run("Glob", func(t *testing.T) {
		protectedTags := []*git_model.ProtectedTag{
			{
				NamePattern:      `*kmup`,
				AllowlistUserIDs: []int64{1},
			},
			{
				NamePattern:      `v-*`,
				AllowlistUserIDs: []int64{2},
			},
			{
				NamePattern: "release",
			},
		}

		for n, c := range cases {
			isAllowed, err := git_model.IsUserAllowedToControlTag(t.Context(), protectedTags, c.name, c.userid)
			assert.NoError(t, err)
			assert.Equal(t, c.allowed, isAllowed, "case %d: error should match", n)
		}
	})

	t.Run("Regex", func(t *testing.T) {
		protectedTags := []*git_model.ProtectedTag{
			{
				NamePattern:      `/kmup\z/`,
				AllowlistUserIDs: []int64{1},
			},
			{
				NamePattern:      `/\Av-/`,
				AllowlistUserIDs: []int64{2},
			},
			{
				NamePattern: "/release/",
			},
		}

		for n, c := range cases {
			isAllowed, err := git_model.IsUserAllowedToControlTag(t.Context(), protectedTags, c.name, c.userid)
			assert.NoError(t, err)
			assert.Equal(t, c.allowed, isAllowed, "case %d: error should match", n)
		}
	})
}
