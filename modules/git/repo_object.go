// Copyright 2014 The Gogs Authors. All rights reserved.
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
	"bytes"
	"io"
	"strings"

	"github.com/kumose/kmup/modules/git/gitcmd"
)

// ObjectType git object type
type ObjectType string

const (
	// ObjectCommit commit object type
	ObjectCommit ObjectType = "commit"
	// ObjectTree tree object type
	ObjectTree ObjectType = "tree"
	// ObjectBlob blob object type
	ObjectBlob ObjectType = "blob"
	// ObjectTag tag object type
	ObjectTag ObjectType = "tag"
	// ObjectBranch branch object type
	ObjectBranch ObjectType = "branch"
)

// Bytes returns the byte array for the Object Type
func (o ObjectType) Bytes() []byte {
	return []byte(o)
}

type EmptyReader struct{}

func (EmptyReader) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func (repo *Repository) GetObjectFormat() (ObjectFormat, error) {
	if repo != nil && repo.objectFormat != nil {
		return repo.objectFormat, nil
	}

	str, err := repo.hashObject(EmptyReader{}, false)
	if err != nil {
		return nil, err
	}
	hash, err := NewIDFromString(str)
	if err != nil {
		return nil, err
	}

	repo.objectFormat = hash.Type()

	return repo.objectFormat, nil
}

// HashObject takes a reader and returns hash for that reader
func (repo *Repository) HashObject(reader io.Reader) (ObjectID, error) {
	idStr, err := repo.hashObject(reader, true)
	if err != nil {
		return nil, err
	}
	return NewIDFromString(idStr)
}

func (repo *Repository) hashObject(reader io.Reader, save bool) (string, error) {
	var cmd *gitcmd.Command
	if save {
		cmd = gitcmd.NewCommand("hash-object", "-w", "--stdin")
	} else {
		cmd = gitcmd.NewCommand("hash-object", "--stdin")
	}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err := cmd.
		WithDir(repo.Path).
		WithStdin(reader).
		WithStdout(stdout).
		WithStderr(stderr).
		Run(repo.Ctx)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}
