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

	activities_model "github.com/kumose/kmup/models/activities"
	perm_model "github.com/kumose/kmup/models/perm"
	access_model "github.com/kumose/kmup/models/perm/access"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	api "github.com/kumose/kmup/modules/structs"
)

func ToActivity(ctx context.Context, ac *activities_model.Action, doer *user_model.User) *api.Activity {
	p, err := access_model.GetUserRepoPermission(ctx, ac.Repo, doer)
	if err != nil {
		log.Error("GetUserRepoPermission[%d]: %v", ac.RepoID, err)
		p.AccessMode = perm_model.AccessModeNone
	}

	result := &api.Activity{
		ID:        ac.ID,
		UserID:    ac.UserID,
		OpType:    ac.OpType.String(),
		ActUserID: ac.ActUserID,
		ActUser:   ToUser(ctx, ac.ActUser, doer),
		RepoID:    ac.RepoID,
		Repo:      ToRepo(ctx, ac.Repo, p),
		RefName:   ac.RefName,
		IsPrivate: ac.IsPrivate,
		Content:   ac.Content,
		Created:   ac.CreatedUnix.AsTime(),
	}

	if ac.Comment != nil {
		result.CommentID = ac.CommentID
		result.Comment = ToAPIComment(ctx, ac.Repo, ac.Comment)
	}

	return result
}

func ToActivities(ctx context.Context, al activities_model.ActionList, doer *user_model.User) []*api.Activity {
	result := make([]*api.Activity, 0, len(al))
	for _, ac := range al {
		result = append(result, ToActivity(ctx, ac, doer))
	}
	return result
}
