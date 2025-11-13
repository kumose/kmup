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
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_GetBlob_Found(t *testing.T) {
	repoPath := filepath.Join(testReposDir, "repo1_bare")
	r, err := OpenRepository(t.Context(), repoPath)
	assert.NoError(t, err)
	defer r.Close()

	testCases := []struct {
		OID  string
		Data []byte
	}{
		{"e2129701f1a4d54dc44f03c93bca0a2aec7c5449", []byte("file1\n")},
		{"6c493ff740f9380390d5c9ddef4af18697ac9375", []byte("file2\n")},
	}

	for _, testCase := range testCases {
		blob, err := r.GetBlob(testCase.OID)
		assert.NoError(t, err)

		dataReader, err := blob.DataAsync()
		assert.NoError(t, err)

		data, err := io.ReadAll(dataReader)
		assert.NoError(t, dataReader.Close())
		assert.NoError(t, err)
		assert.Equal(t, testCase.Data, data)
	}
}

func TestRepository_GetBlob_NotExist(t *testing.T) {
	repoPath := filepath.Join(testReposDir, "repo1_bare")
	r, err := OpenRepository(t.Context(), repoPath)
	assert.NoError(t, err)
	defer r.Close()

	testCase := "0000000000000000000000000000000000000000"
	testError := ErrNotExist{testCase, ""}

	blob, err := r.GetBlob(testCase)
	assert.Nil(t, blob)
	assert.EqualError(t, err, testError.Error())
}

func TestRepository_GetBlob_NoId(t *testing.T) {
	repoPath := filepath.Join(testReposDir, "repo1_bare")
	r, err := OpenRepository(t.Context(), repoPath)
	assert.NoError(t, err)
	defer r.Close()

	testCase := ""
	testError := fmt.Errorf("length %d has no matched object format: %s", len(testCase), testCase)

	blob, err := r.GetBlob(testCase)
	assert.Nil(t, blob)
	assert.EqualError(t, err, testError.Error())
}
