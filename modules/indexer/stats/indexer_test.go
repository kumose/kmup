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

package stats

import (
	"testing"
	"time"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/queue"
	"github.com/kumose/kmup/modules/setting"

	_ "github.com/kumose/kmup/models"
	_ "github.com/kumose/kmup/models/actions"
	_ "github.com/kumose/kmup/models/activities"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

func TestRepoStatsIndex(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	setting.CfgProvider, _ = setting.NewConfigProviderFromData("")

	setting.LoadQueueSettings()

	err := Init()
	assert.NoError(t, err)

	repo, err := repo_model.GetRepositoryByID(t.Context(), 1)
	assert.NoError(t, err)

	err = UpdateRepoIndexer(repo)
	assert.NoError(t, err)

	assert.NoError(t, queue.GetManager().FlushAll(t.Context(), 5*time.Second))

	status, err := repo_model.GetIndexerStatus(t.Context(), repo, repo_model.RepoIndexerTypeStats)
	assert.NoError(t, err)
	assert.Equal(t, "65f1bf27bc3bf70f64657658635e66094edbcb4d", status.CommitSha)
	langs, err := repo_model.GetTopLanguageStats(t.Context(), repo, 5)
	assert.NoError(t, err)
	assert.Empty(t, langs)
}
