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
	"context"
	"testing"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestReadingBlameOutputSha256(t *testing.T) {
	setting.AppDataPath = t.TempDir()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	if isGogit {
		t.Skip("Skipping test since gogit does not support sha256")
		return
	}

	t.Run("Without .git-blame-ignore-revs", func(t *testing.T) {
		repo, err := OpenRepository(ctx, "./tests/repos/repo5_pulls_sha256")
		assert.NoError(t, err)
		defer repo.Close()

		commit, err := repo.GetCommit("0b69b7bb649b5d46e14cabb6468685e5dd721290acc7ffe604d37cde57927345")
		assert.NoError(t, err)

		parts := []*BlamePart{
			{
				Sha: "1e35a51dc00fd7de730344c07061acfe80e8117e075ac979b6a29a3a045190ca",
				Lines: []string{
					"# test_repo",
					"Test repository for testing migration from github to kmup",
				},
			},
			{
				Sha:          "0b69b7bb649b5d46e14cabb6468685e5dd721290acc7ffe604d37cde57927345",
				Lines:        []string{"", "Do not make any changes to this repo it is used for unit testing"},
				PreviousSha:  "1e35a51dc00fd7de730344c07061acfe80e8117e075ac979b6a29a3a045190ca",
				PreviousPath: "README.md",
			},
		}

		for _, bypass := range []bool{false, true} {
			blameReader, err := CreateBlameReader(ctx, Sha256ObjectFormat, "./tests/repos/repo5_pulls_sha256", commit, "README.md", bypass)
			assert.NoError(t, err)
			assert.NotNil(t, blameReader)
			defer blameReader.Close()

			assert.False(t, blameReader.UsesIgnoreRevs())

			for _, part := range parts {
				actualPart, err := blameReader.NextPart()
				assert.NoError(t, err)
				assert.Equal(t, part, actualPart)
			}

			// make sure all parts have been read
			actualPart, err := blameReader.NextPart()
			assert.Nil(t, actualPart)
			assert.NoError(t, err)
		}
	})

	t.Run("With .git-blame-ignore-revs", func(t *testing.T) {
		repo, err := OpenRepository(ctx, "./tests/repos/repo6_blame_sha256")
		assert.NoError(t, err)
		defer repo.Close()

		full := []*BlamePart{
			{
				Sha:   "ab2b57a4fa476fb2edb74dafa577caf918561abbaa8fba0c8dc63c412e17a7cc",
				Lines: []string{"line", "line"},
			},
			{
				Sha:          "9347b0198cd1f25017579b79d0938fa89dba34ad2514f0dd92f6bc975ed1a2fe",
				Lines:        []string{"changed line"},
				PreviousSha:  "ab2b57a4fa476fb2edb74dafa577caf918561abbaa8fba0c8dc63c412e17a7cc",
				PreviousPath: "blame.txt",
			},
			{
				Sha:   "ab2b57a4fa476fb2edb74dafa577caf918561abbaa8fba0c8dc63c412e17a7cc",
				Lines: []string{"line", "line", ""},
			},
		}

		cases := []struct {
			CommitID       string
			UsesIgnoreRevs bool
			Bypass         bool
			Parts          []*BlamePart
		}{
			{
				CommitID:       "e2f5660e15159082902960af0ed74fc144921d2b0c80e069361853b3ece29ba3",
				UsesIgnoreRevs: true,
				Bypass:         false,
				Parts: []*BlamePart{
					{
						Sha:   "ab2b57a4fa476fb2edb74dafa577caf918561abbaa8fba0c8dc63c412e17a7cc",
						Lines: []string{"line", "line", "changed line", "line", "line", ""},
					},
				},
			},
			{
				CommitID:       "e2f5660e15159082902960af0ed74fc144921d2b0c80e069361853b3ece29ba3",
				UsesIgnoreRevs: false,
				Bypass:         true,
				Parts:          full,
			},
			{
				CommitID:       "9347b0198cd1f25017579b79d0938fa89dba34ad2514f0dd92f6bc975ed1a2fe",
				UsesIgnoreRevs: false,
				Bypass:         false,
				Parts:          full,
			},
			{
				CommitID:       "9347b0198cd1f25017579b79d0938fa89dba34ad2514f0dd92f6bc975ed1a2fe",
				UsesIgnoreRevs: false,
				Bypass:         false,
				Parts:          full,
			},
		}

		objectFormat, err := repo.GetObjectFormat()
		assert.NoError(t, err)
		for _, c := range cases {
			commit, err := repo.GetCommit(c.CommitID)
			assert.NoError(t, err)
			blameReader, err := CreateBlameReader(ctx, objectFormat, "./tests/repos/repo6_blame_sha256", commit, "blame.txt", c.Bypass)
			assert.NoError(t, err)
			assert.NotNil(t, blameReader)
			defer blameReader.Close()

			assert.Equal(t, c.UsesIgnoreRevs, blameReader.UsesIgnoreRevs())

			for _, part := range c.Parts {
				actualPart, err := blameReader.NextPart()
				assert.NoError(t, err)
				assert.Equal(t, part, actualPart)
			}

			// make sure all parts have been read
			actualPart, err := blameReader.NextPart()
			assert.Nil(t, actualPart)
			assert.NoError(t, err)
		}
	})
}
