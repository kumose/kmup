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

package attribute

import (
	"path/filepath"
	"testing"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Checker(t *testing.T) {
	setting.AppDataPath = t.TempDir()
	repoPath := "../tests/repos/language_stats_repo"
	gitRepo, err := git.OpenRepository(t.Context(), repoPath)
	require.NoError(t, err)
	defer gitRepo.Close()

	commitID := "8fee858da5796dfb37704761701bb8e800ad9ef3"

	t.Run("Create index file to run git check-attr", func(t *testing.T) {
		defer test.MockVariableValue(&git.DefaultFeatures().SupportCheckAttrOnBare, false)()
		attrs, err := CheckAttributes(t.Context(), gitRepo, commitID, CheckAttributeOpts{
			Filenames:  []string{"i-am-a-python.p"},
			Attributes: LinguistAttributes,
		})
		assert.NoError(t, err)
		assert.Len(t, attrs, 1)
		assert.Equal(t, expectedAttrs(), attrs["i-am-a-python.p"])
	})

	// run git check-attr on work tree
	t.Run("Run git check-attr on git work tree", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "test-repo")
		err := git.Clone(t.Context(), repoPath, dir, git.CloneRepoOptions{
			Shared: true,
			Branch: "master",
		})
		assert.NoError(t, err)

		tempRepo, err := git.OpenRepository(t.Context(), dir)
		assert.NoError(t, err)
		defer tempRepo.Close()

		attrs, err := CheckAttributes(t.Context(), tempRepo, "", CheckAttributeOpts{
			Filenames:  []string{"i-am-a-python.p"},
			Attributes: LinguistAttributes,
		})
		assert.NoError(t, err)
		assert.Len(t, attrs, 1)
		assert.Equal(t, expectedAttrs(), attrs["i-am-a-python.p"])
	})

	t.Run("Run git check-attr in bare repository using index", func(t *testing.T) {
		attrs, err := CheckAttributes(t.Context(), gitRepo, "", CheckAttributeOpts{
			Filenames:  []string{"i-am-a-python.p"},
			Attributes: LinguistAttributes,
		})
		assert.NoError(t, err)
		assert.Len(t, attrs, 1)
		assert.Equal(t, expectedAttrs(), attrs["i-am-a-python.p"])
	})

	if !git.DefaultFeatures().SupportCheckAttrOnBare {
		t.Skip("git version 2.40 is required to support run check-attr on bare repo without using index")
		return
	}

	t.Run("Run git check-attr in bare repository", func(t *testing.T) {
		attrs, err := CheckAttributes(t.Context(), gitRepo, commitID, CheckAttributeOpts{
			Filenames:  []string{"i-am-a-python.p"},
			Attributes: LinguistAttributes,
		})
		assert.NoError(t, err)
		assert.Len(t, attrs, 1)
		assert.Equal(t, expectedAttrs(), attrs["i-am-a-python.p"])
	})
}
