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

package db

import (
	"context"

	"github.com/kumose/kmup/modules/indexer/internal"
)

var _ internal.Indexer = &Indexer{}

// Indexer represents a basic db indexer implementation
type Indexer struct{}

// Init initializes the indexer
func (i *Indexer) Init(_ context.Context) (bool, error) {
	// Return true to indicate that the index was opened/existed.
	// So that the indexer will not try to populate the index, the data is already there.
	return true, nil
}

// Ping checks if the indexer is available
func (i *Indexer) Ping(_ context.Context) error {
	// No need to ping database to check if it is available.
	// If the database goes down, Kmup will go down, so nobody will care if the indexer is available.
	return nil
}

// Close closes the indexer
func (i *Indexer) Close() {
	// nothing to do
}
