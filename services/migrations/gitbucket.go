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
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/kumose/kmup/modules/log"
	base "github.com/kumose/kmup/modules/migration"
	"github.com/kumose/kmup/modules/structs"
)

var (
	_ base.Downloader        = &GitBucketDownloader{}
	_ base.DownloaderFactory = &GitBucketDownloaderFactory{}
)

func init() {
	RegisterDownloaderFactory(&GitBucketDownloaderFactory{})
}

// GitBucketDownloaderFactory defines a GitBucket downloader factory
type GitBucketDownloaderFactory struct{}

// New returns a Downloader related to this factory according MigrateOptions
func (f *GitBucketDownloaderFactory) New(ctx context.Context, opts base.MigrateOptions) (base.Downloader, error) {
	u, err := url.Parse(opts.CloneAddr)
	if err != nil {
		return nil, err
	}

	fields := strings.Split(u.Path, "/")
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid path: %s", u.Path)
	}
	baseURL := u.Scheme + "://" + u.Host + strings.TrimSuffix(strings.Join(fields[:len(fields)-2], "/"), "/git")

	oldOwner := fields[len(fields)-2]
	oldName := strings.TrimSuffix(fields[len(fields)-1], ".git")

	log.Trace("Create GitBucket downloader. BaseURL: %s RepoOwner: %s RepoName: %s", baseURL, oldOwner, oldName)
	return NewGitBucketDownloader(ctx, baseURL, opts.AuthUsername, opts.AuthPassword, opts.AuthToken, oldOwner, oldName), nil
}

// GitServiceType returns the type of git service
func (f *GitBucketDownloaderFactory) GitServiceType() structs.GitServiceType {
	return structs.GitBucketService
}

// GitBucketDownloader implements a Downloader interface to get repository information
// from GitBucket via GithubDownloader
type GitBucketDownloader struct {
	*GithubDownloaderV3
}

// String implements Stringer
func (g *GitBucketDownloader) String() string {
	return fmt.Sprintf("migration from gitbucket server %s %s/%s", g.baseURL, g.repoOwner, g.repoName)
}

func (g *GitBucketDownloader) LogString() string {
	if g == nil {
		return "<GitBucketDownloader nil>"
	}
	return fmt.Sprintf("<GitBucketDownloader %s %s/%s>", g.baseURL, g.repoOwner, g.repoName)
}

// NewGitBucketDownloader creates a GitBucket downloader
func NewGitBucketDownloader(ctx context.Context, baseURL, userName, password, token, repoOwner, repoName string) *GitBucketDownloader {
	githubDownloader := NewGithubDownloaderV3(ctx, baseURL, userName, password, token, repoOwner, repoName)
	// Gitbucket 4.40 uses different internal hard-coded perPage values.
	// Issues, PRs, and other major parts use 25.  Release page uses 10.
	// Some API doesn't support paging yet.  Sounds difficult, but using
	// minimum number among them worked out very well.
	githubDownloader.maxPerPage = 10
	githubDownloader.SkipReactions = true
	githubDownloader.SkipReviews = true
	return &GitBucketDownloader{
		githubDownloader,
	}
}

// SupportGetRepoComments return true if it supports get repo comments
func (g *GitBucketDownloader) SupportGetRepoComments() bool {
	return false
}
