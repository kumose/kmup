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
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/log"
)

type TemplateSubmoduleCommit struct {
	Path   string
	Commit string
}

// GetTemplateSubmoduleCommits returns a list of submodules paths and their commits from a repository
// This function is only for generating new repos based on existing template, the template couldn't be too large.
func GetTemplateSubmoduleCommits(ctx context.Context, repoPath string) (submoduleCommits []TemplateSubmoduleCommit, _ error) {
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	err = gitcmd.NewCommand("ls-tree", "-r", "--", "HEAD").
		WithDir(repoPath).
		WithStdout(stdoutWriter).
		WithPipelineFunc(func(ctx context.Context, cancel context.CancelFunc) error {
			_ = stdoutWriter.Close()
			defer stdoutReader.Close()

			scanner := bufio.NewScanner(stdoutReader)
			for scanner.Scan() {
				entry, err := parseLsTreeLine(scanner.Bytes())
				if err != nil {
					cancel()
					return err
				}
				if entry.EntryMode == EntryModeCommit {
					submoduleCommits = append(submoduleCommits, TemplateSubmoduleCommit{Path: entry.Name, Commit: entry.ID.String()})
				}
			}
			return scanner.Err()
		}).
		Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetTemplateSubmoduleCommits: error running git ls-tree: %v", err)
	}
	return submoduleCommits, nil
}

// AddTemplateSubmoduleIndexes Adds the given submodules to the git index.
// It is only for generating new repos based on existing template, requires the .gitmodules file to be already present in the work dir.
func AddTemplateSubmoduleIndexes(ctx context.Context, repoPath string, submodules []TemplateSubmoduleCommit) error {
	for _, submodule := range submodules {
		cmd := gitcmd.NewCommand("update-index", "--add", "--cacheinfo", "160000").AddDynamicArguments(submodule.Commit, submodule.Path)
		if stdout, _, err := cmd.WithDir(repoPath).RunStdString(ctx); err != nil {
			log.Error("Unable to add %s as submodule to repo %s: stdout %s\nError: %v", submodule.Path, repoPath, stdout, err)
			return err
		}
	}
	return nil
}
