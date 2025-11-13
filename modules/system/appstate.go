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

import "context"

// StateStore is the interface to get/set app state items
type StateStore interface {
	Get(ctx context.Context, item StateItem) error
	Set(ctx context.Context, item StateItem) error
}

// StateItem provides the name for a state item. the name will be used to generate filenames, etc
type StateItem interface {
	Name() string
}

// AppState contains the state items for the app
var AppState StateStore

// Init initialize AppState interface
func Init() error {
	AppState = &DBStore{}
	return nil
}
