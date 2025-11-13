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

package util

import "runtime"

func CallerFuncName(optSkipParent ...int) string {
	pc := make([]uintptr, 1)
	skipParent := 0
	if len(optSkipParent) > 0 {
		skipParent = optSkipParent[0]
	}
	runtime.Callers(skipParent+1 /*this*/ +1 /*runtime*/, pc)
	funcName := runtime.FuncForPC(pc[0]).Name()
	return funcName
}
