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

package pipeline

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/log"
)

// RevListAllObjects runs rev-list --objects --all and writes to a pipewriter
func RevListAllObjects(ctx context.Context, revListWriter *io.PipeWriter, wg *sync.WaitGroup, basePath string, errChan chan<- error) {
	defer wg.Done()
	defer revListWriter.Close()

	stderr := new(bytes.Buffer)
	var errbuf strings.Builder
	cmd := gitcmd.NewCommand("rev-list", "--objects", "--all")
	if err := cmd.WithDir(basePath).
		WithStdout(revListWriter).
		WithStderr(stderr).
		Run(ctx); err != nil {
		log.Error("git rev-list --objects --all [%s]: %v - %s", basePath, err, errbuf.String())
		err = fmt.Errorf("git rev-list --objects --all [%s]: %w - %s", basePath, err, errbuf.String())
		_ = revListWriter.CloseWithError(err)
		errChan <- err
	}
}

// RevListObjects run rev-list --objects from headSHA to baseSHA
func RevListObjects(ctx context.Context, revListWriter *io.PipeWriter, wg *sync.WaitGroup, tmpBasePath, headSHA, baseSHA string, errChan chan<- error) {
	defer wg.Done()
	defer revListWriter.Close()
	stderr := new(bytes.Buffer)
	var errbuf strings.Builder
	cmd := gitcmd.NewCommand("rev-list", "--objects").AddDynamicArguments(headSHA)
	if baseSHA != "" {
		cmd = cmd.AddArguments("--not").AddDynamicArguments(baseSHA)
	}
	if err := cmd.WithDir(tmpBasePath).
		WithStdout(revListWriter).
		WithStderr(stderr).
		Run(ctx); err != nil {
		log.Error("git rev-list [%s]: %v - %s", tmpBasePath, err, errbuf.String())
		errChan <- fmt.Errorf("git rev-list [%s]: %w - %s", tmpBasePath, err, errbuf.String())
	}
}

// BlobsFromRevListObjects reads a RevListAllObjects and only selects blobs
func BlobsFromRevListObjects(revListReader *io.PipeReader, shasToCheckWriter *io.PipeWriter, wg *sync.WaitGroup) {
	defer wg.Done()
	defer revListReader.Close()
	scanner := bufio.NewScanner(revListReader)
	defer func() {
		_ = shasToCheckWriter.CloseWithError(scanner.Err())
	}()
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		fields := strings.Split(line, " ")
		if len(fields) < 2 || len(fields[1]) == 0 {
			continue
		}
		toWrite := []byte(fields[0] + "\n")
		for len(toWrite) > 0 {
			n, err := shasToCheckWriter.Write(toWrite)
			if err != nil {
				_ = revListReader.CloseWithError(err)
				break
			}
			toWrite = toWrite[n:]
		}
	}
}
