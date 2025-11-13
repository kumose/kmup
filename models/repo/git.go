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

	"github.com/kumose/kmup/models/db"
)

// MergeStyle represents the approach to merge commits into base branch.
type MergeStyle string

const (
	// MergeStyleMerge create merge commit
	MergeStyleMerge MergeStyle = "merge"
	// MergeStyleRebase rebase before merging, and fast-forward
	MergeStyleRebase MergeStyle = "rebase"
	// MergeStyleRebaseMerge rebase before merging with merge commit (--no-ff)
	MergeStyleRebaseMerge MergeStyle = "rebase-merge"
	// MergeStyleSquash squash commits into single commit before merging
	MergeStyleSquash MergeStyle = "squash"
	// MergeStyleFastForwardOnly fast-forward merge if possible, otherwise fail
	MergeStyleFastForwardOnly MergeStyle = "fast-forward-only"
	// MergeStyleManuallyMerged pr has been merged manually, just mark it as merged directly
	MergeStyleManuallyMerged MergeStyle = "manually-merged"
	// MergeStyleRebaseUpdate not a merge style, used to update pull head by rebase
	MergeStyleRebaseUpdate MergeStyle = "rebase-update-only"
)

// UpdateDefaultBranch updates the default branch
func UpdateDefaultBranch(ctx context.Context, repo *Repository) error {
	_, err := db.GetEngine(ctx).ID(repo.ID).Cols("default_branch").Update(repo)
	return err
}
