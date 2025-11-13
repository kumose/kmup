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

//go:build gogit

package lfs

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/modules/git"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// SearchPointerBlobs scans the whole repository for LFS pointer files
func SearchPointerBlobs(ctx context.Context, repo *git.Repository, pointerChan chan<- PointerBlob, errChan chan<- error) {
	gitRepo := repo.GoGitRepo()

	err := func() error {
		blobs, err := gitRepo.BlobObjects()
		if err != nil {
			return fmt.Errorf("lfs.SearchPointerBlobs BlobObjects: %w", err)
		}

		return blobs.ForEach(func(blob *object.Blob) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if blob.Size > MetaFileMaxSize {
				return nil
			}

			reader, err := blob.Reader()
			if err != nil {
				return fmt.Errorf("lfs.SearchPointerBlobs blob.Reader: %w", err)
			}
			defer reader.Close()

			pointer, _ := ReadPointer(reader)
			if pointer.IsValid() {
				pointerChan <- PointerBlob{Hash: blob.Hash.String(), Pointer: pointer}
			}

			return nil
		})
	}()
	if err != nil {
		select {
		case <-ctx.Done():
		default:
			errChan <- err
		}
	}

	close(pointerChan)
	close(errChan)
}
