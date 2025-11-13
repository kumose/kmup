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

package gitdiff

import (
	"context"
	"html/template"

	"github.com/kumose/kmup/modules/base"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/htmlutil"
	"github.com/kumose/kmup/modules/log"
)

type SubmoduleDiffInfo struct {
	SubmoduleName string
	SubmoduleFile *git.CommitSubmoduleFile // it might be nil if the submodule is not found or unable to parse
	NewRefID      string
	PreviousRefID string
}

func (si *SubmoduleDiffInfo) PopulateURL(repoLink string, diffFile *DiffFile, leftCommit, rightCommit *git.Commit) {
	si.SubmoduleName = diffFile.Name
	submoduleCommit := rightCommit // If the submodule is added or updated, check at the right commit
	if diffFile.IsDeleted {
		submoduleCommit = leftCommit // If the submodule is deleted, check at the left commit
	}
	if submoduleCommit == nil {
		return
	}

	submoduleFullPath := diffFile.GetDiffFileName()
	submodule, err := submoduleCommit.GetSubModule(submoduleFullPath)
	if err != nil {
		log.Error("Unable to PopulateURL for submodule %q: GetSubModule: %v", submoduleFullPath, err)
		return // ignore the error, do not cause 500 errors for end users
	}
	if submodule != nil {
		si.SubmoduleFile = git.NewCommitSubmoduleFile(repoLink, submoduleFullPath, submodule.URL, submoduleCommit.ID.String())
	}
}

func (si *SubmoduleDiffInfo) CommitRefIDLinkHTML(ctx context.Context, commitID string) template.HTML {
	webLink := si.SubmoduleFile.SubmoduleWebLinkTree(ctx, commitID)
	if webLink == nil {
		return htmlutil.HTMLFormat("%s", base.ShortSha(commitID))
	}
	return htmlutil.HTMLFormat(`<a href="%s">%s</a>`, webLink.CommitWebLink, base.ShortSha(commitID))
}

func (si *SubmoduleDiffInfo) CompareRefIDLinkHTML(ctx context.Context) template.HTML {
	webLink := si.SubmoduleFile.SubmoduleWebLinkCompare(ctx, si.PreviousRefID, si.NewRefID)
	if webLink == nil {
		return htmlutil.HTMLFormat("%s...%s", base.ShortSha(si.PreviousRefID), base.ShortSha(si.NewRefID))
	}
	return htmlutil.HTMLFormat(`<a href="%s">%s...%s</a>`, webLink.CommitWebLink, base.ShortSha(si.PreviousRefID), base.ShortSha(si.NewRefID))
}

func (si *SubmoduleDiffInfo) SubmoduleRepoLinkHTML(ctx context.Context) template.HTML {
	webLink := si.SubmoduleFile.SubmoduleWebLinkTree(ctx)
	if webLink == nil {
		return htmlutil.HTMLFormat("%s", si.SubmoduleName)
	}
	return htmlutil.HTMLFormat(`<a href="%s">%s</a>`, webLink.RepoWebLink, si.SubmoduleName)
}
