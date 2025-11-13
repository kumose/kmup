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

package gitgrep

import (
	"context"
	"fmt"
	"strings"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/indexer"
	code_indexer "github.com/kumose/kmup/modules/indexer/code"
	"github.com/kumose/kmup/modules/setting"
)

func indexSettingToGitGrepPathspecList() (list []string) {
	for _, expr := range setting.Indexer.IncludePatterns {
		list = append(list, ":(glob)"+expr.PatternString())
	}
	for _, expr := range setting.Indexer.ExcludePatterns {
		list = append(list, ":(glob,exclude)"+expr.PatternString())
	}
	return list
}

func PerformSearch(ctx context.Context, page int, repoID int64, gitRepo *git.Repository, ref git.RefName, keyword string, searchMode indexer.SearchModeType) (searchResults []*code_indexer.Result, total int, err error) {
	grepMode := git.GrepModeWords
	switch searchMode {
	case indexer.SearchModeExact:
		grepMode = git.GrepModeExact
	case indexer.SearchModeRegexp:
		grepMode = git.GrepModeRegexp
	}
	res, err := git.GrepSearch(ctx, gitRepo, keyword, git.GrepOptions{
		ContextLineNumber: 1,
		GrepMode:          grepMode,
		RefName:           ref.String(),
		PathspecList:      indexSettingToGitGrepPathspecList(),
	})
	if err != nil {
		// TODO: if no branch exists, it reports: exit status 128, fatal: this operation must be run in a work tree.
		return nil, 0, fmt.Errorf("git.GrepSearch: %w", err)
	}
	commitID, err := gitRepo.GetRefCommitID(ref.String())
	if err != nil {
		return nil, 0, fmt.Errorf("gitRepo.GetRefCommitID: %w", err)
	}

	total = len(res)
	pageStart := min((page-1)*setting.UI.RepoSearchPagingNum, len(res))
	pageEnd := min(page*setting.UI.RepoSearchPagingNum, len(res))
	res = res[pageStart:pageEnd]
	for _, r := range res {
		searchResults = append(searchResults, &code_indexer.Result{
			RepoID:   repoID,
			Filename: r.Filename,
			CommitID: commitID,
			// UpdatedUnix: not supported yet
			// Language:    not supported yet
			// Color:       not supported yet
			Lines: code_indexer.HighlightSearchResultCode(r.Filename, "", r.LineNumbers, strings.Join(r.LineCodes, "\n")),
		})
	}
	return searchResults, total, nil
}
