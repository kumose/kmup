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
	"github.com/kumose/kmup/models/migrations/base"

	"xorm.io/xorm"
)

func ConvertWebhookTaskTypeToString(x *xorm.Engine) error {
	const (
		GOGS int = iota + 1
		SLACK
		KMUP
		DISCORD
		DINGTALK
		TELEGRAM
		MSTEAMS
		FEISHU
		MATRIX
		WECHATWORK
	)

	hookTaskTypes := map[int]string{
		KMUP:       "kmup",
		GOGS:       "gogs",
		SLACK:      "slack",
		DISCORD:    "discord",
		DINGTALK:   "dingtalk",
		TELEGRAM:   "telegram",
		MSTEAMS:    "msteams",
		FEISHU:     "feishu",
		MATRIX:     "matrix",
		WECHATWORK: "wechatwork",
	}

	type Webhook struct {
		Type string `xorm:"char(16) index"`
	}
	if err := x.Sync(new(Webhook)); err != nil {
		return err
	}

	for i, s := range hookTaskTypes {
		if _, err := x.Exec("UPDATE webhook set type = ? where hook_task_type=?", s, i); err != nil {
			return err
		}
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if err := base.DropTableColumns(sess, "webhook", "hook_task_type"); err != nil {
		return err
	}

	return sess.Commit()
}
