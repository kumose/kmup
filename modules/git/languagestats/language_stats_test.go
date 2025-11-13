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

//go:build !gogit

package languagestats

import (
	"testing"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetLanguageStats(t *testing.T) {
	setting.AppDataPath = t.TempDir()
	repoPath := "../tests/repos/language_stats_repo"
	gitRepo, err := git.OpenRepository(t.Context(), repoPath)
	require.NoError(t, err)
	defer gitRepo.Close()

	stats, err := GetLanguageStats(gitRepo, "8fee858da5796dfb37704761701bb8e800ad9ef3")
	require.NoError(t, err)

	assert.Equal(t, map[string]int64{
		"Python": 134,
		"Java":   112,
	}, stats)
}

func TestMergeLanguageStats(t *testing.T) {
	assert.Equal(t, map[string]int64{
		"PHP":    1,
		"python": 10,
		"JAVA":   700,
	}, mergeLanguageStats(map[string]int64{
		"PHP":    1,
		"python": 10,
		"Java":   100,
		"java":   200,
		"JAVA":   400,
	}))
}
