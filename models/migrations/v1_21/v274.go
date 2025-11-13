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

package v1_21

import (
	"time"

	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddExpiredUnixColumnInActionArtifactTable(x *xorm.Engine) error {
	type ActionArtifact struct {
		ExpiredUnix timeutil.TimeStamp `xorm:"index"` // time when the artifact will be expired
	}
	if err := x.Sync(new(ActionArtifact)); err != nil {
		return err
	}
	return updateArtifactsExpiredUnixTo90Days(x)
}

func updateArtifactsExpiredUnixTo90Days(x *xorm.Engine) error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}
	expiredTime := time.Now().AddDate(0, 0, 90).Unix()
	if _, err := sess.Exec(`UPDATE action_artifact SET expired_unix=? WHERE status='2' AND expired_unix is NULL`, expiredTime); err != nil {
		return err
	}

	return sess.Commit()
}
