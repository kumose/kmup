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
)

// Indexer defines an basic indexer interface
type Indexer interface {
	// Init initializes the indexer
	// returns true if the index was opened/existed (with data populated), false if it was created/not-existed (with no data)
	Init(ctx context.Context) (bool, error)
	// Ping checks if the indexer is available
	Ping(ctx context.Context) error
	// Close closes the indexer
	Close()
}

// NewDummyIndexer returns a dummy indexer
func NewDummyIndexer() Indexer {
	return &dummyIndexer{}
}

type dummyIndexer struct{}

func (d *dummyIndexer) Init(ctx context.Context) (bool, error) {
	return false, errors.New("indexer is not ready")
}

func (d *dummyIndexer) Ping(ctx context.Context) error {
	return errors.New("indexer is not ready")
}

func (d *dummyIndexer) Close() {}
