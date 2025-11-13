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

//go:build !gogit

package lfs

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/git/pipeline"
)

// SearchPointerBlobs scans the whole repository for LFS pointer files
func SearchPointerBlobs(ctx context.Context, repo *git.Repository, pointerChan chan<- PointerBlob, errChan chan<- error) {
	basePath := repo.Path

	catFileCheckReader, catFileCheckWriter := io.Pipe()
	shasToBatchReader, shasToBatchWriter := io.Pipe()
	catFileBatchReader, catFileBatchWriter := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(4)

	// Create the go-routines in reverse order.

	// 4. Take the output of cat-file --batch and check if each file in turn
	// to see if they're pointers to files in the LFS store
	go createPointerResultsFromCatFileBatch(ctx, catFileBatchReader, &wg, pointerChan)

	// 3. Take the shas of the blobs and batch read them
	go pipeline.CatFileBatch(ctx, shasToBatchReader, catFileBatchWriter, &wg, basePath)

	// 2. From the provided objects restrict to blobs <=1k
	go pipeline.BlobsLessThan1024FromCatFileBatchCheck(catFileCheckReader, shasToBatchWriter, &wg)

	// 1. Run batch-check on all objects in the repository
	if !git.DefaultFeatures().CheckVersionAtLeast("2.6.0") {
		revListReader, revListWriter := io.Pipe()
		shasToCheckReader, shasToCheckWriter := io.Pipe()
		wg.Add(2)
		go pipeline.CatFileBatchCheck(ctx, shasToCheckReader, catFileCheckWriter, &wg, basePath)
		go pipeline.BlobsFromRevListObjects(revListReader, shasToCheckWriter, &wg)
		go pipeline.RevListAllObjects(ctx, revListWriter, &wg, basePath, errChan)
	} else {
		go pipeline.CatFileBatchCheckAllObjects(ctx, catFileCheckWriter, &wg, basePath, errChan)
	}
	wg.Wait()

	close(pointerChan)
	close(errChan)
}

func createPointerResultsFromCatFileBatch(ctx context.Context, catFileBatchReader *io.PipeReader, wg *sync.WaitGroup, pointerChan chan<- PointerBlob) {
	defer wg.Done()
	defer catFileBatchReader.Close()

	bufferedReader := bufio.NewReader(catFileBatchReader)
	buf := make([]byte, 1025)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
		}

		// File descriptor line: sha
		sha, err := bufferedReader.ReadString(' ')
		if err != nil {
			_ = catFileBatchReader.CloseWithError(err)
			break
		}
		sha = strings.TrimSpace(sha)
		// Throw away the blob
		if _, err := bufferedReader.ReadString(' '); err != nil {
			_ = catFileBatchReader.CloseWithError(err)
			break
		}
		sizeStr, err := bufferedReader.ReadString('\n')
		if err != nil {
			_ = catFileBatchReader.CloseWithError(err)
			break
		}
		size, err := strconv.Atoi(sizeStr[:len(sizeStr)-1])
		if err != nil {
			_ = catFileBatchReader.CloseWithError(err)
			break
		}
		pointerBuf := buf[:size+1]
		if _, err := io.ReadFull(bufferedReader, pointerBuf); err != nil {
			_ = catFileBatchReader.CloseWithError(err)
			break
		}
		pointerBuf = pointerBuf[:size]
		// Now we need to check if the pointerBuf is an LFS pointer
		pointer, _ := ReadPointerFromBuffer(pointerBuf)
		if !pointer.IsValid() {
			continue
		}

		pointerChan <- PointerBlob{Hash: sha, Pointer: pointer}
	}
}
