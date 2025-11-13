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

	"github.com/kumose/kmup/modules/indexer"
	"github.com/kumose/kmup/modules/indexer/internal"
)

// Indexer defines an interface to indexer issues contents
type Indexer interface {
	internal.Indexer
	Index(ctx context.Context, issue ...*IndexerData) error
	Delete(ctx context.Context, ids ...int64) error
	Search(ctx context.Context, options *SearchOptions) (*SearchResult, error)
	SupportedSearchModes() []indexer.SearchMode
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

func (d *dummyIndexer) Index(_ context.Context, _ ...*IndexerData) error {
	return errors.New("indexer is not ready")
}

func (d *dummyIndexer) Delete(_ context.Context, _ ...int64) error {
	return errors.New("indexer is not ready")
}

func (d *dummyIndexer) Search(_ context.Context, _ *SearchOptions) (*SearchResult, error) {
	return nil, errors.New("indexer is not ready")
}
