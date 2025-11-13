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

package unit

import (
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestLoadUnitConfig(t *testing.T) {
	t.Run("regular", func(t *testing.T) {
		defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []Type) {
			DisabledRepoUnitsSet(disabledRepoUnits)
			DefaultRepoUnits = defaultRepoUnits
			DefaultForkRepoUnits = defaultForkRepoUnits
		}(DisabledRepoUnitsGet(), DefaultRepoUnits, DefaultForkRepoUnits)
		defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []string) {
			setting.Repository.DisabledRepoUnits = disabledRepoUnits
			setting.Repository.DefaultRepoUnits = defaultRepoUnits
			setting.Repository.DefaultForkRepoUnits = defaultForkRepoUnits
		}(setting.Repository.DisabledRepoUnits, setting.Repository.DefaultRepoUnits, setting.Repository.DefaultForkRepoUnits)

		setting.Repository.DisabledRepoUnits = []string{"repo.issues"}
		setting.Repository.DefaultRepoUnits = []string{"repo.code", "repo.releases", "repo.issues", "repo.pulls"}
		setting.Repository.DefaultForkRepoUnits = []string{"repo.releases"}
		assert.NoError(t, LoadUnitConfig())
		assert.Equal(t, []Type{TypeIssues}, DisabledRepoUnitsGet())
		assert.Equal(t, []Type{TypeCode, TypeReleases, TypePullRequests}, DefaultRepoUnits)
		assert.Equal(t, []Type{TypeReleases}, DefaultForkRepoUnits)
	})
	t.Run("invalid", func(t *testing.T) {
		defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []Type) {
			DisabledRepoUnitsSet(disabledRepoUnits)
			DefaultRepoUnits = defaultRepoUnits
			DefaultForkRepoUnits = defaultForkRepoUnits
		}(DisabledRepoUnitsGet(), DefaultRepoUnits, DefaultForkRepoUnits)
		defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []string) {
			setting.Repository.DisabledRepoUnits = disabledRepoUnits
			setting.Repository.DefaultRepoUnits = defaultRepoUnits
			setting.Repository.DefaultForkRepoUnits = defaultForkRepoUnits
		}(setting.Repository.DisabledRepoUnits, setting.Repository.DefaultRepoUnits, setting.Repository.DefaultForkRepoUnits)

		setting.Repository.DisabledRepoUnits = []string{"repo.issues", "invalid.1"}
		setting.Repository.DefaultRepoUnits = []string{"repo.code", "invalid.2", "repo.releases", "repo.issues", "repo.pulls"}
		setting.Repository.DefaultForkRepoUnits = []string{"invalid.3", "repo.releases"}
		assert.NoError(t, LoadUnitConfig())
		assert.Equal(t, []Type{TypeIssues}, DisabledRepoUnitsGet())
		assert.Equal(t, []Type{TypeCode, TypeReleases, TypePullRequests}, DefaultRepoUnits)
		assert.Equal(t, []Type{TypeReleases}, DefaultForkRepoUnits)
	})
	t.Run("duplicate", func(t *testing.T) {
		defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []Type) {
			DisabledRepoUnitsSet(disabledRepoUnits)
			DefaultRepoUnits = defaultRepoUnits
			DefaultForkRepoUnits = defaultForkRepoUnits
		}(DisabledRepoUnitsGet(), DefaultRepoUnits, DefaultForkRepoUnits)
		defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []string) {
			setting.Repository.DisabledRepoUnits = disabledRepoUnits
			setting.Repository.DefaultRepoUnits = defaultRepoUnits
			setting.Repository.DefaultForkRepoUnits = defaultForkRepoUnits
		}(setting.Repository.DisabledRepoUnits, setting.Repository.DefaultRepoUnits, setting.Repository.DefaultForkRepoUnits)

		setting.Repository.DisabledRepoUnits = []string{"repo.issues", "repo.issues"}
		setting.Repository.DefaultRepoUnits = []string{"repo.code", "repo.releases", "repo.issues", "repo.pulls", "repo.code"}
		setting.Repository.DefaultForkRepoUnits = []string{"repo.releases", "repo.releases"}
		assert.NoError(t, LoadUnitConfig())
		assert.Equal(t, []Type{TypeIssues}, DisabledRepoUnitsGet())
		assert.Equal(t, []Type{TypeCode, TypeReleases, TypePullRequests}, DefaultRepoUnits)
		assert.Equal(t, []Type{TypeReleases}, DefaultForkRepoUnits)
	})
	t.Run("empty_default", func(t *testing.T) {
		defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []Type) {
			DisabledRepoUnitsSet(disabledRepoUnits)
			DefaultRepoUnits = defaultRepoUnits
			DefaultForkRepoUnits = defaultForkRepoUnits
		}(DisabledRepoUnitsGet(), DefaultRepoUnits, DefaultForkRepoUnits)
		defer func(disabledRepoUnits, defaultRepoUnits, defaultForkRepoUnits []string) {
			setting.Repository.DisabledRepoUnits = disabledRepoUnits
			setting.Repository.DefaultRepoUnits = defaultRepoUnits
			setting.Repository.DefaultForkRepoUnits = defaultForkRepoUnits
		}(setting.Repository.DisabledRepoUnits, setting.Repository.DefaultRepoUnits, setting.Repository.DefaultForkRepoUnits)

		setting.Repository.DisabledRepoUnits = []string{"repo.issues", "repo.issues"}
		setting.Repository.DefaultRepoUnits = []string{}
		setting.Repository.DefaultForkRepoUnits = []string{"repo.releases", "repo.releases"}
		assert.NoError(t, LoadUnitConfig())
		assert.Equal(t, []Type{TypeIssues}, DisabledRepoUnitsGet())
		assert.ElementsMatch(t, []Type{TypeCode, TypePullRequests, TypeReleases, TypeWiki, TypePackages, TypeProjects, TypeActions}, DefaultRepoUnits)
		assert.Equal(t, []Type{TypeReleases}, DefaultForkRepoUnits)
	})
}
