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

package integration

import (
	"context"
	"strings"
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/util"
	files_service "github.com/kumose/kmup/services/repository/files"

	"github.com/stretchr/testify/require"
)

type createFileInBranchOptions struct {
	OldBranch, NewBranch string
}

func testCreateFileInBranch(t *testing.T, user *user_model.User, repo *repo_model.Repository, createOpts createFileInBranchOptions, files map[string]string) *api.FilesResponse {
	resp, err := createFileInBranch(user, repo, createOpts, files)
	require.NoError(t, err)
	return resp
}

func createFileInBranch(user *user_model.User, repo *repo_model.Repository, createOpts createFileInBranchOptions, files map[string]string) (*api.FilesResponse, error) {
	ctx := context.TODO()
	opts := &files_service.ChangeRepoFilesOptions{OldBranch: createOpts.OldBranch, NewBranch: createOpts.NewBranch}
	for path, content := range files {
		opts.Files = append(opts.Files, &files_service.ChangeRepoFile{
			Operation:     "create",
			TreePath:      path,
			ContentReader: strings.NewReader(content),
		})
	}
	return files_service.ChangeRepoFiles(ctx, repo, user, opts)
}

func deleteFileInBranch(user *user_model.User, repo *repo_model.Repository, treePath, branchName string) (*api.FilesResponse, error) {
	ctx := context.TODO()
	opts := &files_service.ChangeRepoFilesOptions{
		Files: []*files_service.ChangeRepoFile{
			{
				Operation: "delete",
				TreePath:  treePath,
			},
		},
		OldBranch: branchName,
		Author:    nil,
		Committer: nil,
	}
	return files_service.ChangeRepoFiles(ctx, repo, user, opts)
}

func createOrReplaceFileInBranch(user *user_model.User, repo *repo_model.Repository, treePath, branchName, content string) error {
	_, err := deleteFileInBranch(user, repo, treePath, branchName)

	if err != nil && !files_service.IsErrRepoFileDoesNotExist(err) {
		return err
	}

	_, err = createFileInBranch(user, repo, createFileInBranchOptions{OldBranch: branchName}, map[string]string{treePath: content})
	return err
}

// TODO: replace all usages of this function with testCreateFileInBranch or testCreateFile
func createFile(user *user_model.User, repo *repo_model.Repository, treePath string, optContent ...string) (*api.FilesResponse, error) {
	content := util.OptionalArg(optContent, "This is a NEW file") // some tests need this default content because its SHA is hardcoded
	return createFileInBranch(user, repo, createFileInBranchOptions{}, map[string]string{treePath: content})
}
