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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getStorageInheritNameSectionTypeForActions(t *testing.T) {
	iniStr := `
	[storage]
	STORAGE_TYPE = minio
	`
	cfg, err := NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadActionsFrom(cfg))

	assert.EqualValues(t, "minio", Actions.LogStorage.Type)
	assert.Equal(t, "actions_log/", Actions.LogStorage.MinioConfig.BasePath)
	assert.EqualValues(t, "minio", Actions.ArtifactStorage.Type)
	assert.Equal(t, "actions_artifacts/", Actions.ArtifactStorage.MinioConfig.BasePath)

	iniStr = `
[storage.actions_log]
STORAGE_TYPE = minio
`
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadActionsFrom(cfg))

	assert.EqualValues(t, "minio", Actions.LogStorage.Type)
	assert.Equal(t, "actions_log/", Actions.LogStorage.MinioConfig.BasePath)
	assert.EqualValues(t, "local", Actions.ArtifactStorage.Type)
	assert.Equal(t, "actions_artifacts", filepath.Base(Actions.ArtifactStorage.Path))

	iniStr = `
[storage.actions_log]
STORAGE_TYPE = my_storage

[storage.my_storage]
STORAGE_TYPE = minio
`
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadActionsFrom(cfg))

	assert.EqualValues(t, "minio", Actions.LogStorage.Type)
	assert.Equal(t, "actions_log/", Actions.LogStorage.MinioConfig.BasePath)
	assert.EqualValues(t, "local", Actions.ArtifactStorage.Type)
	assert.Equal(t, "actions_artifacts", filepath.Base(Actions.ArtifactStorage.Path))

	iniStr = `
[storage.actions_artifacts]
STORAGE_TYPE = my_storage

[storage.my_storage]
STORAGE_TYPE = minio
`
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadActionsFrom(cfg))

	assert.EqualValues(t, "local", Actions.LogStorage.Type)
	assert.Equal(t, "actions_log", filepath.Base(Actions.LogStorage.Path))
	assert.EqualValues(t, "minio", Actions.ArtifactStorage.Type)
	assert.Equal(t, "actions_artifacts/", Actions.ArtifactStorage.MinioConfig.BasePath)

	iniStr = `
[storage.actions_artifacts]
STORAGE_TYPE = my_storage

[storage.my_storage]
STORAGE_TYPE = minio
`
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadActionsFrom(cfg))

	assert.EqualValues(t, "local", Actions.LogStorage.Type)
	assert.Equal(t, "actions_log", filepath.Base(Actions.LogStorage.Path))
	assert.EqualValues(t, "minio", Actions.ArtifactStorage.Type)
	assert.Equal(t, "actions_artifacts/", Actions.ArtifactStorage.MinioConfig.BasePath)

	iniStr = ``
	cfg, err = NewConfigProviderFromData(iniStr)
	assert.NoError(t, err)
	assert.NoError(t, loadActionsFrom(cfg))

	assert.EqualValues(t, "local", Actions.LogStorage.Type)
	assert.Equal(t, "actions_log", filepath.Base(Actions.LogStorage.Path))
	assert.EqualValues(t, "local", Actions.ArtifactStorage.Type)
	assert.Equal(t, "actions_artifacts", filepath.Base(Actions.ArtifactStorage.Path))
}

func Test_getDefaultActionsURLForActions(t *testing.T) {
	oldActions := Actions
	oldAppURL := AppURL
	defer func() {
		Actions = oldActions
		AppURL = oldAppURL
	}()

	AppURL = "http://test_get_default_actions_url_for_actions:3000/"

	tests := []struct {
		name    string
		iniStr  string
		wantErr assert.ErrorAssertionFunc
		wantURL string
	}{
		{
			name: "default",
			iniStr: `
[actions]
`,
			wantErr: assert.NoError,
			wantURL: "https://github.com",
		},
		{
			name: "github",
			iniStr: `
[actions]
DEFAULT_ACTIONS_URL = github
`,
			wantErr: assert.NoError,
			wantURL: "https://github.com",
		},
		{
			name: "self",
			iniStr: `
[actions]
DEFAULT_ACTIONS_URL = self
`,
			wantErr: assert.NoError,
			wantURL: "http://test_get_default_actions_url_for_actions:3000",
		},
		{
			name: "custom url",
			iniStr: `
[actions]
DEFAULT_ACTIONS_URL = https://kmup.com
`,
			wantErr: assert.NoError,
			wantURL: "https://github.com",
		},
		{
			name: "custom urls",
			iniStr: `
[actions]
DEFAULT_ACTIONS_URL = https://kmup.com,https://github.com
`,
			wantErr: assert.NoError,
			wantURL: "https://github.com",
		},
		{
			name: "invalid",
			iniStr: `
[actions]
DEFAULT_ACTIONS_URL = kmup
`,
			wantErr: assert.Error,
			wantURL: "https://github.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConfigProviderFromData(tt.iniStr)
			require.NoError(t, err)
			if !tt.wantErr(t, loadActionsFrom(cfg)) {
				return
			}
			assert.Equal(t, tt.wantURL, Actions.DefaultActionsURL.URL())
		})
	}
}
