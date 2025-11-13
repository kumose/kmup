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
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kumose/kmup/modules/log"
	base "github.com/kumose/kmup/modules/migration"
	"github.com/kumose/kmup/modules/structs"

	kmup_sdk "github.com/kumose-go/kmup"
)

var (
	_ base.Downloader        = &KmupDownloader{}
	_ base.DownloaderFactory = &KmupDownloaderFactory{}
)

func init() {
	RegisterDownloaderFactory(&KmupDownloaderFactory{})
}

// KmupDownloaderFactory defines a kmup downloader factory
type KmupDownloaderFactory struct{}

// New returns a Downloader related to this factory according MigrateOptions
func (f *KmupDownloaderFactory) New(ctx context.Context, opts base.MigrateOptions) (base.Downloader, error) {
	u, err := url.Parse(opts.CloneAddr)
	if err != nil {
		return nil, err
	}

	baseURL := u.Scheme + "://" + u.Host
	repoNameSpace := strings.TrimPrefix(u.Path, "/")
	repoNameSpace = strings.TrimSuffix(repoNameSpace, ".git")

	path := strings.Split(repoNameSpace, "/")
	if len(path) < 2 {
		return nil, fmt.Errorf("invalid path: %s", repoNameSpace)
	}

	repoPath := strings.Join(path[len(path)-2:], "/")
	if len(path) > 2 {
		subPath := strings.Join(path[:len(path)-2], "/")
		baseURL += "/" + subPath
	}

	log.Trace("Create kmup downloader. BaseURL: %s RepoName: %s", baseURL, repoNameSpace)

	return NewKmupDownloader(ctx, baseURL, repoPath, opts.AuthUsername, opts.AuthPassword, opts.AuthToken)
}

// GitServiceType returns the type of git service
func (f *KmupDownloaderFactory) GitServiceType() structs.GitServiceType {
	return structs.KmupService
}

// KmupDownloader implements a Downloader interface to get repository information's
type KmupDownloader struct {
	base.NullDownloader
	client     *kmup_sdk.Client
	baseURL    string
	repoOwner  string
	repoName   string
	pagination bool
	maxPerPage int
}

// NewKmupDownloader creates a kmup Downloader via kmup API
//
//	Use either a username/password or personal token. token is preferred
//	Note: Public access only allows very basic access
func NewKmupDownloader(ctx context.Context, baseURL, repoPath, username, password, token string) (*KmupDownloader, error) {
	kmupClient, err := kmup_sdk.NewClient(
		baseURL,
		kmup_sdk.SetToken(token),
		kmup_sdk.SetBasicAuth(username, password),
		kmup_sdk.SetContext(ctx),
		kmup_sdk.SetHTTPClient(NewMigrationHTTPClient()),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create NewKmupDownloader for: %s. Error: %v", baseURL, err))
		return nil, err
	}

	path := strings.Split(repoPath, "/")

	paginationSupport := true
	if err = kmupClient.CheckServerVersionConstraint(">=1.12"); err != nil {
		paginationSupport = false
	}

	// set small maxPerPage since we can only guess
	// (default would be 50 but this can differ)
	maxPerPage := 10
	// kmup instances >=1.13 can tell us what maximum they have
	apiConf, _, err := kmupClient.GetGlobalAPISettings()
	if err != nil {
		log.Info("Unable to get global API settings. Ignoring these.")
		log.Debug("kmupClient.GetGlobalAPISettings. Error: %v", err)
	}
	if apiConf != nil {
		maxPerPage = apiConf.MaxResponseItems
	}

	return &KmupDownloader{
		client:     kmupClient,
		baseURL:    baseURL,
		repoOwner:  path[0],
		repoName:   path[1],
		pagination: paginationSupport,
		maxPerPage: maxPerPage,
	}, nil
}

// String implements Stringer
func (g *KmupDownloader) String() string {
	return fmt.Sprintf("migration from kmup server %s %s/%s", g.baseURL, g.repoOwner, g.repoName)
}

func (g *KmupDownloader) LogString() string {
	if g == nil {
		return "<KmupDownloader nil>"
	}
	return fmt.Sprintf("<KmupDownloader %s %s/%s>", g.baseURL, g.repoOwner, g.repoName)
}

// GetRepoInfo returns a repository information
func (g *KmupDownloader) GetRepoInfo(_ context.Context) (*base.Repository, error) {
	if g == nil {
		return nil, errors.New("error: KmupDownloader is nil")
	}

	repo, _, err := g.client.GetRepo(g.repoOwner, g.repoName)
	if err != nil {
		return nil, err
	}

	return &base.Repository{
		Name:          repo.Name,
		Owner:         repo.Owner.UserName,
		IsPrivate:     repo.Private,
		Description:   repo.Description,
		CloneURL:      repo.CloneURL,
		OriginalURL:   repo.HTMLURL,
		DefaultBranch: repo.DefaultBranch,
	}, nil
}

// GetTopics return kmup topics
func (g *KmupDownloader) GetTopics(_ context.Context) ([]string, error) {
	topics, _, err := g.client.ListRepoTopics(g.repoOwner, g.repoName, kmup_sdk.ListRepoTopicsOptions{})
	return topics, err
}

// GetMilestones returns milestones
func (g *KmupDownloader) GetMilestones(ctx context.Context) ([]*base.Milestone, error) {
	milestones := make([]*base.Milestone, 0, g.maxPerPage)

	for i := 1; ; i++ {
		// make sure kmup can shutdown gracefully
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		ms, _, err := g.client.ListRepoMilestones(g.repoOwner, g.repoName, kmup_sdk.ListMilestoneOption{
			ListOptions: kmup_sdk.ListOptions{
				PageSize: g.maxPerPage,
				Page:     i,
			},
			State: kmup_sdk.StateAll,
		})
		if err != nil {
			return nil, err
		}

		for i := range ms {
			// old kmup instances dont have this information
			createdAT := time.Time{}
			var updatedAT *time.Time
			if ms[i].Closed != nil {
				createdAT = *ms[i].Closed
				updatedAT = ms[i].Closed
			}

			// new kmup instances (>=1.13) do
			if !ms[i].Created.IsZero() {
				createdAT = ms[i].Created
			}
			if ms[i].Updated != nil && !ms[i].Updated.IsZero() {
				updatedAT = ms[i].Updated
			}

			milestones = append(milestones, &base.Milestone{
				Title:       ms[i].Title,
				Description: ms[i].Description,
				Deadline:    ms[i].Deadline,
				Created:     createdAT,
				Updated:     updatedAT,
				Closed:      ms[i].Closed,
				State:       string(ms[i].State),
			})
		}
		if !g.pagination || len(ms) < g.maxPerPage {
			break
		}
	}
	return milestones, nil
}

func (g *KmupDownloader) convertKmupLabel(label *kmup_sdk.Label) *base.Label {
	return &base.Label{
		Name:        label.Name,
		Color:       label.Color,
		Description: label.Description,
	}
}

// GetLabels returns labels
func (g *KmupDownloader) GetLabels(ctx context.Context) ([]*base.Label, error) {
	labels := make([]*base.Label, 0, g.maxPerPage)

	for i := 1; ; i++ {
		// make sure kmup can shutdown gracefully
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		ls, _, err := g.client.ListRepoLabels(g.repoOwner, g.repoName, kmup_sdk.ListLabelsOptions{ListOptions: kmup_sdk.ListOptions{
			PageSize: g.maxPerPage,
			Page:     i,
		}})
		if err != nil {
			return nil, err
		}

		for i := range ls {
			labels = append(labels, g.convertKmupLabel(ls[i]))
		}
		if !g.pagination || len(ls) < g.maxPerPage {
			break
		}
	}
	return labels, nil
}

func (g *KmupDownloader) convertKmupRelease(rel *kmup_sdk.Release) *base.Release {
	r := &base.Release{
		TagName:         rel.TagName,
		TargetCommitish: rel.Target,
		Name:            rel.Title,
		Body:            rel.Note,
		Draft:           rel.IsDraft,
		Prerelease:      rel.IsPrerelease,
		PublisherID:     rel.Publisher.ID,
		PublisherName:   rel.Publisher.UserName,
		PublisherEmail:  rel.Publisher.Email,
		Published:       rel.PublishedAt,
		Created:         rel.CreatedAt,
	}

	httpClient := NewMigrationHTTPClient()

	for _, asset := range rel.Attachments {
		assetID := asset.ID // Don't optimize this, for closure we need a local variable
		assetDownloadURL := asset.DownloadURL
		size := int(asset.Size)
		dlCount := int(asset.DownloadCount)
		r.Assets = append(r.Assets, &base.ReleaseAsset{
			ID:            asset.ID,
			Name:          asset.Name,
			Size:          &size,
			DownloadCount: &dlCount,
			Created:       asset.Created,
			DownloadURL:   &asset.DownloadURL,
			DownloadFunc: func() (io.ReadCloser, error) {
				asset, _, err := g.client.GetReleaseAttachment(g.repoOwner, g.repoName, rel.ID, assetID)
				if err != nil {
					return nil, err
				}

				if !hasBaseURL(assetDownloadURL, g.baseURL) {
					WarnAndNotice("Unexpected AssetURL for assetID[%d] in %s: %s", assetID, g, assetDownloadURL)
					return io.NopCloser(strings.NewReader(asset.DownloadURL)), nil
				}

				// FIXME: for a private download?
				req, err := http.NewRequest(http.MethodGet, assetDownloadURL, nil)
				if err != nil {
					return nil, err
				}
				resp, err := httpClient.Do(req)
				if err != nil {
					return nil, err
				}

				// resp.Body is closed by the uploader
				return resp.Body, nil
			},
		})
	}
	return r
}

// GetReleases returns releases
func (g *KmupDownloader) GetReleases(ctx context.Context) ([]*base.Release, error) {
	releases := make([]*base.Release, 0, g.maxPerPage)

	for i := 1; ; i++ {
		// make sure kmup can shutdown gracefully
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		rl, _, err := g.client.ListReleases(g.repoOwner, g.repoName, kmup_sdk.ListReleasesOptions{ListOptions: kmup_sdk.ListOptions{
			PageSize: g.maxPerPage,
			Page:     i,
		}})
		if err != nil {
			return nil, err
		}

		for i := range rl {
			releases = append(releases, g.convertKmupRelease(rl[i]))
		}
		if !g.pagination || len(rl) < g.maxPerPage {
			break
		}
	}
	return releases, nil
}

func (g *KmupDownloader) getIssueReactions(index int64) ([]*base.Reaction, error) {
	var reactions []*base.Reaction
	if err := g.client.CheckServerVersionConstraint(">=1.11"); err != nil {
		log.Info("KmupDownloader: instance to old, skip getIssueReactions")
		return reactions, nil
	}
	rl, _, err := g.client.GetIssueReactions(g.repoOwner, g.repoName, index)
	if err != nil {
		return nil, err
	}

	for _, reaction := range rl {
		reactions = append(reactions, &base.Reaction{
			UserID:   reaction.User.ID,
			UserName: reaction.User.UserName,
			Content:  reaction.Reaction,
		})
	}
	return reactions, nil
}

func (g *KmupDownloader) getCommentReactions(commentID int64) ([]*base.Reaction, error) {
	var reactions []*base.Reaction
	if err := g.client.CheckServerVersionConstraint(">=1.11"); err != nil {
		log.Info("KmupDownloader: instance to old, skip getCommentReactions")
		return reactions, nil
	}
	rl, _, err := g.client.GetIssueCommentReactions(g.repoOwner, g.repoName, commentID)
	if err != nil {
		return nil, err
	}

	for i := range rl {
		reactions = append(reactions, &base.Reaction{
			UserID:   rl[i].User.ID,
			UserName: rl[i].User.UserName,
			Content:  rl[i].Reaction,
		})
	}
	return reactions, nil
}

// GetIssues returns issues according start and limit
func (g *KmupDownloader) GetIssues(_ context.Context, page, perPage int) ([]*base.Issue, bool, error) {
	if perPage > g.maxPerPage {
		perPage = g.maxPerPage
	}
	allIssues := make([]*base.Issue, 0, perPage)

	issues, _, err := g.client.ListRepoIssues(g.repoOwner, g.repoName, kmup_sdk.ListIssueOption{
		ListOptions: kmup_sdk.ListOptions{Page: page, PageSize: perPage},
		State:       kmup_sdk.StateAll,
		Type:        kmup_sdk.IssueTypeIssue,
	})
	if err != nil {
		return nil, false, fmt.Errorf("error while listing issues: %w", err)
	}
	for _, issue := range issues {
		labels := make([]*base.Label, 0, len(issue.Labels))
		for i := range issue.Labels {
			labels = append(labels, g.convertKmupLabel(issue.Labels[i]))
		}

		var milestone string
		if issue.Milestone != nil {
			milestone = issue.Milestone.Title
		}

		reactions, err := g.getIssueReactions(issue.Index)
		if err != nil {
			WarnAndNotice("Unable to load reactions during migrating issue #%d in %s. Error: %v", issue.Index, g, err)
		}

		var assignees []string
		for i := range issue.Assignees {
			assignees = append(assignees, issue.Assignees[i].UserName)
		}

		allIssues = append(allIssues, &base.Issue{
			Title:        issue.Title,
			Number:       issue.Index,
			PosterID:     issue.Poster.ID,
			PosterName:   issue.Poster.UserName,
			PosterEmail:  issue.Poster.Email,
			Content:      issue.Body,
			Milestone:    milestone,
			State:        string(issue.State),
			Created:      issue.Created,
			Updated:      issue.Updated,
			Closed:       issue.Closed,
			Reactions:    reactions,
			Labels:       labels,
			Assignees:    assignees,
			IsLocked:     issue.IsLocked,
			ForeignIndex: issue.Index,
		})
	}

	isEnd := len(issues) < perPage
	if !g.pagination {
		isEnd = len(issues) == 0
	}
	return allIssues, isEnd, nil
}

// GetComments returns comments according issueNumber
func (g *KmupDownloader) GetComments(ctx context.Context, commentable base.Commentable) ([]*base.Comment, bool, error) {
	allComments := make([]*base.Comment, 0, g.maxPerPage)

	for i := 1; ; i++ {
		// make sure kmup can shutdown gracefully
		select {
		case <-ctx.Done():
			return nil, false, nil
		default:
		}

		comments, _, err := g.client.ListIssueComments(g.repoOwner, g.repoName, commentable.GetForeignIndex(), kmup_sdk.ListIssueCommentOptions{ListOptions: kmup_sdk.ListOptions{
			PageSize: g.maxPerPage,
			Page:     i,
		}})
		if err != nil {
			return nil, false, fmt.Errorf("error while listing comments for issue #%d. Error: %w", commentable.GetForeignIndex(), err)
		}

		for _, comment := range comments {
			reactions, err := g.getCommentReactions(comment.ID)
			if err != nil {
				WarnAndNotice("Unable to load comment reactions during migrating issue #%d for comment %d in %s. Error: %v", commentable.GetForeignIndex(), comment.ID, g, err)
			}

			allComments = append(allComments, &base.Comment{
				IssueIndex:  commentable.GetLocalIndex(),
				Index:       comment.ID,
				PosterID:    comment.Poster.ID,
				PosterName:  comment.Poster.UserName,
				PosterEmail: comment.Poster.Email,
				Content:     comment.Body,
				Created:     comment.Created,
				Updated:     comment.Updated,
				Reactions:   reactions,
			})
		}

		if !g.pagination || len(comments) < g.maxPerPage {
			break
		}
	}
	return allComments, true, nil
}

// GetPullRequests returns pull requests according page and perPage
func (g *KmupDownloader) GetPullRequests(_ context.Context, page, perPage int) ([]*base.PullRequest, bool, error) {
	if perPage > g.maxPerPage {
		perPage = g.maxPerPage
	}
	allPRs := make([]*base.PullRequest, 0, perPage)

	prs, _, err := g.client.ListRepoPullRequests(g.repoOwner, g.repoName, kmup_sdk.ListPullRequestsOptions{
		ListOptions: kmup_sdk.ListOptions{
			Page:     page,
			PageSize: perPage,
		},
		State: kmup_sdk.StateAll,
	})
	if err != nil {
		return nil, false, fmt.Errorf("error while listing pull requests (page: %d, pagesize: %d). Error: %w", page, perPage, err)
	}
	for _, pr := range prs {
		var milestone string
		if pr.Milestone != nil {
			milestone = pr.Milestone.Title
		}

		labels := make([]*base.Label, 0, len(pr.Labels))
		for i := range pr.Labels {
			labels = append(labels, g.convertKmupLabel(pr.Labels[i]))
		}

		var (
			headUserName string
			headRepoName string
			headCloneURL string
			headRef      string
			headSHA      string
		)
		if pr.Head != nil {
			if pr.Head.Repository != nil {
				headUserName = pr.Head.Repository.Owner.UserName
				headRepoName = pr.Head.Repository.Name
				headCloneURL = pr.Head.Repository.CloneURL
			}
			headSHA = pr.Head.Sha
			headRef = pr.Head.Ref
		}

		var mergeCommitSHA string
		if pr.MergedCommitID != nil {
			mergeCommitSHA = *pr.MergedCommitID
		}

		reactions, err := g.getIssueReactions(pr.Index)
		if err != nil {
			WarnAndNotice("Unable to load reactions during migrating pull #%d in %s. Error: %v", pr.Index, g, err)
		}

		var assignees []string
		for i := range pr.Assignees {
			assignees = append(assignees, pr.Assignees[i].UserName)
		}

		createdAt := time.Time{}
		if pr.Created != nil {
			createdAt = *pr.Created
		}
		updatedAt := time.Time{}
		if pr.Created != nil {
			updatedAt = *pr.Updated
		}

		closedAt := pr.Closed
		if pr.Merged != nil && closedAt == nil {
			closedAt = pr.Merged
		}

		allPRs = append(allPRs, &base.PullRequest{
			Title:          pr.Title,
			Number:         pr.Index,
			PosterID:       pr.Poster.ID,
			PosterName:     pr.Poster.UserName,
			PosterEmail:    pr.Poster.Email,
			Content:        pr.Body,
			State:          string(pr.State),
			Created:        createdAt,
			Updated:        updatedAt,
			Closed:         closedAt,
			Labels:         labels,
			Milestone:      milestone,
			Reactions:      reactions,
			Assignees:      assignees,
			Merged:         pr.HasMerged,
			MergedTime:     pr.Merged,
			MergeCommitSHA: mergeCommitSHA,
			IsLocked:       pr.IsLocked,
			PatchURL:       pr.PatchURL,
			Head: base.PullRequestBranch{
				Ref:       headRef,
				SHA:       headSHA,
				RepoName:  headRepoName,
				OwnerName: headUserName,
				CloneURL:  headCloneURL,
			},
			Base: base.PullRequestBranch{
				Ref:       pr.Base.Ref,
				SHA:       pr.Base.Sha,
				RepoName:  g.repoName,
				OwnerName: g.repoOwner,
			},
			ForeignIndex: pr.Index,
		})
		// SECURITY: Ensure that the PR is safe
		_ = CheckAndEnsureSafePR(allPRs[len(allPRs)-1], g.baseURL, g)
	}

	isEnd := len(prs) < perPage
	if !g.pagination {
		isEnd = len(prs) == 0
	}
	return allPRs, isEnd, nil
}

// GetReviews returns pull requests review
func (g *KmupDownloader) GetReviews(ctx context.Context, reviewable base.Reviewable) ([]*base.Review, error) {
	if err := g.client.CheckServerVersionConstraint(">=1.12"); err != nil {
		log.Info("KmupDownloader: instance to old, skip GetReviews")
		return nil, nil
	}

	allReviews := make([]*base.Review, 0, g.maxPerPage)

	for i := 1; ; i++ {
		// make sure kmup can shutdown gracefully
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		prl, _, err := g.client.ListPullReviews(g.repoOwner, g.repoName, reviewable.GetForeignIndex(), kmup_sdk.ListPullReviewsOptions{ListOptions: kmup_sdk.ListOptions{
			Page:     i,
			PageSize: g.maxPerPage,
		}})
		if err != nil {
			return nil, err
		}

		for _, pr := range prl {
			if pr.Reviewer == nil {
				// Presumably this is a team review which we cannot migrate at present but we have to skip this review as otherwise the review will be mapped on to an incorrect user.
				// TODO: handle team reviews
				continue
			}

			rcl, _, err := g.client.ListPullReviewComments(g.repoOwner, g.repoName, reviewable.GetForeignIndex(), pr.ID)
			if err != nil {
				return nil, err
			}
			var reviewComments []*base.ReviewComment
			for i := range rcl {
				line := int(rcl[i].LineNum)
				if rcl[i].OldLineNum > 0 {
					line = int(rcl[i].OldLineNum) * -1
				}

				reviewComments = append(reviewComments, &base.ReviewComment{
					ID:        rcl[i].ID,
					Content:   rcl[i].Body,
					TreePath:  rcl[i].Path,
					DiffHunk:  rcl[i].DiffHunk,
					Line:      line,
					CommitID:  rcl[i].CommitID,
					PosterID:  rcl[i].Reviewer.ID,
					CreatedAt: rcl[i].Created,
					UpdatedAt: rcl[i].Updated,
				})
			}

			review := &base.Review{
				ID:           pr.ID,
				IssueIndex:   reviewable.GetLocalIndex(),
				ReviewerID:   pr.Reviewer.ID,
				ReviewerName: pr.Reviewer.UserName,
				Official:     pr.Official,
				CommitID:     pr.CommitID,
				Content:      pr.Body,
				CreatedAt:    pr.Submitted,
				State:        string(pr.State),
				Comments:     reviewComments,
			}

			allReviews = append(allReviews, review)
		}

		if len(prl) < g.maxPerPage {
			break
		}
	}
	return allReviews, nil
}
