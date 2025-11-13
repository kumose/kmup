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

package v1_13

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddProjectsInfo(x *xorm.Engine) error {
	// Create new tables
	type (
		ProjectType      uint8
		ProjectBoardType uint8
	)

	type Project struct {
		ID          int64  `xorm:"pk autoincr"`
		Title       string `xorm:"INDEX NOT NULL"`
		Description string `xorm:"TEXT"`
		RepoID      int64  `xorm:"INDEX"`
		CreatorID   int64  `xorm:"NOT NULL"`
		IsClosed    bool   `xorm:"INDEX"`

		BoardType ProjectBoardType
		Type      ProjectType

		ClosedDateUnix timeutil.TimeStamp
		CreatedUnix    timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix    timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	if err := x.Sync(new(Project)); err != nil {
		return err
	}

	type Comment struct {
		OldProjectID int64
		ProjectID    int64
	}

	if err := x.Sync(new(Comment)); err != nil {
		return err
	}

	type Repository struct {
		ID                int64
		NumProjects       int `xorm:"NOT NULL DEFAULT 0"`
		NumClosedProjects int `xorm:"NOT NULL DEFAULT 0"`
	}

	if err := x.Sync(new(Repository)); err != nil {
		return err
	}

	// ProjectIssue saves relation from issue to a project
	type ProjectIssue struct {
		ID             int64 `xorm:"pk autoincr"`
		IssueID        int64 `xorm:"INDEX"`
		ProjectID      int64 `xorm:"INDEX"`
		ProjectBoardID int64 `xorm:"INDEX"`
	}

	if err := x.Sync(new(ProjectIssue)); err != nil {
		return err
	}

	type ProjectBoard struct {
		ID      int64 `xorm:"pk autoincr"`
		Title   string
		Default bool `xorm:"NOT NULL DEFAULT false"`

		ProjectID int64 `xorm:"INDEX NOT NULL"`
		CreatorID int64 `xorm:"NOT NULL"`

		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	return x.Sync(new(ProjectBoard))
}
