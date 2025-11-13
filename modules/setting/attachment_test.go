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

func Test_getStorageCustomType(t *testing.T) {
	iniStr := `
[attachment]
STORAGE_TYPE = my_minio
MINIO_BUCKET = kmup-attachment

[storage.my_minio]
STORAGE_TYPE = minio
MINIO_ENDPOINT = my_minio:9000
`
	cfg, err := NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)

	assert.NoError(t, loadAttachmentFrom(cfg))

	assert.EqualValues(t, "minio", Attachment.Storage.Type)
	assert.Equal(t, "my_minio:9000", Attachment.Storage.MinioConfig.Endpoint)
	assert.Equal(t, "kmup-attachment", Attachment.Storage.MinioConfig.Bucket)
	assert.Equal(t, "attachments/", Attachment.Storage.MinioConfig.BasePath)
}

func Test_getStorageTypeSectionOverridesStorageSection(t *testing.T) {
	iniStr := `
[attachment]
STORAGE_TYPE = minio

[storage.minio]
MINIO_BUCKET = kmup-minio

[storage]
MINIO_BUCKET = kmup
`
	cfg, err := NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)

	assert.NoError(t, loadAttachmentFrom(cfg))

	assert.EqualValues(t, "minio", Attachment.Storage.Type)
	assert.Equal(t, "kmup-minio", Attachment.Storage.MinioConfig.Bucket)
	assert.Equal(t, "attachments/", Attachment.Storage.MinioConfig.BasePath)
}

func Test_getStorageSpecificOverridesStorage(t *testing.T) {
	iniStr := `
[attachment]
STORAGE_TYPE = minio
MINIO_BUCKET = kmup-attachment

[storage.attachments]
MINIO_BUCKET = kmup

[storage]
STORAGE_TYPE = local
`
	cfg, err := NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)

	assert.NoError(t, loadAttachmentFrom(cfg))

	assert.EqualValues(t, "minio", Attachment.Storage.Type)
	assert.Equal(t, "kmup-attachment", Attachment.Storage.MinioConfig.Bucket)
	assert.Equal(t, "attachments/", Attachment.Storage.MinioConfig.BasePath)
}

func Test_getStorageGetDefaults(t *testing.T) {
	cfg, err := NewConfigProviderFromData("")
	assert.NoError(t, err)

	assert.NoError(t, loadAttachmentFrom(cfg))

	// default storage is local, so bucket is empty
	assert.Empty(t, Attachment.Storage.MinioConfig.Bucket)
}

func Test_getStorageInheritNameSectionType(t *testing.T) {
	iniStr := `
[storage.attachments]
STORAGE_TYPE = minio
`
	cfg, err := NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)

	assert.NoError(t, loadAttachmentFrom(cfg))

	assert.EqualValues(t, "minio", Attachment.Storage.Type)
}

func Test_AttachmentStorage(t *testing.T) {
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

	assert.NoError(t, loadAttachmentFrom(cfg))
	storage := Attachment.Storage

	assert.EqualValues(t, "minio", storage.Type)
	assert.Equal(t, "kmup", storage.MinioConfig.Bucket)
}

func Test_AttachmentStorage1(t *testing.T) {
	iniStr := `
[storage]
STORAGE_TYPE = minio
`
	cfg, err := NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)

	assert.NoError(t, loadAttachmentFrom(cfg))
	assert.EqualValues(t, "minio", Attachment.Storage.Type)
	assert.Equal(t, "kmup", Attachment.Storage.MinioConfig.Bucket)
	assert.Equal(t, "attachments/", Attachment.Storage.MinioConfig.BasePath)
}
