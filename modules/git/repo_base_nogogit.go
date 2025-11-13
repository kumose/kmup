// Copyright 2015 The Gogs Authors. All rights reserved.
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

//go:build !gogit

package git

import (
	"bufio"
	"context"
	"path/filepath"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/util"
)

const isGogit = false

// Repository represents a Git repository.
type Repository struct {
	Path string

	tagCache *ObjectCache[*Tag]

	gpgSettings *GPGSettings

	batchInUse bool
	batch      *Batch

	checkInUse bool
	check      *Batch

	Ctx             context.Context
	LastCommitCache *LastCommitCache

	objectFormat ObjectFormat
}

// OpenRepository opens the repository at the given path with the provided context.
func OpenRepository(ctx context.Context, repoPath string) (*Repository, error) {
	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}
	exist, err := util.IsDir(repoPath)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, util.NewNotExistErrorf("no such file or directory")
	}

	return &Repository{
		Path:     repoPath,
		tagCache: newObjectCache[*Tag](),
		Ctx:      ctx,
	}, nil
}

// CatFileBatch obtains a CatFileBatch for this repository
func (repo *Repository) CatFileBatch(ctx context.Context) (WriteCloserError, *bufio.Reader, func(), error) {
	if repo.batch == nil {
		var err error
		repo.batch, err = NewBatch(ctx, repo.Path)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if !repo.batchInUse {
		repo.batchInUse = true
		return repo.batch.Writer, repo.batch.Reader, func() {
			repo.batchInUse = false
		}, nil
	}

	log.Debug("Opening temporary cat file batch for: %s", repo.Path)
	tempBatch, err := NewBatch(ctx, repo.Path)
	if err != nil {
		return nil, nil, nil, err
	}
	return tempBatch.Writer, tempBatch.Reader, tempBatch.Close, nil
}

// CatFileBatchCheck obtains a CatFileBatchCheck for this repository
func (repo *Repository) CatFileBatchCheck(ctx context.Context) (WriteCloserError, *bufio.Reader, func(), error) {
	if repo.check == nil {
		var err error
		repo.check, err = NewBatchCheck(ctx, repo.Path)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if !repo.checkInUse {
		repo.checkInUse = true
		return repo.check.Writer, repo.check.Reader, func() {
			repo.checkInUse = false
		}, nil
	}

	log.Debug("Opening temporary cat file batch-check for: %s", repo.Path)
	tempBatchCheck, err := NewBatchCheck(ctx, repo.Path)
	if err != nil {
		return nil, nil, nil, err
	}
	return tempBatchCheck.Writer, tempBatchCheck.Reader, tempBatchCheck.Close, nil
}

func (repo *Repository) Close() error {
	if repo == nil {
		return nil
	}
	if repo.batch != nil {
		repo.batch.Close()
		repo.batch = nil
		repo.batchInUse = false
	}
	if repo.check != nil {
		repo.check.Close()
		repo.check = nil
		repo.checkInUse = false
	}
	repo.LastCommitCache = nil
	repo.tagCache = nil
	return nil
}
