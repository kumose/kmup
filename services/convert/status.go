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
	"net/url"

	git_model "github.com/kumose/kmup/models/git"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
)

// ToCommitStatus converts git_model.CommitStatus to api.CommitStatus
func ToCommitStatus(ctx context.Context, status *git_model.CommitStatus) *api.CommitStatus {
	apiStatus := &api.CommitStatus{
		Created:     status.CreatedUnix.AsTime(),
		Updated:     status.CreatedUnix.AsTime(),
		State:       status.State,
		TargetURL:   status.TargetURL,
		Description: status.Description,
		ID:          status.Index,
		URL:         status.APIURL(ctx),
		Context:     status.Context,
	}

	if status.CreatorID != 0 {
		creator, _ := user_model.GetUserByID(ctx, status.CreatorID)
		apiStatus.Creator = ToUser(ctx, creator, nil)
	}

	return apiStatus
}

func ToCommitStatuses(ctx context.Context, statuses []*git_model.CommitStatus) []*api.CommitStatus {
	apiStatuses := make([]*api.CommitStatus, len(statuses))
	for i, status := range statuses {
		apiStatuses[i] = ToCommitStatus(ctx, status)
	}
	return apiStatuses
}

// ToCombinedStatus converts List of CommitStatus to a CombinedStatus
func ToCombinedStatus(ctx context.Context, commitID string, statuses []*git_model.CommitStatus, repo *api.Repository) *api.CombinedStatus {
	status := api.CombinedStatus{
		SHA:        commitID,
		TotalCount: len(statuses),
		Repository: repo,
		CommitURL:  repo.URL + "/commits/" + url.PathEscape(commitID),
		URL:        repo.URL + "/commits/" + url.PathEscape(commitID) + "/status",
	}

	combinedStatus := git_model.CalcCommitStatus(statuses)
	if combinedStatus != nil {
		status.Statuses = ToCommitStatuses(ctx, statuses)
		status.State = combinedStatus.State
	}
	return &status
}
