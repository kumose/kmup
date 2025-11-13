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

package git

import (
	"context"
	"path"
	"strings"

	giturl "github.com/kumose/kmup/modules/git/url"
	"github.com/kumose/kmup/modules/util"
)

// CommitSubmoduleFile represents a file with submodule type.
type CommitSubmoduleFile struct {
	repoLink string
	fullPath string
	refURL   string
	refID    string

	parsed           bool
	parsedTargetLink string
}

// NewCommitSubmoduleFile create a new submodule file
func NewCommitSubmoduleFile(repoLink, fullPath, refURL, refID string) *CommitSubmoduleFile {
	return &CommitSubmoduleFile{repoLink: repoLink, fullPath: fullPath, refURL: refURL, refID: refID}
}

// RefID returns the commit ID of the submodule, it returns empty string for nil receiver
func (sf *CommitSubmoduleFile) RefID() string {
	if sf == nil {
		return ""
	}
	return sf.refID
}

func (sf *CommitSubmoduleFile) getWebLinkInTargetRepo(ctx context.Context, moreLinkPath string) *SubmoduleWebLink {
	if sf == nil || sf.refURL == "" {
		return nil
	}
	if strings.HasPrefix(sf.refURL, "../") {
		targetLink := path.Join(sf.repoLink, sf.refURL)
		return &SubmoduleWebLink{RepoWebLink: targetLink, CommitWebLink: targetLink + moreLinkPath}
	}
	if !sf.parsed {
		sf.parsed = true
		parsedURL, err := giturl.ParseRepositoryURL(ctx, sf.refURL)
		if err != nil {
			return nil
		}
		sf.parsedTargetLink = giturl.MakeRepositoryWebLink(parsedURL)
	}
	return &SubmoduleWebLink{RepoWebLink: sf.parsedTargetLink, CommitWebLink: sf.parsedTargetLink + moreLinkPath}
}

// SubmoduleWebLinkTree tries to make the submodule's tree link in its own repo, it also works on "nil" receiver
// It returns nil if the submodule does not have a valid URL or is nil
func (sf *CommitSubmoduleFile) SubmoduleWebLinkTree(ctx context.Context, optCommitID ...string) *SubmoduleWebLink {
	return sf.getWebLinkInTargetRepo(ctx, "/tree/"+util.OptionalArg(optCommitID, sf.RefID()))
}

// SubmoduleWebLinkCompare tries to make the submodule's compare link in its own repo, it also works on "nil" receiver
// It returns nil if the submodule does not have a valid URL or is nil
func (sf *CommitSubmoduleFile) SubmoduleWebLinkCompare(ctx context.Context, commitID1, commitID2 string) *SubmoduleWebLink {
	return sf.getWebLinkInTargetRepo(ctx, "/compare/"+commitID1+"..."+commitID2)
}
