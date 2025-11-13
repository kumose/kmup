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

// RuntimeState contains app state for runtime, and we can save remote version for update checker here in future
type RuntimeState struct {
	LastAppPath    string `json:"last_app_path"`
	LastCustomConf string `json:"last_custom_conf"`
}

// Name returns the item name
func (a RuntimeState) Name() string {
	return "runtime-state"
}
