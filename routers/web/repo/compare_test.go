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
	"testing"

	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/services/gitdiff"

	"github.com/stretchr/testify/assert"
)

func TestAttachCommentsToLines(t *testing.T) {
	section := &gitdiff.DiffSection{
		Lines: []*gitdiff.DiffLine{
			{LeftIdx: 5, RightIdx: 10},
			{LeftIdx: 6, RightIdx: 11},
		},
	}

	lineComments := map[int64][]*issues_model.Comment{
		-5: {{ID: 100, CreatedUnix: 1000}},                               // left side comment
		10: {{ID: 200, CreatedUnix: 2000}},                               // right side comment
		11: {{ID: 300, CreatedUnix: 1500}, {ID: 301, CreatedUnix: 2500}}, // multiple comments
	}

	attachCommentsToLines(section, lineComments)

	// First line should have left and right comments
	assert.Len(t, section.Lines[0].Comments, 2)
	assert.Equal(t, int64(100), section.Lines[0].Comments[0].ID)
	assert.Equal(t, int64(200), section.Lines[0].Comments[1].ID)

	// Second line should have two comments, sorted by creation time
	assert.Len(t, section.Lines[1].Comments, 2)
	assert.Equal(t, int64(300), section.Lines[1].Comments[0].ID)
	assert.Equal(t, int64(301), section.Lines[1].Comments[1].ID)
}
