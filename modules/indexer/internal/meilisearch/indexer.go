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

package meilisearch

import (
	"context"
	"errors"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

// Indexer represents a basic meilisearch indexer implementation
type Indexer struct {
	Client meilisearch.ServiceManager

	url, apiKey string
	indexName   string
	version     int
	settings    *meilisearch.Settings
}

func NewIndexer(url, apiKey, indexName string, version int, settings *meilisearch.Settings) *Indexer {
	return &Indexer{
		url:       url,
		apiKey:    apiKey,
		indexName: indexName,
		version:   version,
		settings:  settings,
	}
}

// Init initializes the indexer
func (i *Indexer) Init(_ context.Context) (bool, error) {
	if i == nil {
		return false, errors.New("cannot init nil indexer")
	}

	if i.Client != nil {
		return false, errors.New("indexer is already initialized")
	}

	i.Client = meilisearch.New(i.url, meilisearch.WithAPIKey(i.apiKey))
	_, err := i.Client.GetIndex(i.VersionedIndexName())
	if err == nil {
		return true, nil
	}
	_, err = i.Client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        i.VersionedIndexName(),
		PrimaryKey: "id",
	})
	if err != nil {
		return false, err
	}

	i.checkOldIndexes()

	_, err = i.Client.Index(i.VersionedIndexName()).UpdateSettings(i.settings)
	return false, err
}

// Ping checks if the indexer is available
func (i *Indexer) Ping(ctx context.Context) error {
	if i == nil {
		return errors.New("cannot ping nil indexer")
	}
	if i.Client == nil {
		return errors.New("indexer is not initialized")
	}
	resp, err := i.Client.Health()
	if err != nil {
		return err
	}
	if resp.Status != "available" {
		// See https://docs.meilisearch.com/reference/api/health.html#status
		return fmt.Errorf("status of meilisearch is not available: %s", resp.Status)
	}
	return nil
}

// Close closes the indexer
func (i *Indexer) Close() {
	if i == nil {
		return
	}
	i.Client = nil
}
