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

	"github.com/kumose/kmup/modules/util"
	kmup_ctx "github.com/kumose/kmup/services/context"
)

type ContainedLinks struct { // TODO: better name?
	Branches      []*namedLink `json:"branches"`
	Tags          []*namedLink `json:"tags"`
	DefaultBranch string       `json:"default_branch"`
}

type namedLink struct { // TODO: better name?
	Name    string `json:"name"`
	WebLink string `json:"web_link"`
}

// LoadBranchesAndTags creates a new repository branch
func LoadBranchesAndTags(ctx context.Context, baseRepo *kmup_ctx.Repository, commitSHA string) (*ContainedLinks, error) {
	containedTags, err := baseRepo.GitRepo.ListOccurrences(ctx, "tag", commitSHA)
	if err != nil {
		return nil, fmt.Errorf("encountered a problem while querying %s: %w", "tags", err)
	}
	containedBranches, err := baseRepo.GitRepo.ListOccurrences(ctx, "branch", commitSHA)
	if err != nil {
		return nil, fmt.Errorf("encountered a problem while querying %s: %w", "branches", err)
	}

	result := &ContainedLinks{
		DefaultBranch: baseRepo.Repository.DefaultBranch,
		Branches:      make([]*namedLink, 0, len(containedBranches)),
		Tags:          make([]*namedLink, 0, len(containedTags)),
	}
	for _, tag := range containedTags {
		// TODO: Use a common method to get the link to a branch/tag instead of hard-coding it here
		result.Tags = append(result.Tags, &namedLink{
			Name:    tag,
			WebLink: fmt.Sprintf("%s/src/tag/%s", baseRepo.RepoLink, util.PathEscapeSegments(tag)),
		})
	}
	for _, branch := range containedBranches {
		result.Branches = append(result.Branches, &namedLink{
			Name:    branch,
			WebLink: fmt.Sprintf("%s/src/branch/%s", baseRepo.RepoLink, util.PathEscapeSegments(branch)),
		})
	}
	return result, nil
}
