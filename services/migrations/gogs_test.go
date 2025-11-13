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

package migrations

import (
	"net/http"
	"os"
	"testing"
	"time"

	base "github.com/kumose/kmup/modules/migration"

	"github.com/stretchr/testify/assert"
)

func TestGogsDownloadRepo(t *testing.T) {
	// Skip tests if Gogs token is not found
	gogsPersonalAccessToken := os.Getenv("GOGS_READ_TOKEN")
	if len(gogsPersonalAccessToken) == 0 {
		t.Skip("skipped test because GOGS_READ_TOKEN was not in the environment")
	}

	resp, err := http.Get("https://try.gogs.io/lunnytest/TESTREPO")
	if err != nil || resp.StatusCode/100 != 2 {
		// skip and don't run test
		t.Skipf("visit test repo failed, ignored")
		return
	}
	ctx := t.Context()
	downloader := NewGogsDownloader(ctx, "https://try.gogs.io", "", "", gogsPersonalAccessToken, "lunnytest", "TESTREPO")
	repo, err := downloader.GetRepoInfo(ctx)
	assert.NoError(t, err)

	assertRepositoryEqual(t, &base.Repository{
		Name:          "TESTREPO",
		Owner:         "lunnytest",
		Description:   "",
		CloneURL:      "https://try.gogs.io/lunnytest/TESTREPO.git",
		OriginalURL:   "https://try.gogs.io/lunnytest/TESTREPO",
		DefaultBranch: "master",
	}, repo)

	milestones, err := downloader.GetMilestones(ctx)
	assert.NoError(t, err)
	assertMilestonesEqual(t, []*base.Milestone{
		{
			Title: "1.0",
			State: "open",
		},
	}, milestones)

	labels, err := downloader.GetLabels(ctx)
	assert.NoError(t, err)
	assertLabelsEqual(t, []*base.Label{
		{
			Name:  "bug",
			Color: "ee0701",
		},
		{
			Name:  "duplicate",
			Color: "cccccc",
		},
		{
			Name:  "enhancement",
			Color: "84b6eb",
		},
		{
			Name:  "help wanted",
			Color: "128a0c",
		},
		{
			Name:  "invalid",
			Color: "e6e6e6",
		},
		{
			Name:  "question",
			Color: "cc317c",
		},
		{
			Name:  "wontfix",
			Color: "ffffff",
		},
	}, labels)

	// downloader.GetIssues()
	issues, isEnd, err := downloader.GetIssues(ctx, 1, 8)
	assert.NoError(t, err)
	assert.False(t, isEnd)
	assertIssuesEqual(t, []*base.Issue{
		{
			Number:      1,
			PosterID:    5331,
			PosterName:  "lunny",
			PosterEmail: "xiaolunwen@gmail.com",
			Title:       "test",
			Content:     "test",
			Milestone:   "",
			State:       "open",
			Created:     time.Date(2019, 6, 11, 8, 16, 44, 0, time.UTC),
			Updated:     time.Date(2019, 10, 26, 11, 7, 2, 0, time.UTC),
			Labels: []*base.Label{
				{
					Name:  "bug",
					Color: "ee0701",
				},
			},
		},
	}, issues)

	// downloader.GetComments()
	comments, _, err := downloader.GetComments(ctx, &base.Issue{Number: 1, ForeignIndex: 1})
	assert.NoError(t, err)
	assertCommentsEqual(t, []*base.Comment{
		{
			IssueIndex:  1,
			PosterID:    5331,
			PosterName:  "lunny",
			PosterEmail: "xiaolunwen@gmail.com",
			Created:     time.Date(2019, 6, 11, 8, 19, 50, 0, time.UTC),
			Updated:     time.Date(2019, 6, 11, 8, 19, 50, 0, time.UTC),
			Content:     "1111",
		},
		{
			IssueIndex:  1,
			PosterID:    15822,
			PosterName:  "clacplouf",
			PosterEmail: "test1234@dbn.re",
			Created:     time.Date(2019, 10, 26, 11, 7, 2, 0, time.UTC),
			Updated:     time.Date(2019, 10, 26, 11, 7, 2, 0, time.UTC),
			Content:     "88888888",
		},
	}, comments)

	// downloader.GetPullRequests()
	_, _, err = downloader.GetPullRequests(ctx, 1, 3)
	assert.Error(t, err)
}
