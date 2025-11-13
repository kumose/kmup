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
	"github.com/kumose/kmup/models/migrations/base"

	"xorm.io/xorm"
)

func DropWebhookColumns(x *xorm.Engine) error {
	// Make sure the columns exist before dropping them
	type Webhook struct {
		Signature string `xorm:"TEXT"`
		IsSSL     bool   `xorm:"is_ssl"`
	}
	if err := x.Sync(new(Webhook)); err != nil {
		return err
	}

	type HookTask struct {
		Typ         string `xorm:"VARCHAR(16) index"`
		URL         string `xorm:"TEXT"`
		Signature   string `xorm:"TEXT"`
		HTTPMethod  string `xorm:"http_method"`
		ContentType int
		IsSSL       bool
	}
	if err := x.Sync(new(HookTask)); err != nil {
		return err
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if err := base.DropTableColumns(sess, "webhook", "signature", "is_ssl"); err != nil {
		return err
	}
	if err := base.DropTableColumns(sess, "hook_task", "typ", "url", "signature", "http_method", "content_type", "is_ssl"); err != nil {
		return err
	}

	return sess.Commit()
}
