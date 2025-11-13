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

package indexer

import (
	"context"

	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	code_indexer "github.com/kumose/kmup/modules/indexer/code"
	issue_indexer "github.com/kumose/kmup/modules/indexer/issues"
	stats_indexer "github.com/kumose/kmup/modules/indexer/stats"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/repository"
	"github.com/kumose/kmup/modules/setting"
	notify_service "github.com/kumose/kmup/services/notify"
)

type indexerNotifier struct {
	notify_service.NullNotifier
}

var _ notify_service.Notifier = &indexerNotifier{}

// NewNotifier create a new indexerNotifier notifier
func NewNotifier() notify_service.Notifier {
	return &indexerNotifier{}
}

func (r *indexerNotifier) AdoptRepository(ctx context.Context, doer, u *user_model.User, repo *repo_model.Repository) {
	r.MigrateRepository(ctx, doer, u, repo)
}

func (r *indexerNotifier) CreateIssueComment(ctx context.Context, doer *user_model.User, repo *repo_model.Repository,
	issue *issues_model.Issue, comment *issues_model.Comment, mentions []*user_model.User,
) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) NewIssue(ctx context.Context, issue *issues_model.Issue, mentions []*user_model.User) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) NewPullRequest(ctx context.Context, pr *issues_model.PullRequest, mentions []*user_model.User) {
	if err := pr.LoadIssue(ctx); err != nil {
		log.Error("LoadIssue: %v", err)
		return
	}
	issue_indexer.UpdateIssueIndexer(ctx, pr.Issue.ID)
}

func (r *indexerNotifier) UpdateComment(ctx context.Context, doer *user_model.User, c *issues_model.Comment, oldContent string) {
	if err := c.LoadIssue(ctx); err != nil {
		log.Error("LoadIssue: %v", err)
		return
	}
	issue_indexer.UpdateIssueIndexer(ctx, c.Issue.ID)
}

func (r *indexerNotifier) DeleteComment(ctx context.Context, doer *user_model.User, comment *issues_model.Comment) {
	if err := comment.LoadIssue(ctx); err != nil {
		log.Error("LoadIssue: %v", err)
		return
	}
	issue_indexer.UpdateIssueIndexer(ctx, comment.Issue.ID)
}

func (r *indexerNotifier) DeleteRepository(ctx context.Context, doer *user_model.User, repo *repo_model.Repository) {
	issue_indexer.DeleteRepoIssueIndexer(ctx, repo.ID)
	if setting.Indexer.RepoIndexerEnabled {
		code_indexer.UpdateRepoIndexer(repo)
	}
}

func (r *indexerNotifier) MigrateRepository(ctx context.Context, doer, u *user_model.User, repo *repo_model.Repository) {
	issue_indexer.UpdateRepoIndexer(ctx, repo.ID)
	if setting.Indexer.RepoIndexerEnabled && !repo.IsEmpty {
		code_indexer.UpdateRepoIndexer(repo)
	}
	if err := stats_indexer.UpdateRepoIndexer(repo); err != nil {
		log.Error("stats_indexer.UpdateRepoIndexer(%d) failed: %v", repo.ID, err)
	}
}

func (r *indexerNotifier) PushCommits(ctx context.Context, pusher *user_model.User, repo *repo_model.Repository, opts *repository.PushUpdateOptions, commits *repository.PushCommits) {
	if !opts.RefFullName.IsBranch() {
		return
	}

	if setting.Indexer.RepoIndexerEnabled && opts.RefFullName.BranchName() == repo.DefaultBranch {
		code_indexer.UpdateRepoIndexer(repo)
	}
	if err := stats_indexer.UpdateRepoIndexer(repo); err != nil {
		log.Error("stats_indexer.UpdateRepoIndexer(%d) failed: %v", repo.ID, err)
	}
}

func (r *indexerNotifier) SyncPushCommits(ctx context.Context, pusher *user_model.User, repo *repo_model.Repository, opts *repository.PushUpdateOptions, commits *repository.PushCommits) {
	if !opts.RefFullName.IsBranch() {
		return
	}

	if setting.Indexer.RepoIndexerEnabled && opts.RefFullName.BranchName() == repo.DefaultBranch {
		code_indexer.UpdateRepoIndexer(repo)
	}
	if err := stats_indexer.UpdateRepoIndexer(repo); err != nil {
		log.Error("stats_indexer.UpdateRepoIndexer(%d) failed: %v", repo.ID, err)
	}
}

func (r *indexerNotifier) ChangeDefaultBranch(ctx context.Context, repo *repo_model.Repository) {
	if setting.Indexer.RepoIndexerEnabled && !repo.IsEmpty {
		code_indexer.UpdateRepoIndexer(repo)
	}
	if err := stats_indexer.UpdateRepoIndexer(repo); err != nil {
		log.Error("stats_indexer.UpdateRepoIndexer(%d) failed: %v", repo.ID, err)
	}
}

func (r *indexerNotifier) IssueChangeContent(ctx context.Context, doer *user_model.User, issue *issues_model.Issue, oldContent string) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) IssueChangeTitle(ctx context.Context, doer *user_model.User, issue *issues_model.Issue, oldTitle string) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) IssueChangeRef(ctx context.Context, doer *user_model.User, issue *issues_model.Issue, oldRef string) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) IssueChangeStatus(ctx context.Context, doer *user_model.User, commitID string, issue *issues_model.Issue, actionComment *issues_model.Comment, closeOrReopen bool) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) IssueChangeAssignee(ctx context.Context, doer *user_model.User, issue *issues_model.Issue, assignee *user_model.User, removed bool, comment *issues_model.Comment) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) IssueChangeMilestone(ctx context.Context, doer *user_model.User, issue *issues_model.Issue, oldMilestoneID int64) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) IssueChangeLabels(ctx context.Context, doer *user_model.User, issue *issues_model.Issue,
	addedLabels, removedLabels []*issues_model.Label,
) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) IssueClearLabels(ctx context.Context, doer *user_model.User, issue *issues_model.Issue) {
	issue_indexer.UpdateIssueIndexer(ctx, issue.ID)
}

func (r *indexerNotifier) MergePullRequest(ctx context.Context, doer *user_model.User, pr *issues_model.PullRequest) {
	if err := pr.LoadIssue(ctx); err != nil {
		log.Error("LoadIssue: %v", err)
		return
	}
	issue_indexer.UpdateIssueIndexer(ctx, pr.Issue.ID)
}

func (r *indexerNotifier) AutoMergePullRequest(ctx context.Context, doer *user_model.User, pr *issues_model.PullRequest) {
	if err := pr.LoadIssue(ctx); err != nil {
		log.Error("LoadIssue: %v", err)
		return
	}
	issue_indexer.UpdateIssueIndexer(ctx, pr.Issue.ID)
}
