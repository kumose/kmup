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
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	code_indexer "github.com/kumose/kmup/modules/indexer/code"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/tests"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func resultFilenames(doc *HTMLDoc) []string {
	filenameSelections := doc.doc.Find(".repository.search").Find(".repo-search-result").Find(".header").Find("span.file")
	result := make([]string, filenameSelections.Length())
	filenameSelections.Each(func(i int, selection *goquery.Selection) {
		result[i] = selection.Text()
	})
	return result
}

func TestSearchRepo(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo, err := repo_model.GetRepositoryByOwnerAndName(t.Context(), "user2", "repo1")
	assert.NoError(t, err)

	code_indexer.UpdateRepoIndexer(repo)

	testSearch(t, "/user2/repo1/search?q=Description&page=1", []string{"README.md"})

	setting.Indexer.IncludePatterns = setting.IndexerGlobFromString("**.txt")
	setting.Indexer.ExcludePatterns = setting.IndexerGlobFromString("**/y/**")

	repo, err = repo_model.GetRepositoryByOwnerAndName(t.Context(), "user2", "glob")
	assert.NoError(t, err)

	code_indexer.UpdateRepoIndexer(repo)

	testSearch(t, "/user2/glob/search?q=loren&page=1", []string{"a.txt"})
	testSearch(t, "/user2/glob/search?q=loren&page=1&t=match", []string{"a.txt"})
	testSearch(t, "/user2/glob/search?q=file3&page=1", []string{"x/b.txt", "a.txt"})
	testSearch(t, "/user2/glob/search?q=file3&page=1&t=match", []string{"x/b.txt", "a.txt"})
	testSearch(t, "/user2/glob/search?q=file4&page=1&t=match", []string{"x/b.txt", "a.txt"})
	testSearch(t, "/user2/glob/search?q=file5&page=1&t=match", []string{"x/b.txt", "a.txt"})
}

func testSearch(t *testing.T, url string, expected []string) {
	req := NewRequest(t, "GET", url)
	resp := MakeRequest(t, req, http.StatusOK)

	filenames := resultFilenames(NewHTMLParser(t, resp.Body))
	assert.Equal(t, expected, filenames)
}
