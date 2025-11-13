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

package private

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/modules/setting"
)

// RestoreParams structure holds a data for restore repository
type RestoreParams struct {
	RepoDir    string
	OwnerName  string
	RepoName   string
	Units      []string
	Validation bool
}

// RestoreRepo calls the internal RestoreRepo function
func RestoreRepo(ctx context.Context, repoDir, ownerName, repoName string, units []string, validation bool) ResponseExtra {
	reqURL := setting.LocalURL + "api/internal/restore_repo"

	req := newInternalRequestAPI(ctx, reqURL, "POST", RestoreParams{
		RepoDir:    repoDir,
		OwnerName:  ownerName,
		RepoName:   repoName,
		Units:      units,
		Validation: validation,
	})
	req.SetReadWriteTimeout(0) // since the request will spend much time, don't timeout
	return requestJSONClientMsg(req, fmt.Sprintf("Restore repo %s/%s successfully", ownerName, repoName))
}
