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
	"fmt"
	"time"

	"github.com/kumose/kmup/modules/git"
)

// PullRequest defines a standard pull request information
type PullRequest struct {
	Number         int64
	Title          string
	PosterName     string `yaml:"poster_name"`
	PosterID       int64  `yaml:"poster_id"`
	PosterEmail    string `yaml:"poster_email"`
	Content        string
	Milestone      string
	State          string
	Created        time.Time
	Updated        time.Time
	Closed         *time.Time
	Labels         []*Label
	PatchURL       string `yaml:"patch_url"` // SECURITY: This must be safe to download directly from
	Merged         bool
	MergedTime     *time.Time `yaml:"merged_time"`
	MergeCommitSHA string     `yaml:"merge_commit_sha"`
	Head           PullRequestBranch
	Base           PullRequestBranch
	Assignees      []string
	IsLocked       bool `yaml:"is_locked"`
	Reactions      []*Reaction
	ForeignIndex   int64
	Context        DownloaderContext `yaml:"-"`
	EnsuredSafe    bool              `yaml:"ensured_safe"`
	IsDraft        bool              `yaml:"is_draft"`
}

func (p *PullRequest) GetLocalIndex() int64          { return p.Number }
func (p *PullRequest) GetForeignIndex() int64        { return p.ForeignIndex }
func (p *PullRequest) GetContext() DownloaderContext { return p.Context }

// IsForkPullRequest returns true if the pull request from a forked repository but not the same repository
func (p *PullRequest) IsForkPullRequest() bool {
	return p.Head.RepoFullName() != p.Base.RepoFullName()
}

// GetGitHeadRefName returns pull request relative path to head
func (p PullRequest) GetGitHeadRefName() string {
	return fmt.Sprintf("%s%d/head", git.PullPrefix, p.Number)
}

// PullRequestBranch represents a pull request branch
type PullRequestBranch struct {
	CloneURL  string `yaml:"clone_url"` // SECURITY: This must be safe to download from
	Ref       string // SECURITY: this must be a git.IsValidRefPattern
	SHA       string // SECURITY: this must be a git.IsValidSHAPattern
	RepoName  string `yaml:"repo_name"`
	OwnerName string `yaml:"owner_name"`
}

// RepoFullName returns pull request repo full name
func (p PullRequestBranch) RepoFullName() string {
	return fmt.Sprintf("%s/%s", p.OwnerName, p.RepoName)
}

// GetExternalName ExternalUserMigrated interface
func (p *PullRequest) GetExternalName() string { return p.PosterName }

// ExternalID ExternalUserMigrated interface
func (p *PullRequest) GetExternalID() int64 { return p.PosterID }
