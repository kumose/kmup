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

package issue

import (
	"strings"

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/base"
	"github.com/kumose/kmup/services/context"
)

// PrepareFilterIssueLabels reads the "labels" query parameter, sets `ctx.Data["Labels"]` and `ctx.Data["SelectLabels"]`
func PrepareFilterIssueLabels(ctx *context.Context, repoID int64, owner *user_model.User) (ret struct {
	AllLabels        []*issues_model.Label
	SelectedLabelIDs []int64
},
) {
	// 1,-2 means including label 1 and excluding label 2
	// 0 means issues with no label
	// blank means labels will not be filtered for issues
	selectLabels := ctx.FormString("labels")
	if selectLabels != "" {
		var err error
		ret.SelectedLabelIDs, err = base.StringsToInt64s(strings.Split(selectLabels, ","))
		if err != nil {
			ctx.Flash.Error(ctx.Tr("invalid_data", selectLabels), true)
		}
	}

	var allLabels []*issues_model.Label
	if repoID != 0 {
		repoLabels, err := issues_model.GetLabelsByRepoID(ctx, repoID, "", db.ListOptions{})
		if err != nil {
			ctx.ServerError("GetLabelsByRepoID", err)
			return ret
		}
		allLabels = append(allLabels, repoLabels...)
	}

	if owner != nil && owner.IsOrganization() {
		orgLabels, err := issues_model.GetLabelsByOrgID(ctx, owner.ID, "", db.ListOptions{})
		if err != nil {
			ctx.ServerError("GetLabelsByOrgID", err)
			return ret
		}
		allLabels = append(allLabels, orgLabels...)
	}

	// Get the exclusive scope for every label ID
	labelExclusiveScopes := make([]string, 0, len(ret.SelectedLabelIDs))
	for _, labelID := range ret.SelectedLabelIDs {
		foundExclusiveScope := false
		for _, label := range allLabels {
			if label.ID == labelID || label.ID == -labelID {
				labelExclusiveScopes = append(labelExclusiveScopes, label.ExclusiveScope())
				foundExclusiveScope = true
				break
			}
		}
		if !foundExclusiveScope {
			labelExclusiveScopes = append(labelExclusiveScopes, "")
		}
	}

	for _, l := range allLabels {
		l.LoadSelectedLabelsAfterClick(ret.SelectedLabelIDs, labelExclusiveScopes)
	}
	ctx.Data["Labels"] = allLabels
	ctx.Data["SelectLabels"] = selectLabels
	ret.AllLabels = allLabels
	return ret
}
