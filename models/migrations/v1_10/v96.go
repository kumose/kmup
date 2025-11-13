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

package v1_10

import (
	"path/filepath"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"

	"xorm.io/xorm"
)

func DeleteOrphanedAttachments(x *xorm.Engine) error {
	type Attachment struct {
		ID        int64  `xorm:"pk autoincr"`
		UUID      string `xorm:"uuid UNIQUE"`
		IssueID   int64  `xorm:"INDEX"`
		ReleaseID int64  `xorm:"INDEX"`
		CommentID int64
	}

	sess := x.NewSession()
	defer sess.Close()

	limit := setting.Database.IterateBufferSize
	if limit <= 0 {
		limit = 50
	}

	for {
		attachments := make([]Attachment, 0, limit)
		if err := sess.Where("`issue_id` = 0 and (`release_id` = 0 or `release_id` not in (select `id` from `release`))").
			Cols("id, uuid").Limit(limit).
			Asc("id").
			Find(&attachments); err != nil {
			return err
		}
		if len(attachments) == 0 {
			return nil
		}

		ids := make([]int64, 0, limit)
		for _, attachment := range attachments {
			ids = append(ids, attachment.ID)
		}
		if len(ids) > 0 {
			if _, err := sess.In("id", ids).Delete(new(Attachment)); err != nil {
				return err
			}
		}

		for _, attachment := range attachments {
			uuid := attachment.UUID
			if err := util.RemoveAll(filepath.Join(setting.Attachment.Storage.Path, uuid[0:1], uuid[1:2], uuid)); err != nil {
				return err
			}
		}
		if len(attachments) < limit {
			return nil
		}
	}
}
