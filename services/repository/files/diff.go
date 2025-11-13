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

package files

import (
	"context"
	"strings"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/services/gitdiff"
)

// GetDiffPreview produces and returns diff result of a file which is not yet committed.
func GetDiffPreview(ctx context.Context, repo *repo_model.Repository, branch, treePath, content string) (*gitdiff.Diff, error) {
	if branch == "" {
		branch = repo.DefaultBranch
	}
	t, err := NewTemporaryUploadRepository(repo)
	if err != nil {
		return nil, err
	}
	defer t.Close()
	if err := t.Clone(ctx, branch, true); err != nil {
		return nil, err
	}
	if err := t.SetDefaultIndex(ctx); err != nil {
		return nil, err
	}

	// Add the object to the database
	objectHash, err := t.HashObjectAndWrite(ctx, strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	// Add the object to the index
	if err := t.AddObjectToIndex(ctx, "100644", objectHash, treePath); err != nil {
		return nil, err
	}
	return t.DiffIndex(ctx)
}
