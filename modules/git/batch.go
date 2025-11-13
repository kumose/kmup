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

package git

import (
	"bufio"
	"context"
)

type Batch struct {
	cancel context.CancelFunc
	Reader *bufio.Reader
	Writer WriteCloserError
}

// NewBatch creates a new batch for the given repository, the Close must be invoked before release the batch
func NewBatch(ctx context.Context, repoPath string) (*Batch, error) {
	// Now because of some insanity with git cat-file not immediately failing if not run in a valid git directory we need to run git rev-parse first!
	if err := ensureValidGitRepository(ctx, repoPath); err != nil {
		return nil, err
	}

	var batch Batch
	batch.Writer, batch.Reader, batch.cancel = catFileBatch(ctx, repoPath)
	return &batch, nil
}

func NewBatchCheck(ctx context.Context, repoPath string) (*Batch, error) {
	// Now because of some insanity with git cat-file not immediately failing if not run in a valid git directory we need to run git rev-parse first!
	if err := ensureValidGitRepository(ctx, repoPath); err != nil {
		return nil, err
	}

	var check Batch
	check.Writer, check.Reader, check.cancel = catFileBatchCheck(ctx, repoPath)
	return &check, nil
}

func (b *Batch) Close() {
	if b.cancel != nil {
		b.cancel()
		b.Reader = nil
		b.Writer = nil
		b.cancel = nil
	}
}
