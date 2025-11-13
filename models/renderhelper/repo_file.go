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

package renderhelper

import (
	"context"
	"fmt"
	"path"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/markup"
	"github.com/kumose/kmup/modules/util"
)

type RepoFile struct {
	ctx  *markup.RenderContext
	opts RepoFileOptions

	commitChecker *commitChecker
	repoLink      string
}

func (r *RepoFile) CleanUp() {
	_ = r.commitChecker.Close()
}

func (r *RepoFile) IsCommitIDExisting(commitID string) bool {
	return r.commitChecker.IsCommitIDExisting(commitID)
}

func (r *RepoFile) ResolveLink(link, preferLinkType string) (finalLink string) {
	linkType, link := markup.ParseRenderedLink(link, preferLinkType)
	switch linkType {
	case markup.LinkTypeRoot:
		finalLink = r.ctx.ResolveLinkRoot(link)
	case markup.LinkTypeRaw:
		finalLink = r.ctx.ResolveLinkRelative(path.Join(r.repoLink, "raw", r.opts.CurrentRefPath), r.opts.CurrentTreePath, link)
	case markup.LinkTypeMedia:
		finalLink = r.ctx.ResolveLinkRelative(path.Join(r.repoLink, "media", r.opts.CurrentRefPath), r.opts.CurrentTreePath, link)
	default:
		finalLink = r.ctx.ResolveLinkRelative(path.Join(r.repoLink, "src", r.opts.CurrentRefPath), r.opts.CurrentTreePath, link)
	}
	return finalLink
}

var _ markup.RenderHelper = (*RepoFile)(nil)

type RepoFileOptions struct {
	DeprecatedRepoName  string // it is only a patch for the non-standard "markup" api
	DeprecatedOwnerName string // it is only a patch for the non-standard "markup" api

	CurrentRefPath  string // eg: "branch/main"
	CurrentTreePath string // eg: "path/to/file" in the repo
}

func NewRenderContextRepoFile(ctx context.Context, repo *repo_model.Repository, opts ...RepoFileOptions) *markup.RenderContext {
	helper := &RepoFile{opts: util.OptionalArg(opts)}
	rctx := markup.NewRenderContext(ctx)
	helper.ctx = rctx
	if repo != nil {
		helper.repoLink = repo.Link()
		helper.commitChecker = newCommitChecker(ctx, repo)
		rctx = rctx.WithMetas(repo.ComposeRepoFileMetas(ctx))
	} else {
		// this is almost dead code, only to pass the incorrect tests
		helper.repoLink = fmt.Sprintf("%s/%s", helper.opts.DeprecatedOwnerName, helper.opts.DeprecatedRepoName)
		rctx = rctx.WithMetas(map[string]string{
			"user": helper.opts.DeprecatedOwnerName,
			"repo": helper.opts.DeprecatedRepoName,
		})
	}
	rctx = rctx.WithHelper(helper)
	return rctx
}
