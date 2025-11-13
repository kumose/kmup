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

import "context"

// Uploader uploads all the information of one repository
type Uploader interface {
	MaxBatchInsertSize(tp string) int
	CreateRepo(ctx context.Context, repo *Repository, opts MigrateOptions) error
	CreateTopics(ctx context.Context, topic ...string) error
	CreateMilestones(ctx context.Context, milestones ...*Milestone) error
	CreateReleases(ctx context.Context, releases ...*Release) error
	SyncTags(ctx context.Context) error
	CreateLabels(ctx context.Context, labels ...*Label) error
	CreateIssues(ctx context.Context, issues ...*Issue) error
	CreateComments(ctx context.Context, comments ...*Comment) error
	CreatePullRequests(ctx context.Context, prs ...*PullRequest) error
	CreateReviews(ctx context.Context, reviews ...*Review) error
	Rollback() error
	Finish(ctx context.Context) error
	Close()
}
