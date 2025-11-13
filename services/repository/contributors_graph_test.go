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

package repository

import (
	"slices"
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/cache"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestRepository_ContributorsGraph(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	assert.NoError(t, repo.LoadOwner(t.Context()))
	mockCache, err := cache.NewStringCache(setting.Cache{})
	assert.NoError(t, err)

	generateContributorStats(nil, mockCache, "key", repo, "404ref")
	var data map[string]*ContributorData
	_, getErr := mockCache.GetJSON("key", &data)
	assert.NotNil(t, getErr)
	assert.ErrorContains(t, getErr.ToError(), "object does not exist")

	generateContributorStats(nil, mockCache, "key2", repo, "master")
	exist, _ := mockCache.GetJSON("key2", &data)
	assert.True(t, exist)
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	assert.Equal(t, []string{
		"ethantkoenig@gmail.com",
		"jimmy.praet@telenet.be",
		"jon@allspice.io",
		"total", // generated summary
	}, keys)

	assert.Equal(t, &ContributorData{
		Name:         "Ethan Koenig",
		AvatarLink:   "/assets/img/avatar_default.png",
		TotalCommits: 1,
		Weeks: map[int64]*WeekData{
			1511654400000: {
				Week:      1511654400000, // sunday 2017-11-26
				Additions: 3,
				Deletions: 0,
				Commits:   1,
			},
		},
	}, data["ethantkoenig@gmail.com"])
	assert.Equal(t, &ContributorData{
		Name:         "Total",
		AvatarLink:   "",
		TotalCommits: 3,
		Weeks: map[int64]*WeekData{
			1511654400000: {
				Week:      1511654400000, // sunday 2017-11-26 (2017-11-26 20:31:18 -0800)
				Additions: 3,
				Deletions: 0,
				Commits:   1,
			},
			1607817600000: {
				Week:      1607817600000, // sunday 2020-12-13 (2020-12-15 15:23:11 -0500)
				Additions: 10,
				Deletions: 0,
				Commits:   1,
			},
			1624752000000: {
				Week:      1624752000000, // sunday 2021-06-27 (2021-06-29 21:54:09 +0200)
				Additions: 2,
				Deletions: 0,
				Commits:   1,
			},
		},
	}, data["total"])
}
