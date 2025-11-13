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

package db

import (
	"context"

	"github.com/kumose/kmup/modules/setting"

	"xorm.io/builder"
)

// Iterate iterates all the Bean object
func Iterate[Bean any](ctx context.Context, cond builder.Cond, f func(ctx context.Context, bean *Bean) error) error {
	var start int
	batchSize := setting.Database.IterateBufferSize
	sess := GetEngine(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			beans := make([]*Bean, 0, batchSize)
			if cond != nil {
				sess = sess.Where(cond)
			}
			if err := sess.Limit(batchSize, start).Find(&beans); err != nil {
				return err
			}
			if len(beans) == 0 {
				return nil
			}
			start += len(beans)

			for _, bean := range beans {
				if err := f(ctx, bean); err != nil {
					return err
				}
			}
		}
	}
}
