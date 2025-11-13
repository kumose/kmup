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

package v1_9

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddUploaderIDForAttachment(x *xorm.Engine) error {
	type Attachment struct {
		ID            int64  `xorm:"pk autoincr"`
		UUID          string `xorm:"uuid UNIQUE"`
		IssueID       int64  `xorm:"INDEX"`
		ReleaseID     int64  `xorm:"INDEX"`
		UploaderID    int64  `xorm:"INDEX DEFAULT 0"`
		CommentID     int64
		Name          string
		DownloadCount int64              `xorm:"DEFAULT 0"`
		Size          int64              `xorm:"DEFAULT 0"`
		CreatedUnix   timeutil.TimeStamp `xorm:"created"`
	}

	return x.Sync(new(Attachment))
}
