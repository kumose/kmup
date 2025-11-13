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
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kumose/kmup/modules/indexer/issues/internal/tests"

	"github.com/stretchr/testify/require"
)

func TestElasticsearchIndexer(t *testing.T) {
	// The elasticsearch instance started by pull-db-tests.yml > test-unit > services > elasticsearch
	url := "http://elastic:changeme@elasticsearch:9200"

	if os.Getenv("CI") == "" {
		// Make it possible to run tests against a local elasticsearch instance
		url = os.Getenv("TEST_ELASTICSEARCH_URL")
		if url == "" {
			t.Skip("TEST_ELASTICSEARCH_URL not set and not running in CI")
			return
		}
	}

	require.Eventually(t, func() bool {
		resp, err := http.Get(url)
		return err == nil && resp.StatusCode == http.StatusOK
	}, time.Minute, time.Second, "Expected elasticsearch to be up")

	indexer := NewIndexer(url, fmt.Sprintf("test_elasticsearch_indexer_%d", time.Now().Unix()))
	defer indexer.Close()

	tests.TestIndexer(t, indexer)
}
