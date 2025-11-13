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

package repository

import (
	"strings"
	"testing"

	repo_module "github.com/kumose/kmup/modules/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_detectLicense(t *testing.T) {
	type DetectLicenseTest struct {
		name string
		arg  string
		want []string
	}

	tests := []DetectLicenseTest{
		{
			name: "empty",
			arg:  "",
			want: nil,
		},
		{
			name: "no detected license",
			arg:  "Copyright (c) 2023 Kmup",
			want: nil,
		},
	}

	require.NoError(t, repo_module.LoadRepoConfig())
	for _, licenseName := range repo_module.Licenses {
		license, err := repo_module.GetLicense(licenseName, &repo_module.LicenseValues{
			Owner: "Kmup",
			Email: "teabot@kmup.io",
			Repo:  "kmup",
			Year:  "2024",
		})
		assert.NoError(t, err)

		tests = append(tests, DetectLicenseTest{
			name: "single license test: " + licenseName,
			arg:  string(license),
			want: []string{licenseName},
		})
	}

	require.NoError(t, InitLicenseClassifier())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			license, err := detectLicense(strings.NewReader(tt.arg))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, license)
		})
	}

	result, err := detectLicense(strings.NewReader(tests[2].arg + tests[3].arg + tests[4].arg))
	assert.NoError(t, err)
	t.Run("multiple licenses test", func(t *testing.T) {
		assert.Len(t, result, 3)
		assert.Contains(t, result, tests[2].want[0])
		assert.Contains(t, result, tests[3].want[0])
		assert.Contains(t, result, tests[4].want[0])
	})
}
