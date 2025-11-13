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
	"path/filepath"
	"testing"

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/setting"

	_ "github.com/kumose/kmup/cmd" // for TestPrimaryKeys

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDumpDatabase(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	dir := t.TempDir()

	type Version struct {
		ID      int64 `xorm:"pk autoincr"`
		Version int64
	}
	assert.NoError(t, db.GetEngine(t.Context()).Sync(new(Version)))

	for _, dbType := range setting.SupportedDatabaseTypes {
		assert.NoError(t, db.DumpDatabase(filepath.Join(dir, dbType+".sql"), dbType))
	}
}

func TestDeleteOrphanedObjects(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	countBefore, err := db.GetEngine(t.Context()).Count(&issues_model.PullRequest{})
	assert.NoError(t, err)

	_, err = db.GetEngine(t.Context()).Insert(&issues_model.PullRequest{IssueID: 1000}, &issues_model.PullRequest{IssueID: 1001}, &issues_model.PullRequest{IssueID: 1003})
	assert.NoError(t, err)

	orphaned, err := db.CountOrphanedObjects(t.Context(), "pull_request", "issue", "pull_request.issue_id=issue.id")
	assert.NoError(t, err)
	assert.EqualValues(t, 3, orphaned)

	err = db.DeleteOrphanedObjects(t.Context(), "pull_request", "issue", "pull_request.issue_id=issue.id")
	assert.NoError(t, err)

	countAfter, err := db.GetEngine(t.Context()).Count(&issues_model.PullRequest{})
	assert.NoError(t, err)
	assert.Equal(t, countBefore, countAfter)
}

func TestPrimaryKeys(t *testing.T) {
	// Some dbs require that all tables have primary keys, see
	// To avoid creating tables without primary key again, this test will check them.
	// Import "github.com/kumose/kmup/cmd" to make sure each db.RegisterModel in init functions has been called.

	beans, err := db.NamesToBean()
	require.NoError(t, err)

	whitelist := map[string]string{
		"the_table_name_to_skip_checking": "Write a note here to explain why",
	}

	for _, bean := range beans {
		table, err := db.GetXORMEngineForTesting().TableInfo(bean)
		if err != nil {
			t.Fatal(err)
		}
		if why, ok := whitelist[table.Name]; ok {
			t.Logf("ignore %q because %q", table.Name, why)
			continue
		}
		assert.NotEmpty(t, table.PrimaryKeys, "table %q has no primary key", table.Name)
	}
}
