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

package v1_23

import (
	"testing"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/timeutil"

	"github.com/stretchr/testify/assert"
)

func Test_AddIndexForReleaseSha1(t *testing.T) {
	type Release struct {
		ID               int64  `xorm:"pk autoincr"`
		RepoID           int64  `xorm:"INDEX UNIQUE(n)"`
		PublisherID      int64  `xorm:"INDEX"`
		TagName          string `xorm:"INDEX UNIQUE(n)"`
		OriginalAuthor   string
		OriginalAuthorID int64 `xorm:"index"`
		LowerTagName     string
		Target           string
		Title            string
		Sha1             string `xorm:"VARCHAR(64)"`
		NumCommits       int64
		Note             string             `xorm:"TEXT"`
		IsDraft          bool               `xorm:"NOT NULL DEFAULT false"`
		IsPrerelease     bool               `xorm:"NOT NULL DEFAULT false"`
		IsTag            bool               `xorm:"NOT NULL DEFAULT false"` // will be true only if the record is a tag and has no related releases
		CreatedUnix      timeutil.TimeStamp `xorm:"INDEX"`
	}

	// Prepare and load the testing database
	x, deferable := base.PrepareTestEnv(t, 0, new(Release))
	defer deferable()

	assert.NoError(t, AddIndexForReleaseSha1(x))
}
