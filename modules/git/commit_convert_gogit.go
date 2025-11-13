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

//go:build gogit

package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func convertPGPSignature(c *object.Commit) *CommitSignature {
	if c.PGPSignature == "" {
		return nil
	}

	var w strings.Builder
	var err error

	if _, err = fmt.Fprintf(&w, "tree %s\n", c.TreeHash.String()); err != nil {
		return nil
	}

	for _, parent := range c.ParentHashes {
		if _, err = fmt.Fprintf(&w, "parent %s\n", parent.String()); err != nil {
			return nil
		}
	}

	if _, err = fmt.Fprint(&w, "author "); err != nil {
		return nil
	}

	if err = c.Author.Encode(&w); err != nil {
		return nil
	}

	if _, err = fmt.Fprint(&w, "\ncommitter "); err != nil {
		return nil
	}

	if err = c.Committer.Encode(&w); err != nil {
		return nil
	}

	if c.Encoding != "" && c.Encoding != "UTF-8" {
		if _, err = fmt.Fprintf(&w, "\nencoding %s\n", c.Encoding); err != nil {
			return nil
		}
	}

	if _, err = fmt.Fprintf(&w, "\n\n%s", c.Message); err != nil {
		return nil
	}

	return &CommitSignature{
		Signature: c.PGPSignature,
		Payload:   w.String(),
	}
}

func convertCommit(c *object.Commit) *Commit {
	return &Commit{
		ID:            ParseGogitHash(c.Hash),
		CommitMessage: c.Message,
		Committer:     &c.Committer,
		Author:        &c.Author,
		Signature:     convertPGPSignature(c),
		Parents:       ParseGogitHashArray(c.ParentHashes),
	}
}
