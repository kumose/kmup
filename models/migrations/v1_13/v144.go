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
	"github.com/kumose/kmup/modules/log"

	"xorm.io/builder"
	"xorm.io/xorm"
)

func UpdateMatrixWebhookHTTPMethod(x *xorm.Engine) error {
	matrixHookTaskType := 9 // value comes from the models package
	type Webhook struct {
		HTTPMethod string
	}

	cond := builder.Eq{"hook_task_type": matrixHookTaskType}.And(builder.Neq{"http_method": "PUT"})
	count, err := x.Where(cond).Cols("http_method").Update(&Webhook{HTTPMethod: "PUT"})
	if err == nil {
		log.Debug("Updated %d Matrix webhooks with http_method 'PUT'", count)
	}
	return err
}
