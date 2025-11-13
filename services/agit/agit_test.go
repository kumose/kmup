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

package agit

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/util"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

func TestParseAgitPushOptionValue(t *testing.T) {
	assert.Equal(t, "a", parseAgitPushOptionValue("a"))
	assert.Equal(t, "a", parseAgitPushOptionValue("{base64}YQ=="))
	assert.Equal(t, "{base64}invalid value", parseAgitPushOptionValue("{base64}invalid value"))
}

func TestGetAgitBranchInfo(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	_, _, err := GetAgitBranchInfo(t.Context(), 1, "non-exist-basebranch")
	assert.ErrorIs(t, err, util.ErrNotExist)

	baseBranch, currentTopicBranch, err := GetAgitBranchInfo(t.Context(), 1, "master")
	assert.NoError(t, err)
	assert.Equal(t, "master", baseBranch)
	assert.Empty(t, currentTopicBranch)

	baseBranch, currentTopicBranch, err = GetAgitBranchInfo(t.Context(), 1, "master/topicbranch")
	assert.NoError(t, err)
	assert.Equal(t, "master", baseBranch)
	assert.Equal(t, "topicbranch", currentTopicBranch)

	baseBranch, currentTopicBranch, err = GetAgitBranchInfo(t.Context(), 1, "master/")
	assert.NoError(t, err)
	assert.Equal(t, "master", baseBranch)
	assert.Empty(t, currentTopicBranch)

	_, _, err = GetAgitBranchInfo(t.Context(), 1, "/")
	assert.ErrorIs(t, err, util.ErrNotExist)

	_, _, err = GetAgitBranchInfo(t.Context(), 1, "//")
	assert.ErrorIs(t, err, util.ErrNotExist)

	baseBranch, currentTopicBranch, err = GetAgitBranchInfo(t.Context(), 1, "master/topicbranch/")
	assert.NoError(t, err)
	assert.Equal(t, "master", baseBranch)
	assert.Equal(t, "topicbranch/", currentTopicBranch)

	baseBranch, currentTopicBranch, err = GetAgitBranchInfo(t.Context(), 1, "master/topicbranch/1")
	assert.NoError(t, err)
	assert.Equal(t, "master", baseBranch)
	assert.Equal(t, "topicbranch/1", currentTopicBranch)
}
