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
	"crypto/sha1"
	"fmt"

	"xorm.io/xorm"
)

func hashContext(context string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(context)))
}

func AddCommitStatusContext(x *xorm.Engine) error {
	type CommitStatus struct {
		ID          int64  `xorm:"pk autoincr"`
		ContextHash string `xorm:"char(40) index"`
		Context     string `xorm:"TEXT"`
	}

	if err := x.Sync(new(CommitStatus)); err != nil {
		return err
	}

	sess := x.NewSession()
	defer sess.Close()

	start := 0
	for {
		statuses := make([]*CommitStatus, 0, 100)
		err := sess.OrderBy("id").Limit(100, start).Find(&statuses)
		if err != nil {
			return err
		}
		if len(statuses) == 0 {
			break
		}

		if err = sess.Begin(); err != nil {
			return err
		}

		for _, status := range statuses {
			status.ContextHash = hashContext(status.Context)
			if _, err := sess.ID(status.ID).Cols("context_hash").Update(status); err != nil {
				return err
			}
		}

		if err := sess.Commit(); err != nil {
			return err
		}

		if len(statuses) < 100 {
			break
		}

		start += len(statuses)
	}

	return nil
}
