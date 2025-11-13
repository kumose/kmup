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

package bleve

import (
	"context"
	"errors"

	"github.com/kumose/kmup/modules/indexer/internal"
	"github.com/kumose/kmup/modules/log"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/ethantkoenig/rupture"
)

var _ internal.Indexer = &Indexer{}

// Indexer represents a basic bleve indexer implementation
type Indexer struct {
	Indexer bleve.Index

	indexDir      string
	version       int
	mappingGetter MappingGetter
}

type MappingGetter func() (mapping.IndexMapping, error)

func NewIndexer(indexDir string, version int, mappingGetter func() (mapping.IndexMapping, error)) *Indexer {
	return &Indexer{
		indexDir:      indexDir,
		version:       version,
		mappingGetter: mappingGetter,
	}
}

// Init initializes the indexer
func (i *Indexer) Init(_ context.Context) (bool, error) {
	if i == nil {
		return false, errors.New("cannot init nil indexer")
	}

	if i.Indexer != nil {
		return false, errors.New("indexer is already initialized")
	}

	indexer, version, err := openIndexer(i.indexDir, i.version)
	if err != nil {
		return false, err
	}
	if indexer != nil {
		i.Indexer = indexer
		return true, nil
	}

	if version != 0 {
		log.Warn("Found older bleve index with version %d, Kmup will remove it and rebuild", version)
	}

	indexMapping, err := i.mappingGetter()
	if err != nil {
		return false, err
	}

	indexer, err = bleve.New(i.indexDir, indexMapping)
	if err != nil {
		return false, err
	}

	if err = rupture.WriteIndexMetadata(i.indexDir, &rupture.IndexMetadata{
		Version: i.version,
	}); err != nil {
		return false, err
	}

	i.Indexer = indexer

	return false, nil
}

// Ping checks if the indexer is available
func (i *Indexer) Ping(_ context.Context) error {
	if i == nil {
		return errors.New("cannot ping nil indexer")
	}
	if i.Indexer == nil {
		return errors.New("indexer is not initialized")
	}
	return nil
}

func (i *Indexer) Close() {
	if i == nil || i.Indexer == nil {
		return
	}

	if err := i.Indexer.Close(); err != nil {
		log.Error("Failed to close bleve indexer in %q: %v", i.indexDir, err)
	}
	i.Indexer = nil
}
