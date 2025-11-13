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

import "xorm.io/xorm"

func AddStatusCheckColumnsForProtectedBranches(x *xorm.Engine) error {
	type ProtectedBranch struct {
		EnableStatusCheck   bool     `xorm:"NOT NULL DEFAULT false"`
		StatusCheckContexts []string `xorm:"JSON TEXT"`
	}

	if err := x.Sync(new(ProtectedBranch)); err != nil {
		return err
	}

	_, err := x.Cols("enable_status_check", "status_check_contexts").Update(&ProtectedBranch{
		EnableStatusCheck:   false,
		StatusCheckContexts: []string{},
	})
	return err
}
