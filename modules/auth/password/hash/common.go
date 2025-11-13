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

package hash

import (
	"strconv"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/util"
)

func parseIntParam(value, param, algorithmName, config string, previousErr error) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Error("invalid integer for %s representation in %s hash spec %s", param, algorithmName, config)
		return 0, err
	}
	return parsed, previousErr // <- Keep the previous error as this function should still return an error once everything has been checked if any call failed
}

func parseUintParam[T uint32 | uint8](value, param, algorithmName, config string, previousErr error) (ret T, _ error) {
	_, isUint32 := any(ret).(uint32)
	parsed, err := strconv.ParseUint(value, 10, util.Iif(isUint32, 32, 8))
	if err != nil {
		log.Error("invalid integer for %s representation in %s hash spec %s", param, algorithmName, config)
		return 0, err
	}
	return T(parsed), previousErr // <- Keep the previous error as this function should still return an error once everything has been checked if any call failed
}
