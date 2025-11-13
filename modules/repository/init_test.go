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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeCustomLabels(t *testing.T) {
	files := mergeCustomLabelFiles(optionFileList{
		all:    []string{"a", "a.yaml", "a.yml"},
		custom: nil,
	})
	assert.Equal(t, []string{"a.yaml"}, files, "yaml file should win")

	files = mergeCustomLabelFiles(optionFileList{
		all:    []string{"a", "a.yaml"},
		custom: []string{"a"},
	})
	assert.Equal(t, []string{"a"}, files, "custom file should win")

	files = mergeCustomLabelFiles(optionFileList{
		all:    []string{"a", "a.yml", "a.yaml"},
		custom: []string{"a", "a.yml"},
	})
	assert.Equal(t, []string{"a.yml"}, files, "custom yml file should win if no yaml")
}
