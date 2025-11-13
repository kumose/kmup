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

package repository

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
)

// CreateTemporaryPath creates a temporary path
func CreateTemporaryPath(prefix string) (string, context.CancelFunc, error) {
	basePath, cleanup, err := setting.AppDataTempDir("local-repo").MkdirTempRandom(prefix + ".git")
	if err != nil {
		log.Error("Unable to create temporary directory: %s-*.git (%v)", prefix, err)
		return "", nil, fmt.Errorf("failed to create dir %s-*.git: %w", prefix, err)
	}
	return basePath, cleanup, nil
}
