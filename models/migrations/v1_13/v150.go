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
	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddPrimaryKeyToRepoTopic(x *xorm.Engine) error {
	// Topic represents a topic of repositories
	type Topic struct {
		ID          int64  `xorm:"pk autoincr"`
		Name        string `xorm:"UNIQUE VARCHAR(25)"`
		RepoCount   int
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	// RepoTopic represents associated repositories and topics
	type RepoTopic struct {
		RepoID  int64 `xorm:"pk"`
		TopicID int64 `xorm:"pk"`
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	base.RecreateTable(sess, &Topic{})
	base.RecreateTable(sess, &RepoTopic{})

	return sess.Commit()
}
