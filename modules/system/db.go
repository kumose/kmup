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

package system

import (
	"context"

	"github.com/kumose/kmup/models/system"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/util"
)

// DBStore can be used to store app state items in local filesystem
type DBStore struct{}

// Get reads the state item
func (f *DBStore) Get(ctx context.Context, item StateItem) error {
	content, err := system.GetAppStateContent(ctx, item.Name())
	if err != nil {
		return err
	}
	if content == "" {
		return nil
	}
	return json.Unmarshal(util.UnsafeStringToBytes(content), item)
}

// Set saves the state item
func (f *DBStore) Set(ctx context.Context, item StateItem) error {
	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return system.SaveAppStateContent(ctx, item.Name(), util.UnsafeBytesToString(b))
}
