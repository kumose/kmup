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
	"fmt"
	"testing"
	"time"

	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/timeutil"

	"github.com/stretchr/testify/assert"
)

func TestLabel_ToLabel(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	label := unittest.AssertExistsAndLoadBean(t, &issues_model.Label{ID: 1})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: label.RepoID})
	assert.Equal(t, &api.Label{
		ID:    label.ID,
		Name:  label.Name,
		Color: "abcdef",
		URL:   fmt.Sprintf("%sapi/v1/repos/user2/repo1/labels/%d", setting.AppURL, label.ID),
	}, ToLabel(label, repo, nil))
}

func TestMilestone_APIFormat(t *testing.T) {
	milestone := &issues_model.Milestone{
		ID:              3,
		RepoID:          4,
		Name:            "milestoneName",
		Content:         "milestoneContent",
		IsClosed:        false,
		NumOpenIssues:   5,
		NumClosedIssues: 6,
		CreatedUnix:     timeutil.TimeStamp(time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()),
		UpdatedUnix:     timeutil.TimeStamp(time.Date(1999, time.March, 1, 0, 0, 0, 0, time.UTC).Unix()),
		DeadlineUnix:    timeutil.TimeStamp(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()),
	}
	assert.Equal(t, api.Milestone{
		ID:           milestone.ID,
		State:        api.StateOpen,
		Title:        milestone.Name,
		Description:  milestone.Content,
		OpenIssues:   milestone.NumOpenIssues,
		ClosedIssues: milestone.NumClosedIssues,
		Created:      milestone.CreatedUnix.AsTime(),
		Updated:      milestone.UpdatedUnix.AsTimePtr(),
		Deadline:     milestone.DeadlineUnix.AsTimePtr(),
	}, *ToAPIMilestone(milestone))
}
