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

package v1_15

import (
	"fmt"
	"time"

	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func CreatePushMirrorTable(x *xorm.Engine) error {
	type PushMirror struct {
		ID         int64 `xorm:"pk autoincr"`
		RepoID     int64 `xorm:"INDEX"`
		RemoteName string

		Interval       time.Duration
		CreatedUnix    timeutil.TimeStamp `xorm:"created"`
		LastUpdateUnix timeutil.TimeStamp `xorm:"INDEX last_update"`
		LastError      string             `xorm:"text"`
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := sess.Sync(new(PushMirror)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}

	return sess.Commit()
}
