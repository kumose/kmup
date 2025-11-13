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

package setting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getStorageInheritNameSectionTypeForRepoArchive(t *testing.T) {
	// packages storage inherits from storage if nothing configured
	iniStr := `
[storage]
STORAGE_TYPE = minio
`
	cfg, err := NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadRepoArchiveFrom(cfg))

	assert.EqualValues(t, "minio", RepoArchive.Storage.Type)
	assert.Equal(t, "repo-archive/", RepoArchive.Storage.MinioConfig.BasePath)

	// we can also configure packages storage directly
	iniStr = `
[storage.repo-archive]
STORAGE_TYPE = minio
`
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadRepoArchiveFrom(cfg))

	assert.EqualValues(t, "minio", RepoArchive.Storage.Type)
	assert.Equal(t, "repo-archive/", RepoArchive.Storage.MinioConfig.BasePath)

	// or we can indicate the storage type in the packages section
	iniStr = `
[repo-archive]
STORAGE_TYPE = my_minio

[storage.my_minio]
STORAGE_TYPE = minio
`
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadRepoArchiveFrom(cfg))

	assert.EqualValues(t, "minio", RepoArchive.Storage.Type)
	assert.Equal(t, "repo-archive/", RepoArchive.Storage.MinioConfig.BasePath)

	// or we can indicate the storage type  and minio base path in the packages section
	iniStr = `
[repo-archive]
STORAGE_TYPE = my_minio
MINIO_BASE_PATH = my_archive/

[storage.my_minio]
STORAGE_TYPE = minio
`
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadRepoArchiveFrom(cfg))

	assert.EqualValues(t, "minio", RepoArchive.Storage.Type)
	assert.Equal(t, "my_archive/", RepoArchive.Storage.MinioConfig.BasePath)
}

func Test_RepoArchiveStorage(t *testing.T) {
	iniStr := `
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
[storage]
STORAGE_TYPE            = minio
MINIO_ENDPOINT          = s3.my-domain.net
MINIO_BUCKET            = kmup
MINIO_LOCATION          = homenet
MINIO_USE_SSL           = true
MINIO_ACCESS_KEY_ID     = correct_key
MINIO_SECRET_ACCESS_KEY = correct_key
`
	cfg, err := NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)

	assert.NoError(t, loadRepoArchiveFrom(cfg))
	storage := RepoArchive.Storage

	assert.EqualValues(t, "minio", storage.Type)
	assert.Equal(t, "kmup", storage.MinioConfig.Bucket)

	iniStr = `
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
[storage.repo-archive]
STORAGE_TYPE = s3
[storage.s3]
STORAGE_TYPE            = minio
MINIO_ENDPOINT          = s3.my-domain.net
MINIO_BUCKET            = kmup
MINIO_LOCATION          = homenet
MINIO_USE_SSL           = true
MINIO_ACCESS_KEY_ID     = correct_key
MINIO_SECRET_ACCESS_KEY = correct_key
`
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)

	assert.NoError(t, loadRepoArchiveFrom(cfg))
	storage = RepoArchive.Storage

	assert.EqualValues(t, "minio", storage.Type)
	assert.Equal(t, "kmup", storage.MinioConfig.Bucket)
}
