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

package v1_14

import (
	"xorm.io/xorm"
)

func FixRepoTopics(x *xorm.Engine) error {
	type Repository struct {
		ID     int64    `xorm:"pk autoincr"`
		Topics []string `xorm:"TEXT JSON"`
	}

	const batchSize = 100
	sess := x.NewSession()
	defer sess.Close()
	repos := make([]*Repository, 0, batchSize)
	topics := make([]string, 0, batchSize)
	for start := 0; ; start += batchSize {
		repos = repos[:0]

		if err := sess.Begin(); err != nil {
			return err
		}

		if err := sess.Limit(batchSize, start).Find(&repos); err != nil {
			return err
		}

		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			topics = topics[:0]
			if err := sess.Select("name").Table("topic").
				Join("INNER", "repo_topic", "repo_topic.topic_id = topic.id").
				Where("repo_topic.repo_id = ?", repo.ID).Desc("topic.repo_count").Find(&topics); err != nil {
				return err
			}
			repo.Topics = topics
			if _, err := sess.ID(repo.ID).Cols("topics").Update(repo); err != nil {
				return err
			}
		}

		if err := sess.Commit(); err != nil {
			return err
		}
	}

	return nil
}
