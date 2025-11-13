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

package migration

import (
	"context"

	"github.com/kumose/kmup/modules/structs"
)

// Downloader downloads the site repo information
type Downloader interface {
	GetRepoInfo(ctx context.Context) (*Repository, error)
	GetTopics(ctx context.Context) ([]string, error)
	GetMilestones(ctx context.Context) ([]*Milestone, error)
	GetReleases(ctx context.Context) ([]*Release, error)
	GetLabels(ctx context.Context) ([]*Label, error)
	GetIssues(ctx context.Context, page, perPage int) ([]*Issue, bool, error)
	GetComments(ctx context.Context, commentable Commentable) ([]*Comment, bool, error)
	GetAllComments(ctx context.Context, page, perPage int) ([]*Comment, bool, error)
	SupportGetRepoComments() bool
	GetPullRequests(ctx context.Context, page, perPage int) ([]*PullRequest, bool, error)
	GetReviews(ctx context.Context, reviewable Reviewable) ([]*Review, error)
	FormatCloneURL(opts MigrateOptions, remoteAddr string) (string, error)
}

// DownloaderFactory defines an interface to match a downloader implementation and create a downloader
type DownloaderFactory interface {
	New(ctx context.Context, opts MigrateOptions) (Downloader, error)
	GitServiceType() structs.GitServiceType
}

// DownloaderContext has opaque information only relevant to a given downloader
type DownloaderContext any
