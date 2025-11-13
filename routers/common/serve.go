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

package common

import (
	"io"
	"path"
	"time"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/httpcache"
	"github.com/kumose/kmup/modules/httplib"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/services/context"
)

// ServeBlob download a git.Blob
func ServeBlob(ctx *context.Base, repo *repo_model.Repository, filePath string, blob *git.Blob, lastModified *time.Time) error {
	if httpcache.HandleGenericETagTimeCache(ctx.Req, ctx.Resp, `"`+blob.ID.String()+`"`, lastModified) {
		return nil
	}

	dataRc, err := blob.DataAsync()
	if err != nil {
		return err
	}
	defer func() {
		if err = dataRc.Close(); err != nil {
			log.Error("ServeBlob: Close: %v", err)
		}
	}()

	_ = repo.LoadOwner(ctx)
	httplib.ServeContentByReader(ctx.Req, ctx.Resp, blob.Size(), dataRc, &httplib.ServeHeaderOptions{
		Filename:      path.Base(filePath),
		CacheIsPublic: !repo.IsPrivate && repo.Owner != nil && repo.Owner.Visibility == structs.VisibleTypePublic,
		CacheDuration: setting.StaticCacheTime,
	})
	return nil
}

func ServeContentByReader(ctx *context.Base, filePath string, size int64, reader io.Reader) {
	httplib.ServeContentByReader(ctx.Req, ctx.Resp, size, reader, &httplib.ServeHeaderOptions{Filename: path.Base(filePath)})
}

func ServeContentByReadSeeker(ctx *context.Base, filePath string, modTime *time.Time, reader io.ReadSeeker) {
	httplib.ServeContentByReadSeeker(ctx.Req, ctx.Resp, modTime, reader, &httplib.ServeHeaderOptions{Filename: path.Base(filePath)})
}
