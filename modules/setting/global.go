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

package setting

// Global settings
var (
	// RunUser is the OS user that Kmup is running as. ini:"RUN_USER"
	RunUser string
	// RunMode is the running mode of Kmup, it only accepts two values: "dev" and "prod".
	// Non-dev values will be replaced by "prod". ini: "RUN_MODE"
	RunMode string
	// IsProd is true if RunMode is not "dev"
	IsProd bool

	// AppName is the Application name, used in the page title. ini: "APP_NAME"
	AppName string
)
