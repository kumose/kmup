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

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefName(t *testing.T) {
	// Test branch names (with and without slash).
	assert.Equal(t, "foo", RefName("refs/heads/foo").BranchName())
	assert.Equal(t, "feature/foo", RefName("refs/heads/feature/foo").BranchName())

	// Test tag names (with and without slash).
	assert.Equal(t, "foo", RefName("refs/tags/foo").TagName())
	assert.Equal(t, "release/foo", RefName("refs/tags/release/foo").TagName())

	// Test pull names
	assert.Equal(t, "1", RefName("refs/pull/1/head").PullName())
	assert.True(t, RefName("refs/pull/1/head").IsPull())
	assert.True(t, RefName("refs/pull/1/merge").IsPull())
	assert.Equal(t, "my/pull", RefName("refs/pull/my/pull/head").PullName())

	// Test for branch names
	assert.Equal(t, "master", RefName("refs/for/master").ForBranchName())
	assert.Equal(t, "my/branch", RefName("refs/for/my/branch").ForBranchName())

	// Test commit hashes.
	assert.Equal(t, "c0ffee", RefName("c0ffee").ShortName())
}

func TestRefWebLinkPath(t *testing.T) {
	assert.Equal(t, "branch/foo", RefName("refs/heads/foo").RefWebLinkPath())
	assert.Equal(t, "tag/foo", RefName("refs/tags/foo").RefWebLinkPath())
	assert.Equal(t, "commit/c0ffee", RefName("c0ffee").RefWebLinkPath())
}
