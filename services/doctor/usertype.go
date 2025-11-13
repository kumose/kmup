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

package doctor

import (
	"context"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
)

func checkUserType(ctx context.Context, logger log.Logger, autofix bool) error {
	count, err := user_model.CountWrongUserType(ctx)
	if err != nil {
		logger.Critical("Error: %v whilst counting wrong user types")
		return err
	}
	if count > 0 {
		if autofix {
			if count, err = user_model.FixWrongUserType(ctx); err != nil {
				logger.Critical("Error: %v whilst fixing wrong user types")
				return err
			}
			logger.Info("%d users with wrong type fixed", count)
		} else {
			logger.Warn("%d users with wrong type exist", count)
		}
	}
	return nil
}

func init() {
	Register(&Check{
		Title:     "Check if user with wrong type exist",
		Name:      "check-user-type",
		IsDefault: true,
		Run:       checkUserType,
		Priority:  3,
	})
}
