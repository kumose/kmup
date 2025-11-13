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

package repo

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestMigrate_InsertReleases(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	a := &Attachment{
		UUID: "a0eebc91-9c0c-4ef7-bb6e-6bb9bd380a12",
	}
	r := &Release{
		Attachments: []*Attachment{a},
	}

	err := InsertReleases(t.Context(), r)
	assert.NoError(t, err)
}

func Test_FindTagsByCommitIDs(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	sha1Rels, err := FindTagsByCommitIDs(t.Context(), 1, "65f1bf27bc3bf70f64657658635e66094edbcb4d")
	assert.NoError(t, err)
	assert.Len(t, sha1Rels, 1)
	rels := sha1Rels["65f1bf27bc3bf70f64657658635e66094edbcb4d"]
	assert.Len(t, rels, 3)
	assert.Equal(t, "v1.1", rels[0].TagName)
	assert.Equal(t, "delete-tag", rels[1].TagName)
	assert.Equal(t, "v1.0", rels[2].TagName)
}
