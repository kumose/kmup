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

package repository

import (
	"context"
	"fmt"
	"io"

	"github.com/kumose/kmup/models/db"
	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/gitrepo"
	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/options"
	"github.com/kumose/kmup/modules/queue"

	licenseclassifier "github.com/google/licenseclassifier/v2"
)

var (
	classifier      *licenseclassifier.Classifier
	LicenseFileName = "LICENSE"

	// licenseUpdaterQueue represents a queue to handle update repo licenses
	licenseUpdaterQueue *queue.WorkerPoolQueue[*LicenseUpdaterOptions]
)

func AddRepoToLicenseUpdaterQueue(opts *LicenseUpdaterOptions) error {
	if opts == nil {
		return nil
	}
	return licenseUpdaterQueue.Push(opts)
}

func InitLicenseClassifier() error {
	// threshold should be 0.84~0.86 or the test will be failed
	classifier = licenseclassifier.NewClassifier(.85)
	licenseFiles, err := options.AssetFS().ListFiles("license", true)
	if err != nil {
		return err
	}

	for _, licenseFile := range licenseFiles {
		licenseName := licenseFile
		data, err := options.License(licenseFile)
		if err != nil {
			return err
		}
		classifier.AddContent("License", licenseName, licenseName, data)
	}
	return nil
}

type LicenseUpdaterOptions struct {
	RepoID int64
}

func repoLicenseUpdater(items ...*LicenseUpdaterOptions) []*LicenseUpdaterOptions {
	ctx := graceful.GetManager().ShutdownContext()

	for _, opts := range items {
		repo, err := repo_model.GetRepositoryByID(ctx, opts.RepoID)
		if err != nil {
			log.Error("repoLicenseUpdater [%d] failed: GetRepositoryByID: %v", opts.RepoID, err)
			continue
		}
		if repo.IsEmpty {
			continue
		}

		gitRepo, err := gitrepo.OpenRepository(ctx, repo)
		if err != nil {
			log.Error("repoLicenseUpdater [%d] failed: OpenRepository: %v", opts.RepoID, err)
			continue
		}
		defer gitRepo.Close()

		commit, err := gitRepo.GetBranchCommit(repo.DefaultBranch)
		if err != nil {
			log.Error("repoLicenseUpdater [%d] failed: GetBranchCommit: %v", opts.RepoID, err)
			continue
		}
		if err = UpdateRepoLicenses(ctx, repo, commit); err != nil {
			log.Error("repoLicenseUpdater [%d] failed: updateRepoLicenses: %v", opts.RepoID, err)
		}
	}
	return nil
}

func SyncRepoLicenses(ctx context.Context) error {
	log.Trace("Doing: SyncRepoLicenses")

	if err := db.Iterate(
		ctx,
		nil,
		func(ctx context.Context, repo *repo_model.Repository) error {
			select {
			case <-ctx.Done():
				return db.ErrCancelledf("before sync repo licenses for %s", repo.FullName())
			default:
			}
			return AddRepoToLicenseUpdaterQueue(&LicenseUpdaterOptions{RepoID: repo.ID})
		},
	); err != nil {
		log.Trace("Error: SyncRepoLicenses: %v", err)
		return err
	}

	log.Trace("Finished: SyncReposLicenses")
	return nil
}

// UpdateRepoLicenses will update repository licenses col if license file exists
func UpdateRepoLicenses(ctx context.Context, repo *repo_model.Repository, commit *git.Commit) error {
	if commit == nil {
		return nil
	}

	b, err := commit.GetBlobByPath(LicenseFileName)
	if err != nil && !git.IsErrNotExist(err) {
		return fmt.Errorf("GetBlobByPath: %w", err)
	}

	if git.IsErrNotExist(err) {
		return repo_model.CleanRepoLicenses(ctx, repo)
	}

	licenses := make([]string, 0)
	if b != nil {
		r, err := b.DataAsync()
		if err != nil {
			return err
		}
		defer r.Close()

		licenses, err = detectLicense(r)
		if err != nil {
			return fmt.Errorf("detectLicense: %w", err)
		}
	}
	return repo_model.UpdateRepoLicenses(ctx, repo, commit.ID.String(), licenses)
}

// detectLicense returns the licenses detected by the given content buff
func detectLicense(r io.Reader) ([]string, error) {
	if r == nil {
		return nil, nil
	}

	matches, err := classifier.MatchFrom(r)
	if err != nil {
		return nil, err
	}
	if len(matches.Matches) > 0 {
		results := make(container.Set[string], len(matches.Matches))
		for _, r := range matches.Matches {
			if r.MatchType == "License" && !results.Contains(r.Variant) {
				results.Add(r.Variant)
			}
		}
		return results.Values(), nil
	}
	return nil, nil
}
