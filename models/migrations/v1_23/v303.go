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

package v1_23

import (
	"xorm.io/xorm"
)

// CommentMetaData stores metadata for a comment, these data will not be changed once inserted into database
type CommentMetaData struct {
	ProjectColumnID    int64  `json:"project_column_id"`
	ProjectColumnTitle string `json:"project_column_title"`
	ProjectTitle       string `json:"project_title"`
}

func AddCommentMetaDataColumn(x *xorm.Engine) error {
	type Comment struct {
		CommentMetaData *CommentMetaData `xorm:"JSON TEXT"` // put all non-index metadata in a single field
	}

	_, err := x.SyncWithOptions(xorm.SyncOptions{
		IgnoreConstrains: true,
		IgnoreIndices:    true,
	}, new(Comment))
	return err
}
