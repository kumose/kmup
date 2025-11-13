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
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/kumose/kmup/modules/git/gitcmd"
)

// NameRevStdin runs name-rev --stdin
func NameRevStdin(ctx context.Context, shasToNameReader *io.PipeReader, nameRevStdinWriter *io.PipeWriter, wg *sync.WaitGroup, tmpBasePath string) {
	defer wg.Done()
	defer shasToNameReader.Close()
	defer nameRevStdinWriter.Close()

	stderr := new(bytes.Buffer)
	var errbuf strings.Builder
	if err := gitcmd.NewCommand("name-rev", "--stdin", "--name-only", "--always").
		WithDir(tmpBasePath).
		WithStdin(shasToNameReader).
		WithStdout(nameRevStdinWriter).
		WithStderr(stderr).
		Run(ctx); err != nil {
		_ = shasToNameReader.CloseWithError(fmt.Errorf("git name-rev [%s]: %w - %s", tmpBasePath, err, errbuf.String()))
	}
}
