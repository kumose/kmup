// Copyright 2015 The Gogs Authors. All rights reserved.
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
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlob_Data(t *testing.T) {
	output := "file2\n"
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	repo, err := OpenRepository(t.Context(), bareRepo1Path)
	require.NoError(t, err)
	defer repo.Close()

	testBlob, err := repo.GetBlob("6c493ff740f9380390d5c9ddef4af18697ac9375")
	assert.NoError(t, err)

	r, err := testBlob.DataAsync()
	assert.NoError(t, err)
	require.NotNil(t, r)

	data, err := io.ReadAll(r)
	assert.NoError(t, r.Close())

	assert.NoError(t, err)
	assert.Equal(t, output, string(data))
}

func Benchmark_Blob_Data(b *testing.B) {
	bareRepo1Path := filepath.Join(testReposDir, "repo1_bare")
	repo, err := OpenRepository(b.Context(), bareRepo1Path)
	if err != nil {
		b.Fatal(err)
	}
	defer repo.Close()

	testBlob, err := repo.GetBlob("6c493ff740f9380390d5c9ddef4af18697ac9375")
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		r, err := testBlob.DataAsync()
		if err != nil {
			b.Fatal(err)
		}
		io.ReadAll(r)
		_ = r.Close()
	}
}
