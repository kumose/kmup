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

package convert

import (
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestRelease_ToRelease(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	release1 := unittest.AssertExistsAndLoadBean(t, &repo_model.Release{ID: 1})
	release1.LoadAttributes(t.Context())

	apiRelease := ToAPIRelease(t.Context(), repo1, release1)
	assert.NotNil(t, apiRelease)
	assert.EqualValues(t, 1, apiRelease.ID)
	assert.Equal(t, "https://try.kmup.io/api/v1/repos/user2/repo1/releases/1", apiRelease.URL)
	assert.Equal(t, "https://try.kmup.io/api/v1/repos/user2/repo1/releases/1/assets", apiRelease.UploadURL)
}
