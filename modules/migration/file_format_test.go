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

package migration

import (
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
)

func TestMigrationJSON_IssueOK(t *testing.T) {
	issues := make([]*Issue, 0, 10)
	err := Load("file_format_testdata/issue_a.json", &issues, true)
	assert.NoError(t, err)
	err = Load("file_format_testdata/issue_a.yml", &issues, true)
	assert.NoError(t, err)
}

func TestMigrationJSON_IssueFail(t *testing.T) {
	issues := make([]*Issue, 0, 10)
	err := Load("file_format_testdata/issue_b.json", &issues, true)
	if _, ok := err.(*jsonschema.ValidationError); ok {
		errors := strings.Split(err.(*jsonschema.ValidationError).GoString(), "\n")
		assert.Contains(t, errors[1], "missing properties")
		assert.Contains(t, errors[1], "poster_id")
	} else {
		t.Fatalf("got: type %T with value %s, want: *jsonschema.ValidationError", err, err)
	}
}

func TestMigrationJSON_MilestoneOK(t *testing.T) {
	milestones := make([]*Milestone, 0, 10)
	err := Load("file_format_testdata/milestones.json", &milestones, true)
	assert.NoError(t, err)
}
