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

package doctor

import (
	"context"
	"os"
	"path/filepath"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/util"
)

func checkOldArchives(ctx context.Context, logger log.Logger, autofix bool) error {
	numRepos := 0
	numReposUpdated := 0
	err := iterateRepositories(ctx, func(repo *repo_model.Repository) error {
		if repo.IsEmpty {
			return nil
		}

		p := filepath.Join(repo.RepoPath(), "archives")
		isDir, err := util.IsDir(p)
		if err != nil {
			log.Warn("check if %s is directory failed: %v", p, err)
		}
		if isDir {
			numRepos++
			if autofix {
				if err := os.RemoveAll(p); err == nil {
					numReposUpdated++
				} else {
					log.Warn("remove %s failed: %v", p, err)
				}
			}
		}
		return nil
	})

	if autofix {
		logger.Info("%d / %d old archives in repository deleted", numReposUpdated, numRepos)
	} else {
		logger.Info("%d old archives in repository need to be deleted", numRepos)
	}

	return err
}

func init() {
	Register(&Check{
		Title:     "Check old archives",
		Name:      "check-old-archives",
		IsDefault: false,
		Run:       checkOldArchives,
		Priority:  7,
	})
}
