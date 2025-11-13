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

package markup

import (
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestResolveLinkRelative(t *testing.T) {
	ctx := t.Context()
	setting.AppURL = "http://localhost:3000"
	assert.Equal(t, "/a", resolveLinkRelative(ctx, "/a", "", "", false))
	assert.Equal(t, "/a/b", resolveLinkRelative(ctx, "/a", "b", "", false))
	assert.Equal(t, "/a/b/c", resolveLinkRelative(ctx, "/a", "b", "c", false))
	assert.Equal(t, "/a/c", resolveLinkRelative(ctx, "/a", "b", "/c", false))
	assert.Equal(t, "http://localhost:3000/a", resolveLinkRelative(ctx, "/a", "", "", true))

	// some users might have used absolute paths a lot, so if the prefix overlaps and has enough slashes, we should tolerate it
	assert.Equal(t, "/owner/repo/foo/owner/repo/foo/bar/xxx", resolveLinkRelative(ctx, "/owner/repo/foo", "", "/owner/repo/foo/bar/xxx", false))
	assert.Equal(t, "/owner/repo/foo/bar/xxx", resolveLinkRelative(ctx, "/owner/repo/foo/bar", "", "/owner/repo/foo/bar/xxx", false))
}
