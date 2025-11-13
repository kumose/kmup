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

package repo

import (
	"html/template"
	"testing"

	pull_model "github.com/kumose/kmup/models/pull"
	"github.com/kumose/kmup/modules/fileicon"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/services/gitdiff"

	"github.com/stretchr/testify/assert"
)

func TestTransformDiffTreeForWeb(t *testing.T) {
	renderedIconPool := fileicon.NewRenderedIconPool()
	ret := transformDiffTreeForWeb(renderedIconPool, &gitdiff.DiffTree{Files: []*gitdiff.DiffTreeRecord{
		{
			Status:   "changed",
			HeadPath: "dir-a/dir-a-x/file-deep",
			HeadMode: git.EntryModeBlob,
		},
		{
			Status:   "added",
			HeadPath: "file1",
			HeadMode: git.EntryModeBlob,
		},
	}}, map[string]pull_model.ViewedState{
		"dir-a/dir-a-x/file-deep": pull_model.Viewed,
	})

	mockIconForFile := func(id string) template.HTML {
		return template.HTML(`<svg class="svg git-entry-icon octicon-file" width="16" height="16" aria-hidden="true"><use xlink:href="#` + id + `"></use></svg>`)
	}
	assert.Equal(t, WebDiffFileTree{
		TreeRoot: WebDiffFileItem{
			Children: []*WebDiffFileItem{
				{
					EntryMode:   "tree",
					DisplayName: "dir-a/dir-a-x",
					FullName:    "dir-a/dir-a-x",
					Children: []*WebDiffFileItem{
						{
							EntryMode:   "",
							DisplayName: "file-deep",
							FullName:    "dir-a/dir-a-x/file-deep",
							NameHash:    "4acf7eef1c943a09e9f754e93ff190db8583236b",
							DiffStatus:  "changed",
							IsViewed:    true,
							FileIcon:    mockIconForFile(`svg-mfi-file`),
						},
					},
				},
				{
					EntryMode:   "",
					DisplayName: "file1",
					FullName:    "file1",
					NameHash:    "60b27f004e454aca81b0480209cce5081ec52390",
					DiffStatus:  "added",
					FileIcon:    mockIconForFile(`svg-mfi-file`),
				},
			},
		},
	}, ret)
}
