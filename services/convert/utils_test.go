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

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCorrectPageSize(t *testing.T) {
	assert.Equal(t, 30, ToCorrectPageSize(0))
	assert.Equal(t, 30, ToCorrectPageSize(-10))
	assert.Equal(t, 20, ToCorrectPageSize(20))
	assert.Equal(t, 50, ToCorrectPageSize(100))
}

func TestToGitServiceType(t *testing.T) {
	tc := []struct {
		typ  string
		enum int
	}{{
		typ: "trash", enum: 1,
	}, {
		typ: "github", enum: 2,
	}, {
		typ: "kmup", enum: 3,
	}, {
		typ: "gitlab", enum: 4,
	}, {
		typ: "gogs", enum: 5,
	}, {
		typ: "onedev", enum: 6,
	}, {
		typ: "gitbucket", enum: 7,
	}, {
		typ: "codebase", enum: 8,
	}, {
		typ: "codecommit", enum: 9,
	}}
	for _, test := range tc {
		assert.EqualValues(t, test.enum, ToGitServiceType(test.typ))
	}
}
