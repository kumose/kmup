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

package v1_11

import (
	"path/filepath"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"

	"xorm.io/builder"
	"xorm.io/xorm"
)

func RemoveAttachmentMissedRepo(x *xorm.Engine) error {
	type Attachment struct {
		UUID string `xorm:"uuid"`
	}
	var start int
	attachments := make([]*Attachment, 0, 50)
	for {
		err := x.Select("uuid").Where(builder.NotIn("release_id", builder.Select("id").From("`release`"))).
			And("release_id > 0").
			OrderBy("id").Limit(50, start).Find(&attachments)
		if err != nil {
			return err
		}

		for i := 0; i < len(attachments); i++ {
			uuid := attachments[i].UUID
			if err = util.RemoveAll(filepath.Join(setting.Attachment.Storage.Path, uuid[0:1], uuid[1:2], uuid)); err != nil {
				log.Warn("Unable to remove attachment file by UUID %s: %v", uuid, err)
			}
		}

		if len(attachments) < 50 {
			break
		}
		start += 50
		attachments = attachments[:0]
	}

	_, err := x.Exec("DELETE FROM attachment WHERE release_id > 0 AND release_id NOT IN (SELECT id FROM `release`)")
	return err
}
