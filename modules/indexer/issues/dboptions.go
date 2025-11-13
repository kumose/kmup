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

package issues

import (
	"strings"

	"github.com/kumose/kmup/models/db"
	issues_model "github.com/kumose/kmup/models/issues"
	"github.com/kumose/kmup/modules/indexer/issues/internal"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"
)

func ToSearchOptions(keyword string, opts *issues_model.IssuesOptions) *SearchOptions {
	if opts.IssueIDs != nil {
		setting.PanicInDevOrTesting("Indexer SearchOptions doesn't support IssueIDs")
	}
	searchOpt := &SearchOptions{
		Keyword:    keyword,
		RepoIDs:    opts.RepoIDs,
		AllPublic:  opts.AllPublic,
		IsPull:     opts.IsPull,
		IsClosed:   opts.IsClosed,
		IsArchived: opts.IsArchived,
	}

	if len(opts.LabelIDs) == 1 && opts.LabelIDs[0] == 0 {
		searchOpt.NoLabelOnly = true
	} else {
		for _, labelID := range opts.LabelIDs {
			if labelID > 0 {
				searchOpt.IncludedLabelIDs = append(searchOpt.IncludedLabelIDs, labelID)
			} else {
				searchOpt.ExcludedLabelIDs = append(searchOpt.ExcludedLabelIDs, -labelID)
			}
		}
		// opts.IncludedLabelNames and opts.ExcludedLabelNames are not supported here.
		// It's not a TO DO, it's just unnecessary.
	}

	if len(opts.MilestoneIDs) == 1 && opts.MilestoneIDs[0] == db.NoConditionID {
		searchOpt.MilestoneIDs = []int64{0}
	} else {
		searchOpt.MilestoneIDs = opts.MilestoneIDs
	}

	if opts.ProjectID > 0 {
		searchOpt.ProjectID = optional.Some(opts.ProjectID)
	} else if opts.ProjectID == db.NoConditionID { // FIXME: this is inconsistent from other places
		searchOpt.ProjectID = optional.Some[int64](0) // Those issues with no project(projectid==0)
	}

	searchOpt.AssigneeID = opts.AssigneeID

	// See the comment of issues_model.SearchOptions for the reason why we need to convert
	convertID := func(id int64) optional.Option[int64] {
		if id > 0 {
			return optional.Some(id)
		}
		if id == db.NoConditionID {
			return optional.None[int64]()
		}
		return nil
	}

	searchOpt.ProjectColumnID = convertID(opts.ProjectColumnID)
	searchOpt.PosterID = opts.PosterID
	searchOpt.MentionID = convertID(opts.MentionedID)
	searchOpt.ReviewedID = convertID(opts.ReviewedID)
	searchOpt.ReviewRequestedID = convertID(opts.ReviewRequestedID)
	searchOpt.SubscriberID = convertID(opts.SubscriberID)

	if opts.UpdatedAfterUnix > 0 {
		searchOpt.UpdatedAfterUnix = optional.Some(opts.UpdatedAfterUnix)
	}
	if opts.UpdatedBeforeUnix > 0 {
		searchOpt.UpdatedBeforeUnix = optional.Some(opts.UpdatedBeforeUnix)
	}

	searchOpt.Paginator = opts.Paginator

	switch opts.SortType {
	case "", "latest":
		searchOpt.SortBy = SortByCreatedDesc
	case "oldest":
		searchOpt.SortBy = SortByCreatedAsc
	case "recentupdate":
		searchOpt.SortBy = SortByUpdatedDesc
	case "leastupdate":
		searchOpt.SortBy = SortByUpdatedAsc
	case "mostcomment":
		searchOpt.SortBy = SortByCommentsDesc
	case "leastcomment":
		searchOpt.SortBy = SortByCommentsAsc
	case "nearduedate":
		searchOpt.SortBy = SortByDeadlineAsc
	case "farduedate":
		searchOpt.SortBy = SortByDeadlineDesc
	case "priority", "priorityrepo", "project-column-sorting":
		// Unsupported sort type for search
		fallthrough
	default:
		if strings.HasPrefix(opts.SortType, issues_model.ScopeSortPrefix) {
			searchOpt.SortBy = internal.SortBy(opts.SortType)
		} else {
			searchOpt.SortBy = SortByUpdatedDesc
		}
	}

	return searchOpt
}
