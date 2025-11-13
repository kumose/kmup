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

	base "github.com/kumose/kmup/modules/migration"
)

var _ base.Downloader = &PlainGitDownloader{}

// PlainGitDownloader implements a Downloader interface to clone git from a http/https URL
type PlainGitDownloader struct {
	base.NullDownloader
	ownerName string
	repoName  string
	remoteURL string
}

// NewPlainGitDownloader creates a git Downloader
func NewPlainGitDownloader(ownerName, repoName, remoteURL string) *PlainGitDownloader {
	return &PlainGitDownloader{
		ownerName: ownerName,
		repoName:  repoName,
		remoteURL: remoteURL,
	}
}

// GetRepoInfo returns a repository information
func (g *PlainGitDownloader) GetRepoInfo(_ context.Context) (*base.Repository, error) {
	// convert github repo to stand Repo
	return &base.Repository{
		Owner:    g.ownerName,
		Name:     g.repoName,
		CloneURL: g.remoteURL,
	}, nil
}

// GetTopics return empty string slice
func (g PlainGitDownloader) GetTopics(_ context.Context) ([]string, error) {
	return []string{}, nil
}
