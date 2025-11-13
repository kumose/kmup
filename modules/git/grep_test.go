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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrepSearch(t *testing.T) {
	repo, err := OpenRepository(t.Context(), filepath.Join(testReposDir, "language_stats_repo"))
	assert.NoError(t, err)
	defer repo.Close()

	res, err := GrepSearch(t.Context(), repo, "void", GrepOptions{})
	assert.NoError(t, err)
	assert.Equal(t, []*GrepResult{
		{
			Filename:    "java-hello/main.java",
			LineNumbers: []int{3},
			LineCodes:   []string{" public static void main(String[] args)"},
		},
		{
			Filename:    "main.vendor.java",
			LineNumbers: []int{3},
			LineCodes:   []string{" public static void main(String[] args)"},
		},
	}, res)

	res, err = GrepSearch(t.Context(), repo, "void", GrepOptions{PathspecList: []string{":(glob)java-hello/*"}})
	assert.NoError(t, err)
	assert.Equal(t, []*GrepResult{
		{
			Filename:    "java-hello/main.java",
			LineNumbers: []int{3},
			LineCodes:   []string{" public static void main(String[] args)"},
		},
	}, res)

	res, err = GrepSearch(t.Context(), repo, "void", GrepOptions{PathspecList: []string{":(glob,exclude)java-hello/*"}})
	assert.NoError(t, err)
	assert.Equal(t, []*GrepResult{
		{
			Filename:    "main.vendor.java",
			LineNumbers: []int{3},
			LineCodes:   []string{" public static void main(String[] args)"},
		},
	}, res)

	res, err = GrepSearch(t.Context(), repo, "void", GrepOptions{MaxResultLimit: 1})
	assert.NoError(t, err)
	assert.Equal(t, []*GrepResult{
		{
			Filename:    "java-hello/main.java",
			LineNumbers: []int{3},
			LineCodes:   []string{" public static void main(String[] args)"},
		},
	}, res)

	res, err = GrepSearch(t.Context(), repo, "void", GrepOptions{MaxResultLimit: 1, MaxLineLength: 39})
	assert.NoError(t, err)
	assert.Equal(t, []*GrepResult{
		{
			Filename:    "java-hello/main.java",
			LineNumbers: []int{3},
			LineCodes:   []string{" public static void main(String[] arg"},
		},
	}, res)

	res, err = GrepSearch(t.Context(), repo, "no-such-content", GrepOptions{})
	assert.NoError(t, err)
	assert.Empty(t, res)

	res, err = GrepSearch(t.Context(), &Repository{Path: "no-such-git-repo"}, "no-such-content", GrepOptions{})
	assert.Error(t, err)
	assert.Empty(t, res)
}
