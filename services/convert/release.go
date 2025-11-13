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

	repo_model "github.com/kumose/kmup/models/repo"
	api "github.com/kumose/kmup/modules/structs"
)

// ToAPIRelease convert a repo_model.Release to api.Release
func ToAPIRelease(ctx context.Context, repo *repo_model.Repository, r *repo_model.Release) *api.Release {
	return &api.Release{
		ID:           r.ID,
		TagName:      r.TagName,
		Target:       r.Target,
		Title:        r.Title,
		Note:         r.Note,
		URL:          r.APIURL(),
		HTMLURL:      r.HTMLURL(),
		TarURL:       r.TarURL(),
		ZipURL:       r.ZipURL(),
		UploadURL:    r.APIUploadURL(),
		IsDraft:      r.IsDraft,
		IsPrerelease: r.IsPrerelease,
		CreatedAt:    r.CreatedUnix.AsTime(),
		PublishedAt:  r.CreatedUnix.AsTime(),
		Publisher:    ToUser(ctx, r.Publisher, nil),
		Attachments:  ToAPIAttachments(repo, r.Attachments),
	}
}
