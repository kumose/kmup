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

package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/kumose/kmup/models/packages"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	packages_module "github.com/kumose/kmup/modules/packages"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/storage"
	packages_service "github.com/kumose/kmup/services/packages"

	"github.com/stretchr/testify/assert"
)

func TestMigratePackages(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	creator := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})

	content := "package main\n\nfunc main() {\nfmt.Println(\"hi\")\n}\n"
	buf, err := packages_module.CreateHashedBufferFromReaderWithSize(strings.NewReader(content), 1024)
	assert.NoError(t, err)
	defer buf.Close()

	v, f, err := packages_service.CreatePackageAndAddFile(t.Context(), &packages_service.PackageCreationInfo{
		PackageInfo: packages_service.PackageInfo{
			Owner:       creator,
			PackageType: packages.TypeGeneric,
			Name:        "test",
			Version:     "1.0.0",
		},
		Creator:           creator,
		SemverCompatible:  true,
		VersionProperties: map[string]string{},
	}, &packages_service.PackageFileCreationInfo{
		PackageFileInfo: packages_service.PackageFileInfo{
			Filename: "a.go",
		},
		Creator: creator,
		Data:    buf,
		IsLead:  true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.NotNil(t, f)

	ctx := t.Context()

	p := t.TempDir()

	dstStorage, err := storage.NewLocalStorage(
		ctx,
		&setting.Storage{
			Path: p,
		})
	assert.NoError(t, err)

	err = migratePackages(ctx, dstStorage)
	assert.NoError(t, err)

	entries, err := os.ReadDir(p)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "01", entries[0].Name())
	assert.Equal(t, "tmp", entries[1].Name())
}
