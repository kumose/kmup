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

package convert

import (
	"context"

	repo_model "github.com/kumose/kmup/models/repo"
	api "github.com/kumose/kmup/modules/structs"
)

// ToPushMirror convert from repo_model.PushMirror and remoteAddress to api.TopicResponse
func ToPushMirror(ctx context.Context, pm *repo_model.PushMirror) (*api.PushMirror, error) {
	repo := pm.GetRepository(ctx)
	return &api.PushMirror{
		RepoName:       repo.Name,
		RemoteName:     pm.RemoteName,
		RemoteAddress:  pm.RemoteAddress,
		CreatedUnix:    pm.CreatedUnix.AsTime(),
		LastUpdateUnix: pm.LastUpdateUnix.AsTimePtr(),
		LastError:      pm.LastError,
		Interval:       pm.Interval.String(),
		SyncOnCommit:   pm.SyncOnCommit,
	}, nil
}
