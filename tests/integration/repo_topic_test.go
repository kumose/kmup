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

package integration

import (
	"net/http"
	"net/url"
	"testing"

	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestTopicSearch(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	searchURL, _ := url.Parse("/explore/topics/search")
	var topics struct {
		TopicNames []*api.TopicResponse `json:"topics"`
	}

	// search all topics
	res := MakeRequest(t, NewRequest(t, "GET", searchURL.String()), http.StatusOK)
	DecodeJSON(t, res, &topics)
	assert.Len(t, topics.TopicNames, 6)
	assert.Equal(t, "6", res.Header().Get("x-total-count"))

	// pagination search topics
	topics.TopicNames = nil
	query := url.Values{"page": []string{"1"}, "limit": []string{"4"}}

	searchURL.RawQuery = query.Encode()
	res = MakeRequest(t, NewRequest(t, "GET", searchURL.String()), http.StatusOK)
	DecodeJSON(t, res, &topics)
	assert.Len(t, topics.TopicNames, 4)
	assert.Equal(t, "6", res.Header().Get("x-total-count"))

	// second page
	topics.TopicNames = nil
	query = url.Values{"page": []string{"2"}, "limit": []string{"4"}}

	searchURL.RawQuery = query.Encode()
	res = MakeRequest(t, NewRequest(t, "GET", searchURL.String()), http.StatusOK)
	DecodeJSON(t, res, &topics)
	assert.Len(t, topics.TopicNames, 2)
	assert.Equal(t, "6", res.Header().Get("x-total-count"))

	// add keyword search
	topics.TopicNames = nil
	query = url.Values{"page": []string{"1"}, "limit": []string{"4"}}
	query.Add("q", "topic")
	searchURL.RawQuery = query.Encode()
	res = MakeRequest(t, NewRequest(t, "GET", searchURL.String()), http.StatusOK)
	DecodeJSON(t, res, &topics)
	assert.Len(t, topics.TopicNames, 2)

	topics.TopicNames = nil
	query.Set("q", "database")
	searchURL.RawQuery = query.Encode()
	res = MakeRequest(t, NewRequest(t, "GET", searchURL.String()), http.StatusOK)
	DecodeJSON(t, res, &topics)
	if assert.Len(t, topics.TopicNames, 1) {
		assert.EqualValues(t, 2, topics.TopicNames[0].ID)
		assert.Equal(t, "database", topics.TopicNames[0].Name)
		assert.Equal(t, 1, topics.TopicNames[0].RepoCount)
	}
}
