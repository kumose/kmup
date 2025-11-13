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

package private

import (
	"fmt"
	"net/http"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/private"
	kmup_context "github.com/kumose/kmup/services/context"
)

// This file contains common functions relating to setting the Repository for the internal routes

// RepoAssignment assigns the repository and git repository to the private context
func RepoAssignment(ctx *kmup_context.PrivateContext) {
	ownerName := ctx.PathParam("owner")
	repoName := ctx.PathParam("repo")

	repo := loadRepository(ctx, ownerName, repoName)
	if ctx.Written() {
		// Error handled in loadRepository
		return
	}

	gitRepo, err := gitrepo.RepositoryFromRequestContextOrOpen(ctx, repo)
	if err != nil {
		log.Error("Failed to open repository: %s/%s Error: %v", ownerName, repoName, err)
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: fmt.Sprintf("Failed to open repository: %s/%s Error: %v", ownerName, repoName, err),
		})
		return
	}
	ctx.Repo = &kmup_context.Repository{
		Repository: repo,
		GitRepo:    gitRepo,
	}
}

func loadRepository(ctx *kmup_context.PrivateContext, ownerName, repoName string) *repo_model.Repository {
	repo, err := repo_model.GetRepositoryByOwnerAndName(ctx, ownerName, repoName)
	if err != nil {
		log.Error("Failed to get repository: %s/%s Error: %v", ownerName, repoName, err)
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: fmt.Sprintf("Failed to get repository: %s/%s Error: %v", ownerName, repoName, err),
		})
		return nil
	}
	if repo.OwnerName == "" {
		repo.OwnerName = ownerName
	}
	return repo
}
