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
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestInTransaction(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	assert.False(t, db.InTransaction(t.Context()))
	assert.NoError(t, db.WithTx(t.Context(), func(ctx context.Context) error {
		assert.True(t, db.InTransaction(ctx))
		return nil
	}))

	ctx, committer, err := db.TxContext(t.Context())
	assert.NoError(t, err)
	defer committer.Close()
	assert.True(t, db.InTransaction(ctx))
	assert.NoError(t, db.WithTx(ctx, func(ctx context.Context) error {
		assert.True(t, db.InTransaction(ctx))
		return nil
	}))
}

func TestTxContext(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	{ // create new transaction
		ctx, committer, err := db.TxContext(t.Context())
		assert.NoError(t, err)
		assert.True(t, db.InTransaction(ctx))
		assert.NoError(t, committer.Commit())
	}

	{ // reuse the transaction created by TxContext and commit it
		ctx, committer, err := db.TxContext(t.Context())
		engine := db.GetEngine(ctx)
		assert.NoError(t, err)
		assert.True(t, db.InTransaction(ctx))
		{
			ctx, committer, err := db.TxContext(ctx)
			assert.NoError(t, err)
			assert.True(t, db.InTransaction(ctx))
			assert.Equal(t, engine, db.GetEngine(ctx))
			assert.NoError(t, committer.Commit())
		}
		assert.NoError(t, committer.Commit())
	}

	{ // reuse the transaction created by TxContext and close it
		ctx, committer, err := db.TxContext(t.Context())
		engine := db.GetEngine(ctx)
		assert.NoError(t, err)
		assert.True(t, db.InTransaction(ctx))
		{
			ctx, committer, err := db.TxContext(ctx)
			assert.NoError(t, err)
			assert.True(t, db.InTransaction(ctx))
			assert.Equal(t, engine, db.GetEngine(ctx))
			assert.NoError(t, committer.Close())
		}
		assert.NoError(t, committer.Close())
	}

	{ // reuse the transaction created by WithTx
		assert.NoError(t, db.WithTx(t.Context(), func(ctx context.Context) error {
			assert.True(t, db.InTransaction(ctx))
			{
				ctx, committer, err := db.TxContext(ctx)
				assert.NoError(t, err)
				assert.True(t, db.InTransaction(ctx))
				assert.NoError(t, committer.Commit())
			}
			return nil
		}))
	}
}

func TestContextSafety(t *testing.T) {
	type TestModel1 struct {
		ID int64
	}
	type TestModel2 struct {
		ID int64
	}
	assert.NoError(t, unittest.GetXORMEngine().Sync(&TestModel1{}, &TestModel2{}))
	assert.NoError(t, db.TruncateBeans(t.Context(), &TestModel1{}, &TestModel2{}))
	testCount := 10
	for i := 1; i <= testCount; i++ {
		assert.NoError(t, db.Insert(t.Context(), &TestModel1{ID: int64(i)}))
		assert.NoError(t, db.Insert(t.Context(), &TestModel2{ID: int64(-i)}))
	}

	t.Run("Show-XORM-Bug", func(t *testing.T) {
		actualCount := 0
		// here: db.GetEngine(t.Context()) is a new *Session created from *Engine
		_ = db.WithTx(t.Context(), func(ctx context.Context) error {
			_ = db.GetEngine(ctx).Iterate(&TestModel1{}, func(i int, bean any) error {
				// here: db.GetEngine(ctx) is always the unclosed "Iterate" *Session with autoResetStatement=false,
				// and the internal states (including "cond" and others) are always there and not be reset in this callback.
				m1 := bean.(*TestModel1)
				assert.EqualValues(t, i+1, m1.ID)

				// here: XORM bug, it fails because the SQL becomes "WHERE id=-1", "WHERE id=-1 AND id=-2", "WHERE id=-1 AND id=-2 AND id=-3" ...
				// and it conflicts with the "Iterate"'s internal states.
				// has, err := db.GetEngine(ctx).Get(&TestModel2{ID: -m1.ID})

				actualCount++
				return nil
			})
			return nil
		})
		assert.Equal(t, testCount, actualCount)
	})

	t.Run("DenyBadUsage", func(t *testing.T) {
		assert.PanicsWithError(t, "using session context in an iterator would cause corrupted results", func() {
			_ = db.WithTx(t.Context(), func(ctx context.Context) error {
				return db.GetEngine(ctx).Iterate(&TestModel1{}, func(i int, bean any) error {
					_ = db.GetEngine(ctx)
					return nil
				})
			})
		})
	})
}
