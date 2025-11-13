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

package system_test

import (
	"testing"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/models/system"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestNotice_TrStr(t *testing.T) {
	notice := &system.Notice{
		Type:        system.NoticeRepository,
		Description: "test description",
	}
	assert.Equal(t, "admin.notices.type_1", notice.TrStr())
}

func TestCreateNotice(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	noticeBean := &system.Notice{
		Type:        system.NoticeRepository,
		Description: "test description",
	}
	unittest.AssertNotExistsBean(t, noticeBean)
	assert.NoError(t, system.CreateNotice(t.Context(), noticeBean.Type, noticeBean.Description))
	unittest.AssertExistsAndLoadBean(t, noticeBean)
}

func TestCreateRepositoryNotice(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	noticeBean := &system.Notice{
		Type:        system.NoticeRepository,
		Description: "test description",
	}
	unittest.AssertNotExistsBean(t, noticeBean)
	assert.NoError(t, system.CreateRepositoryNotice(noticeBean.Description))
	unittest.AssertExistsAndLoadBean(t, noticeBean)
}

func TestCountNotices(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	assert.Equal(t, int64(3), system.CountNotices(t.Context()))
}

func TestNotices(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	notices, err := system.Notices(t.Context(), 1, 2)
	assert.NoError(t, err)
	if assert.Len(t, notices, 2) {
		assert.Equal(t, int64(3), notices[0].ID)
		assert.Equal(t, int64(2), notices[1].ID)
	}

	notices, err = system.Notices(t.Context(), 2, 2)
	assert.NoError(t, err)
	if assert.Len(t, notices, 1) {
		assert.Equal(t, int64(1), notices[0].ID)
	}
}

func TestDeleteNotices(t *testing.T) {
	// delete a non-empty range
	assert.NoError(t, unittest.PrepareTestDatabase())

	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 1})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 2})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 3})
	assert.NoError(t, system.DeleteNotices(t.Context(), 1, 2))
	unittest.AssertNotExistsBean(t, &system.Notice{ID: 1})
	unittest.AssertNotExistsBean(t, &system.Notice{ID: 2})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 3})
}

func TestDeleteNotices2(t *testing.T) {
	// delete an empty range
	assert.NoError(t, unittest.PrepareTestDatabase())

	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 1})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 2})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 3})
	assert.NoError(t, system.DeleteNotices(t.Context(), 3, 2))
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 1})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 2})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 3})
}

func TestDeleteNoticesByIDs(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 1})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 2})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 3})
	err := db.DeleteByIDs[system.Notice](t.Context(), 1, 3)
	assert.NoError(t, err)
	unittest.AssertNotExistsBean(t, &system.Notice{ID: 1})
	unittest.AssertExistsAndLoadBean(t, &system.Notice{ID: 2})
	unittest.AssertNotExistsBean(t, &system.Notice{ID: 3})
}
