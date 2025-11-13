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
	repo_model "github.com/kumose/kmup/models/repo"
	api "github.com/kumose/kmup/modules/structs"
)

func WebAssetDownloadURL(repo *repo_model.Repository, attach *repo_model.Attachment) string {
	return attach.DownloadURL()
}

func APIAssetDownloadURL(repo *repo_model.Repository, attach *repo_model.Attachment) string {
	return attach.DownloadURL()
}

// ToAttachment converts models.Attachment to api.Attachment for API usage
func ToAttachment(repo *repo_model.Repository, a *repo_model.Attachment) *api.Attachment {
	return toAttachment(repo, a, WebAssetDownloadURL)
}

// ToAPIAttachment converts models.Attachment to api.Attachment for API usage
func ToAPIAttachment(repo *repo_model.Repository, a *repo_model.Attachment) *api.Attachment {
	return toAttachment(repo, a, APIAssetDownloadURL)
}

// toAttachment converts models.Attachment to api.Attachment for API usage
func toAttachment(repo *repo_model.Repository, a *repo_model.Attachment, getDownloadURL func(repo *repo_model.Repository, attach *repo_model.Attachment) string) *api.Attachment {
	return &api.Attachment{
		ID:            a.ID,
		Name:          a.Name,
		Created:       a.CreatedUnix.AsTime(),
		DownloadCount: a.DownloadCount,
		Size:          a.Size,
		UUID:          a.UUID,
		DownloadURL:   getDownloadURL(repo, a), // for web request json and api request json, return different download urls
	}
}

func ToAPIAttachments(repo *repo_model.Repository, attachments []*repo_model.Attachment) []*api.Attachment {
	return toAttachments(repo, attachments, APIAssetDownloadURL)
}

func toAttachments(repo *repo_model.Repository, attachments []*repo_model.Attachment, getDownloadURL func(repo *repo_model.Repository, attach *repo_model.Attachment) string) []*api.Attachment {
	converted := make([]*api.Attachment, 0, len(attachments))
	for _, attachment := range attachments {
		converted = append(converted, toAttachment(repo, attachment, getDownloadURL))
	}
	return converted
}
