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

package paginator

import (
	"testing"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestPaginator(t *testing.T) {
	cases := []struct {
		db.Paginator
		Skip  int
		Take  int
		Start int
		End   int
	}{
		{
			Paginator: &db.ListOptions{Page: -1, PageSize: -1},
			Skip:      0,
			Take:      setting.API.DefaultPagingNum,
			Start:     0,
			End:       setting.API.DefaultPagingNum,
		},
		{
			Paginator: &db.ListOptions{Page: 2, PageSize: 10},
			Skip:      10,
			Take:      10,
			Start:     10,
			End:       20,
		},
		{
			Paginator: db.NewAbsoluteListOptions(-1, -1),
			Skip:      0,
			Take:      setting.API.DefaultPagingNum,
			Start:     0,
			End:       setting.API.DefaultPagingNum,
		},
		{
			Paginator: db.NewAbsoluteListOptions(2, 10),
			Skip:      2,
			Take:      10,
			Start:     2,
			End:       12,
		},
	}

	for _, c := range cases {
		skip, take := c.Paginator.GetSkipTake()

		assert.Equal(t, c.Skip, skip)
		assert.Equal(t, c.Take, take)
	}
}
