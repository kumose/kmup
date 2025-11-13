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
	"fmt"

	"github.com/kumose/kmup/modules/log"
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

	// The format of the index name is <index_name>_v<version>, not <index_name>.v<version> like elasticsearch.
	// Because meilisearch does not support "." in index name, it should contain only alphanumeric characters, hyphens (-) and underscores (_).
	// See https://www.meilisearch.com/docs/learn/core_concepts/indexes#index-uid

	return fmt.Sprintf("%s_v%d", indexName, version)
}

func (i *Indexer) checkOldIndexes() {
	for v := 0; v < i.version; v++ {
		indexName := versionedIndexName(i.indexName, v)
		_, err := i.Client.GetIndex(indexName)
		if err == nil {
			log.Warn("Found older meilisearch index named %q, Kmup will keep the old NOT DELETED. You can delete the old version after the upgrade succeed.", indexName)
		}
	}
}
