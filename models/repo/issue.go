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

package repo

import (
	"context"

	"github.com/kumose/kmup/models/unit"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
)

// ___________.__             ___________                     __
// \__    ___/|__| _____   ___\__    ___/___________    ____ |  | __ ___________
// |    |   |  |/     \_/ __ \|    |  \_  __ \__  \ _/ ___\|  |/ // __ \_  __ \
// |    |   |  |  Y Y  \  ___/|    |   |  | \// __ \\  \___|    <\  ___/|  | \/
// |____|   |__|__|_|  /\___  >____|   |__|  (____  /\___  >__|_ \\___  >__|
// \/     \/                    \/     \/     \/    \/

// CanEnableTimetracker returns true when the server admin enabled time tracking
// This overrules IsTimetrackerEnabled
func (repo *Repository) CanEnableTimetracker() bool {
	return setting.Service.EnableTimetracking
}

// IsTimetrackerEnabled returns whether or not the timetracker is enabled. It returns the default value from config if an error occurs.
func (repo *Repository) IsTimetrackerEnabled(ctx context.Context) bool {
	if !setting.Service.EnableTimetracking {
		return false
	}

	var u *RepoUnit
	var err error
	if u, err = repo.GetUnit(ctx, unit.TypeIssues); err != nil {
		return setting.Service.DefaultEnableTimetracking
	}
	return u.IssuesConfig().EnableTimetracker
}

// AllowOnlyContributorsToTrackTime returns value of IssuesConfig or the default value
func (repo *Repository) AllowOnlyContributorsToTrackTime(ctx context.Context) bool {
	var u *RepoUnit
	var err error
	if u, err = repo.GetUnit(ctx, unit.TypeIssues); err != nil {
		return setting.Service.DefaultAllowOnlyContributorsToTrackTime
	}
	return u.IssuesConfig().AllowOnlyContributorsToTrackTime
}

// IsDependenciesEnabled returns if dependencies are enabled and returns the default setting if not set.
func (repo *Repository) IsDependenciesEnabled(ctx context.Context) bool {
	var u *RepoUnit
	var err error
	if u, err = repo.GetUnit(ctx, unit.TypeIssues); err != nil {
		log.Trace("IsDependenciesEnabled: %v", err)
		return setting.Service.DefaultEnableDependencies
	}
	return u.IssuesConfig().EnableDependencies
}
