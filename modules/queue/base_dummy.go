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

package queue

import "context"

type baseDummy struct{}

var _ baseQueue = (*baseDummy)(nil)

func newBaseDummy(cfg *BaseConfig, unique bool) (baseQueue, error) {
	return &baseDummy{}, nil
}

func (q *baseDummy) PushItem(ctx context.Context, data []byte) error {
	return nil
}

func (q *baseDummy) PopItem(ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (q *baseDummy) Len(ctx context.Context) (int, error) {
	return 0, nil
}

func (q *baseDummy) HasItem(ctx context.Context, data []byte) (bool, error) {
	return false, nil
}

func (q *baseDummy) Close() error {
	return nil
}

func (q *baseDummy) RemoveAll(ctx context.Context) error {
	return nil
}
