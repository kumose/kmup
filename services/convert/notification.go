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

package convert

import (
	"context"
	"net/url"

	activities_model "github.com/kumose/kmup/models/activities"
	"github.com/kumose/kmup/models/perm"
	access_model "github.com/kumose/kmup/models/perm/access"
	api "github.com/kumose/kmup/modules/structs"
)

// ToNotificationThread convert a Notification to api.NotificationThread
func ToNotificationThread(ctx context.Context, n *activities_model.Notification) *api.NotificationThread {
	result := &api.NotificationThread{
		ID:        n.ID,
		Unread:    !(n.Status == activities_model.NotificationStatusRead || n.Status == activities_model.NotificationStatusPinned),
		Pinned:    n.Status == activities_model.NotificationStatusPinned,
		UpdatedAt: n.UpdatedUnix.AsTime(),
		URL:       n.APIURL(),
	}

	// since user only get notifications when he has access to use minimal access mode
	if n.Repository != nil {
		result.Repository = ToRepo(ctx, n.Repository, access_model.Permission{AccessMode: perm.AccessModeRead})

		// This permission is not correct and we should not be reporting it
		for repository := result.Repository; repository != nil; repository = repository.Parent {
			repository.Permissions = nil
		}
	}

	// handle Subject
	switch n.Source {
	case activities_model.NotificationSourceIssue:
		result.Subject = &api.NotificationSubject{Type: api.NotifySubjectIssue}
		if n.Issue != nil {
			result.Subject.Title = n.Issue.Title
			result.Subject.URL = n.Issue.APIURL(ctx)
			result.Subject.HTMLURL = n.Issue.HTMLURL(ctx)
			result.Subject.State = n.Issue.State()
			comment, err := n.Issue.GetLastComment(ctx)
			if err == nil && comment != nil {
				result.Subject.LatestCommentURL = comment.APIURL(ctx)
				result.Subject.LatestCommentHTMLURL = comment.HTMLURL(ctx)
			}
		}
	case activities_model.NotificationSourcePullRequest:
		result.Subject = &api.NotificationSubject{Type: api.NotifySubjectPull}
		if n.Issue != nil {
			result.Subject.Title = n.Issue.Title
			result.Subject.URL = n.Issue.APIURL(ctx)
			result.Subject.HTMLURL = n.Issue.HTMLURL(ctx)
			result.Subject.State = n.Issue.State()
			comment, err := n.Issue.GetLastComment(ctx)
			if err == nil && comment != nil {
				result.Subject.LatestCommentURL = comment.APIURL(ctx)
				result.Subject.LatestCommentHTMLURL = comment.HTMLURL(ctx)
			}

			if err := n.Issue.LoadPullRequest(ctx); err == nil &&
				n.Issue.PullRequest != nil &&
				n.Issue.PullRequest.HasMerged {
				result.Subject.State = "merged"
			}
		}
	case activities_model.NotificationSourceCommit:
		url := n.Repository.HTMLURL() + "/commit/" + url.PathEscape(n.CommitID)
		result.Subject = &api.NotificationSubject{
			Type:    api.NotifySubjectCommit,
			Title:   n.CommitID,
			URL:     url,
			HTMLURL: url,
		}
	case activities_model.NotificationSourceRepository:
		result.Subject = &api.NotificationSubject{
			Type:  api.NotifySubjectRepository,
			Title: n.Repository.FullName(),
			// FIXME: this is a relative URL, rather useless and inconsistent, but keeping for backwards compat
			URL:     n.Repository.Link(),
			HTMLURL: n.Repository.HTMLURL(),
		}
	}

	return result
}

// ToNotifications convert list of Notification to api.NotificationThread list
func ToNotifications(ctx context.Context, nl activities_model.NotificationList) []*api.NotificationThread {
	result := make([]*api.NotificationThread, 0, len(nl))
	for _, n := range nl {
		result = append(result, ToNotificationThread(ctx, n))
	}
	return result
}
