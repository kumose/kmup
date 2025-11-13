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

package packages_test

import (
	"testing"

	packages_model "github.com/kumose/kmup/models/packages"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	_ "github.com/kumose/kmup/models"
	_ "github.com/kumose/kmup/models/actions"
	_ "github.com/kumose/kmup/models/activities"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}

func TestHasOwnerPackages(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})

	p, err := packages_model.TryInsertPackage(t.Context(), &packages_model.Package{
		OwnerID:   owner.ID,
		LowerName: "package",
	})
	assert.NotNil(t, p)
	assert.NoError(t, err)

	// A package without package versions gets automatically cleaned up and should return false
	has, err := packages_model.HasOwnerPackages(t.Context(), owner.ID)
	assert.False(t, has)
	assert.NoError(t, err)

	pv, err := packages_model.GetOrInsertVersion(t.Context(), &packages_model.PackageVersion{
		PackageID:    p.ID,
		LowerVersion: "internal",
		IsInternal:   true,
	})
	assert.NotNil(t, pv)
	assert.NoError(t, err)

	// A package with an internal package version gets automatically cleaned up and should return false
	has, err = packages_model.HasOwnerPackages(t.Context(), owner.ID)
	assert.False(t, has)
	assert.NoError(t, err)

	pv, err = packages_model.GetOrInsertVersion(t.Context(), &packages_model.PackageVersion{
		PackageID:    p.ID,
		LowerVersion: "normal",
		IsInternal:   false,
	})
	assert.NotNil(t, pv)
	assert.NoError(t, err)

	// A package with a normal package version should return true
	has, err = packages_model.HasOwnerPackages(t.Context(), owner.ID)
	assert.True(t, has)
	assert.NoError(t, err)
}
