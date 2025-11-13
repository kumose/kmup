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

package v1_22

import (
	"strconv"
	"testing"

	"github.com/kumose/kmup/models/migrations/base"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateBadgeColName(t *testing.T) {
	type Badge struct {
		ID          int64 `xorm:"pk autoincr"`
		Description string
		ImageURL    string
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(Badge))
	defer deferable()
	if x == nil || t.Failed() {
		return
	}

	oldBadges := []*Badge{
		{Description: "Test Badge 1", ImageURL: "https://example.com/badge1.png"},
		{Description: "Test Badge 2", ImageURL: "https://example.com/badge2.png"},
		{Description: "Test Badge 3", ImageURL: "https://example.com/badge3.png"},
	}

	for _, badge := range oldBadges {
		_, err := x.Insert(badge)
		assert.NoError(t, err)
	}

	if err := UseSlugInsteadOfIDForBadges(x); err != nil {
		assert.NoError(t, err)
		return
	}

	got := []BadgeUnique{}
	if err := x.Table("badge").Asc("id").Find(&got); !assert.NoError(t, err) {
		return
	}

	for i, e := range oldBadges {
		got := got[i+1] // 1 is in the badge.yml
		assert.Equal(t, e.ID, got.ID)
		assert.Equal(t, strconv.FormatInt(e.ID, 10), got.Slug)
	}

	// TODO: check if badges have been updated
}
