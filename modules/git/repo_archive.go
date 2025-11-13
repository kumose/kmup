// Copyright 2015 The Gogs Authors. All rights reserved.
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
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/kumose/kmup/modules/git/gitcmd"
)

// ArchiveType archive types
type ArchiveType int

const (
	ArchiveUnknown ArchiveType = iota
	ArchiveZip                 // 1
	ArchiveTarGz               // 2
	ArchiveBundle              // 3
)

// String converts an ArchiveType to string: the extension of the archive file without prefix dot
func (a ArchiveType) String() string {
	switch a {
	case ArchiveZip:
		return "zip"
	case ArchiveTarGz:
		return "tar.gz"
	case ArchiveBundle:
		return "bundle"
	}
	return "unknown"
}

func SplitArchiveNameType(s string) (string, ArchiveType) {
	switch {
	case strings.HasSuffix(s, ".zip"):
		return strings.TrimSuffix(s, ".zip"), ArchiveZip
	case strings.HasSuffix(s, ".tar.gz"):
		return strings.TrimSuffix(s, ".tar.gz"), ArchiveTarGz
	case strings.HasSuffix(s, ".bundle"):
		return strings.TrimSuffix(s, ".bundle"), ArchiveBundle
	}
	return s, ArchiveUnknown
}

// CreateArchive create archive content to the target path
func (repo *Repository) CreateArchive(ctx context.Context, format ArchiveType, target io.Writer, usePrefix bool, commitID string) error {
	if format.String() == "unknown" {
		return fmt.Errorf("unknown format: %v", format)
	}

	cmd := gitcmd.NewCommand("archive")
	if usePrefix {
		cmd.AddOptionFormat("--prefix=%s", filepath.Base(strings.TrimSuffix(repo.Path, ".git"))+"/")
	}
	cmd.AddOptionFormat("--format=%s", format.String())
	cmd.AddDynamicArguments(commitID)

	var stderr strings.Builder
	err := cmd.WithDir(repo.Path).
		WithStdout(target).
		WithStderr(&stderr).
		Run(ctx)
	if err != nil {
		return gitcmd.ConcatenateError(err, stderr.String())
	}
	return nil
}
