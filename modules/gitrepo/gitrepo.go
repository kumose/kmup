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

package gitrepo

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/reqctx"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

// Repository represents a git repository which stored in a disk
type Repository interface {
	RelativePath() string // We don't assume how the directory structure of the repository is, so we only need the relative path
}

// repoPath resolves the Repository.RelativePath (which is a unix-style path like "username/reponame.git")
// to a local filesystem path according to setting.RepoRootPath
var repoPath = func(repo Repository) string {
	return filepath.Join(setting.RepoRootPath, filepath.FromSlash(repo.RelativePath()))
}

// OpenRepository opens the repository at the given relative path with the provided context.
func OpenRepository(ctx context.Context, repo Repository) (*git.Repository, error) {
	return git.OpenRepository(ctx, repoPath(repo))
}

// contextKey is a value for use with context.WithValue.
type contextKey struct {
	repoPath string
}

// RepositoryFromContextOrOpen attempts to get the repository from the context or just opens it
// The caller must call "defer gitRepo.Close()"
func RepositoryFromContextOrOpen(ctx context.Context, repo Repository) (*git.Repository, io.Closer, error) {
	reqCtx := reqctx.FromContext(ctx)
	if reqCtx != nil {
		gitRepo, err := RepositoryFromRequestContextOrOpen(reqCtx, repo)
		return gitRepo, util.NopCloser{}, err
	}
	gitRepo, err := OpenRepository(ctx, repo)
	return gitRepo, gitRepo, err
}

// RepositoryFromRequestContextOrOpen opens the repository at the given relative path in the provided request context.
// Caller shouldn't close the git repo manually, the git repo will be automatically closed when the request context is done.
func RepositoryFromRequestContextOrOpen(ctx reqctx.RequestContext, repo Repository) (*git.Repository, error) {
	ck := contextKey{repoPath: repoPath(repo)}
	if gitRepo, ok := ctx.Value(ck).(*git.Repository); ok {
		return gitRepo, nil
	}
	gitRepo, err := git.OpenRepository(ctx, ck.repoPath)
	if err != nil {
		return nil, err
	}
	ctx.AddCloser(gitRepo)
	ctx.SetContextValue(ck, gitRepo)
	return gitRepo, nil
}

// IsRepositoryExist returns true if the repository directory exists in the disk
func IsRepositoryExist(ctx context.Context, repo Repository) (bool, error) {
	return util.IsExist(repoPath(repo))
}

// DeleteRepository deletes the repository directory from the disk, it will return
// nil if the repository does not exist.
func DeleteRepository(ctx context.Context, repo Repository) error {
	return util.RemoveAll(repoPath(repo))
}

// RenameRepository renames a repository's name on disk
func RenameRepository(ctx context.Context, repo, newRepo Repository) error {
	if err := util.Rename(repoPath(repo), repoPath(newRepo)); err != nil {
		return fmt.Errorf("rename repository directory: %w", err)
	}
	return nil
}

func InitRepository(ctx context.Context, repo Repository, objectFormatName string) error {
	return git.InitRepository(ctx, repoPath(repo), true, objectFormatName)
}

func UpdateServerInfo(ctx context.Context, repo Repository) error {
	_, _, err := RunCmdBytes(ctx, repo, gitcmd.NewCommand("update-server-info"))
	return err
}

func GetRepoFS(repo Repository) fs.FS {
	return os.DirFS(repoPath(repo))
}
