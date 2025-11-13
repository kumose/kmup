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

package internal

import (
	"context"
	"errors"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/indexer"
	"github.com/kumose/kmup/modules/indexer/internal"
)

// Indexer defines an interface to index and search code contents
type Indexer interface {
	internal.Indexer
	Index(ctx context.Context, repo *repo_model.Repository, sha string, changes *RepoChanges) error
	Delete(ctx context.Context, repoID int64) error
	Search(ctx context.Context, opts *SearchOptions) (int64, []*SearchResult, []*SearchResultLanguages, error)
	SupportedSearchModes() []indexer.SearchMode
}

type SearchOptions struct {
	RepoIDs  []int64
	Keyword  string
	Language string

	SearchMode indexer.SearchModeType

	db.Paginator
}

// NewDummyIndexer returns a dummy indexer
func NewDummyIndexer() Indexer {
	return &dummyIndexer{
		Indexer: internal.NewDummyIndexer(),
	}
}

type dummyIndexer struct {
	internal.Indexer
}

func (d *dummyIndexer) SupportedSearchModes() []indexer.SearchMode {
	return nil
}

func (d *dummyIndexer) Index(ctx context.Context, repo *repo_model.Repository, sha string, changes *RepoChanges) error {
	return errors.New("indexer is not ready")
}

func (d *dummyIndexer) Delete(ctx context.Context, repoID int64) error {
	return errors.New("indexer is not ready")
}

func (d *dummyIndexer) Search(ctx context.Context, opts *SearchOptions) (int64, []*SearchResult, []*SearchResultLanguages, error) {
	return 0, nil, nil, errors.New("indexer is not ready")
}
