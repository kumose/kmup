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

package doctor

import (
	"context"

	"github.com/kumose/kmup/models/db"
	org_model "github.com/kumose/kmup/models/organization"
	"github.com/kumose/kmup/models/perm"
	"github.com/kumose/kmup/modules/log"
	org_service "github.com/kumose/kmup/services/org"

	"xorm.io/builder"
)

func fixOwnerTeamCreateOrgRepo(ctx context.Context, logger log.Logger, autofix bool) error {
	count := 0

	err := db.Iterate(
		ctx,
		builder.Eq{"authorize": perm.AccessModeOwner, "can_create_org_repo": false},
		func(ctx context.Context, team *org_model.Team) error {
			team.CanCreateOrgRepo = true
			count++

			if !autofix {
				return nil
			}

			return org_service.UpdateTeam(ctx, team, false, false)
		},
	)
	if err != nil {
		logger.Critical("Unable to iterate across repounits to fix incorrect can_create_org_repo: Error %v", err)
		return err
	}

	if !autofix {
		if count == 0 {
			logger.Info("Found no team with incorrect can_create_org_repo")
		} else {
			logger.Warn("Found %d teams with incorrect can_create_org_repo", count)
		}
		return nil
	}
	logger.Info("Fixed %d teams with incorrect can_create_org_repo", count)

	return nil
}

func init() {
	Register(&Check{
		Title:     "Check for incorrect can_create_org_repo for org owner teams",
		Name:      "fix-owner-team-create-org-repo",
		IsDefault: false,
		Run:       fixOwnerTeamCreateOrgRepo,
		Priority:  7,
	})
}
