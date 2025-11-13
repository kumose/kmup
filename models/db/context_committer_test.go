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

package db // it's not db_test, because this file is for testing the private type halfCommitter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCommitter struct {
	wants []string
	gots  []string
}

func NewMockCommitter(wants ...string) *MockCommitter {
	return &MockCommitter{
		wants: wants,
	}
}

func (c *MockCommitter) Commit() error {
	c.gots = append(c.gots, "commit")
	return nil
}

func (c *MockCommitter) Close() error {
	c.gots = append(c.gots, "close")
	return nil
}

func (c *MockCommitter) Assert(t *testing.T) {
	assert.Equal(t, c.wants, c.gots, "want operations %v, but got %v", c.wants, c.gots)
}

func Test_halfCommitter(t *testing.T) {
	/*
		Do something like:

		ctx, committer, err := db.TxContext(t.Context())
		if err != nil {
			return nil
		}
		defer committer.Close()

		// ...

		if err != nil {
			return nil
		}

		// ...

		return committer.Commit()
	*/

	testWithCommitter := func(committer Committer, f func(committer Committer) error) {
		if err := f(&halfCommitter{committer: committer}); err == nil {
			committer.Commit()
		}
		committer.Close()
	}

	t.Run("commit and close", func(t *testing.T) {
		mockCommitter := NewMockCommitter("commit", "close")

		testWithCommitter(mockCommitter, func(committer Committer) error {
			defer committer.Close()
			return committer.Commit()
		})

		mockCommitter.Assert(t)
	})

	t.Run("rollback and close", func(t *testing.T) {
		mockCommitter := NewMockCommitter("close", "close")

		testWithCommitter(mockCommitter, func(committer Committer) error {
			defer committer.Close()
			if true {
				return errors.New("error")
			}
			return committer.Commit()
		})

		mockCommitter.Assert(t)
	})

	t.Run("close and commit", func(t *testing.T) {
		mockCommitter := NewMockCommitter("close", "close")

		testWithCommitter(mockCommitter, func(committer Committer) error {
			committer.Close()
			committer.Commit()
			return errors.New("error")
		})

		mockCommitter.Assert(t)
	})
}
