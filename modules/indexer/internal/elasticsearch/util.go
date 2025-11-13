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

package elasticsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/kumose/kmup/modules/log"

	"github.com/olivere/elastic/v7"
)

// VersionedIndexName returns the full index name with version
func (i *Indexer) VersionedIndexName() string {
	return versionedIndexName(i.indexName, i.version)
}

func versionedIndexName(indexName string, version int) string {
	if version == 0 {
		// Old index name without version
		return indexName
	}
	return fmt.Sprintf("%s.v%d", indexName, version)
}

func (i *Indexer) createIndex(ctx context.Context) error {
	createIndex, err := i.Client.CreateIndex(i.VersionedIndexName()).BodyString(i.mapping).Do(ctx)
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return fmt.Errorf("create index %s with %s failed", i.VersionedIndexName(), i.mapping)
	}

	i.checkOldIndexes(ctx)

	return nil
}

func (i *Indexer) initClient() (*elastic.Client, error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(i.url),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10 * time.Second),
		elastic.SetGzip(false),
	}

	logger := log.GetLogger(log.DEFAULT)

	opts = append(opts, elastic.SetTraceLog(&log.PrintfLogger{Logf: logger.Trace}))
	opts = append(opts, elastic.SetInfoLog(&log.PrintfLogger{Logf: logger.Info}))
	opts = append(opts, elastic.SetErrorLog(&log.PrintfLogger{Logf: logger.Error}))

	return elastic.NewClient(opts...)
}

func (i *Indexer) checkOldIndexes(ctx context.Context) {
	for v := 0; v < i.version; v++ {
		indexName := versionedIndexName(i.indexName, v)
		exists, err := i.Client.IndexExists(indexName).Do(ctx)
		if err == nil && exists {
			log.Warn("Found older elasticsearch index named %q, Kmup will keep the old NOT DELETED. You can delete the old version after the upgrade succeed.", indexName)
		}
	}
}
