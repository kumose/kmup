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
	"bytes"
	"os"
	"strings"
	"time"

	"github.com/kumose/kmup/modules/git/gitcmd"
)

// CommitTreeOpts represents the possible options to CommitTree
type CommitTreeOpts struct {
	Parents    []string
	Message    string
	Key        *SigningKey
	NoGPGSign  bool
	AlwaysSign bool
}

// CommitTree creates a commit from a given tree id for the user with provided message
func (repo *Repository) CommitTree(author, committer *Signature, tree *Tree, opts CommitTreeOpts) (ObjectID, error) {
	commitTimeStr := time.Now().Format(time.RFC3339)

	// Because this may call hooks we should pass in the environment
	env := append(os.Environ(),
		"GIT_AUTHOR_NAME="+author.Name,
		"GIT_AUTHOR_EMAIL="+author.Email,
		"GIT_AUTHOR_DATE="+commitTimeStr,
		"GIT_COMMITTER_NAME="+committer.Name,
		"GIT_COMMITTER_EMAIL="+committer.Email,
		"GIT_COMMITTER_DATE="+commitTimeStr,
	)
	cmd := gitcmd.NewCommand("commit-tree").AddDynamicArguments(tree.ID.String())

	for _, parent := range opts.Parents {
		cmd.AddArguments("-p").AddDynamicArguments(parent)
	}

	messageBytes := new(bytes.Buffer)
	_, _ = messageBytes.WriteString(opts.Message)
	_, _ = messageBytes.WriteString("\n")

	if opts.Key != nil {
		if opts.Key.Format != "" {
			cmd.AddConfig("gpg.format", opts.Key.Format)
		}
		cmd.AddOptionFormat("-S%s", opts.Key.KeyID)
	} else if opts.AlwaysSign {
		cmd.AddOptionFormat("-S")
	}

	if opts.NoGPGSign {
		cmd.AddArguments("--no-gpg-sign")
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err := cmd.WithEnv(env).
		WithDir(repo.Path).
		WithStdin(messageBytes).
		WithStdout(stdout).
		WithStderr(stderr).
		Run(repo.Ctx)
	if err != nil {
		return nil, gitcmd.ConcatenateError(err, stderr.String())
	}
	return NewIDFromString(strings.TrimSpace(stdout.String()))
}
