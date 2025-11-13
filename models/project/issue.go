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

package project

import (
	"context"
	"errors"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/util"
)

// ProjectIssue saves relation from issue to a project
type ProjectIssue struct { //revive:disable-line:exported
	ID        int64 `xorm:"pk autoincr"`
	IssueID   int64 `xorm:"INDEX"`
	ProjectID int64 `xorm:"INDEX"`

	// ProjectColumnID should not be zero since 1.22. If it's zero, the issue will not be displayed on UI and it might result in errors.
	ProjectColumnID int64 `xorm:"'project_board_id' INDEX"`

	// the sorting order on the column
	Sorting int64 `xorm:"NOT NULL DEFAULT 0"`
}

func init() {
	db.RegisterModel(new(ProjectIssue))
}

func deleteProjectIssuesByProjectID(ctx context.Context, projectID int64) error {
	_, err := db.GetEngine(ctx).Where("project_id=?", projectID).Delete(&ProjectIssue{})
	return err
}

func (c *Column) moveIssuesToAnotherColumn(ctx context.Context, newColumn *Column) error {
	if c.ProjectID != newColumn.ProjectID {
		return errors.New("columns have to be in the same project")
	}

	if c.ID == newColumn.ID {
		return nil
	}

	res := struct {
		MaxSorting int64
		IssueCount int64
	}{}
	if _, err := db.GetEngine(ctx).Select("max(sorting) as max_sorting, count(*) as issue_count").
		Table("project_issue").
		Where("project_id=?", newColumn.ProjectID).
		And("project_board_id=?", newColumn.ID).
		Get(&res); err != nil {
		return err
	}

	issues, err := c.GetIssues(ctx)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		return nil
	}

	nextSorting := util.Iif(res.IssueCount > 0, res.MaxSorting+1, 0)
	return db.WithTx(ctx, func(ctx context.Context) error {
		for i, issue := range issues {
			issue.ProjectColumnID = newColumn.ID
			issue.Sorting = nextSorting + int64(i)
			if _, err := db.GetEngine(ctx).ID(issue.ID).Cols("project_board_id", "sorting").Update(issue); err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteAllProjectIssueByIssueIDsAndProjectIDs delete all project's issues by issue's and project's ids
func DeleteAllProjectIssueByIssueIDsAndProjectIDs(ctx context.Context, issueIDs, projectIDs []int64) error {
	_, err := db.GetEngine(ctx).In("project_id", projectIDs).In("issue_id", issueIDs).Delete(&ProjectIssue{})
	return err
}
