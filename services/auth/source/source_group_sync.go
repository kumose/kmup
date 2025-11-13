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

package source

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/models/organization"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/log"
	org_service "github.com/kumose/kmup/services/org"
)

type syncType int

const (
	syncAdd syncType = iota
	syncRemove
)

// SyncGroupsToTeams maps authentication source groups to organization and team memberships
func SyncGroupsToTeams(ctx context.Context, user *user_model.User, sourceUserGroups container.Set[string], sourceGroupTeamMapping map[string]map[string][]string, performRemoval bool) error {
	orgCache := make(map[string]*organization.Organization)
	teamCache := make(map[string]*organization.Team)
	return SyncGroupsToTeamsCached(ctx, user, sourceUserGroups, sourceGroupTeamMapping, performRemoval, orgCache, teamCache)
}

// SyncGroupsToTeamsCached maps authentication source groups to organization and team memberships
func SyncGroupsToTeamsCached(ctx context.Context, user *user_model.User, sourceUserGroups container.Set[string], sourceGroupTeamMapping map[string]map[string][]string, performRemoval bool, orgCache map[string]*organization.Organization, teamCache map[string]*organization.Team) error {
	membershipsToAdd, membershipsToRemove := resolveMappedMemberships(sourceUserGroups, sourceGroupTeamMapping)

	if performRemoval {
		if err := syncGroupsToTeamsCached(ctx, user, membershipsToRemove, syncRemove, orgCache, teamCache); err != nil {
			return fmt.Errorf("could not sync[remove] user groups: %w", err)
		}
	}

	if err := syncGroupsToTeamsCached(ctx, user, membershipsToAdd, syncAdd, orgCache, teamCache); err != nil {
		return fmt.Errorf("could not sync[add] user groups: %w", err)
	}

	return nil
}

func resolveMappedMemberships(sourceUserGroups container.Set[string], sourceGroupTeamMapping map[string]map[string][]string) (map[string][]string, map[string][]string) {
	membershipsToAdd := map[string][]string{}
	membershipsToRemove := map[string][]string{}
	for group, memberships := range sourceGroupTeamMapping {
		isUserInGroup := sourceUserGroups.Contains(group)
		if isUserInGroup {
			for org, teams := range memberships {
				membershipsToAdd[org] = append(membershipsToAdd[org], teams...)
			}
		} else {
			for org, teams := range memberships {
				membershipsToRemove[org] = append(membershipsToRemove[org], teams...)
			}
		}
	}
	return membershipsToAdd, membershipsToRemove
}

func syncGroupsToTeamsCached(ctx context.Context, user *user_model.User, orgTeamMap map[string][]string, action syncType, orgCache map[string]*organization.Organization, teamCache map[string]*organization.Team) error {
	for orgName, teamNames := range orgTeamMap {
		var err error
		org, ok := orgCache[orgName]
		if !ok {
			org, err = organization.GetOrgByName(ctx, orgName)
			if err != nil {
				if organization.IsErrOrgNotExist(err) {
					// organization must be created before group sync
					log.Warn("group sync: Could not find organisation %s: %v", orgName, err)
					continue
				}
				return err
			}
			orgCache[orgName] = org
		}
		for _, teamName := range teamNames {
			team, ok := teamCache[orgName+teamName]
			if !ok {
				team, err = org.GetTeam(ctx, teamName)
				if err != nil {
					if organization.IsErrTeamNotExist(err) {
						// team must be created before group sync
						log.Warn("group sync: Could not find team %s: %v", teamName, err)
						continue
					}
					return err
				}
				teamCache[orgName+teamName] = team
			}

			isMember, err := organization.IsTeamMember(ctx, org.ID, team.ID, user.ID)
			if err != nil {
				return err
			}

			if action == syncAdd && !isMember {
				if err := org_service.AddTeamMember(ctx, team, user); err != nil {
					log.Error("group sync: Could not add user to team: %v", err)
					return err
				}
			} else if action == syncRemove && isMember {
				if err := org_service.RemoveTeamMember(ctx, team, user); err != nil {
					log.Error("group sync: Could not remove user from team: %v", err)
					return err
				}
			}
		}
	}
	return nil
}
