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

package db_test

import (
	"context"
	"testing"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestIterate(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	xe := unittest.GetXORMEngine()
	assert.NoError(t, xe.Sync(&repo_model.RepoUnit{}))

	cnt, err := db.GetEngine(t.Context()).Count(&repo_model.RepoUnit{})
	assert.NoError(t, err)

	var repoUnitCnt int
	err = db.Iterate(t.Context(), nil, func(ctx context.Context, repo *repo_model.RepoUnit) error {
		repoUnitCnt++
		return nil
	})
	assert.NoError(t, err)
	assert.EqualValues(t, cnt, repoUnitCnt)

	err = db.Iterate(t.Context(), nil, func(ctx context.Context, repoUnit *repo_model.RepoUnit) error {
		has, err := db.ExistByID[repo_model.RepoUnit](ctx, repoUnit.ID)
		if err != nil {
			return err
		}
		if !has {
			return db.ErrNotExist{Resource: "repo_unit", ID: repoUnit.ID}
		}
		return nil
	})
	assert.NoError(t, err)
}
