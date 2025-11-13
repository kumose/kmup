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

package git

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/kumose/kmup/modules/git/gitcmd"
	"github.com/kumose/kmup/modules/log"
)

// ResolveReference resolves a name to a reference
func (repo *Repository) ResolveReference(name string) (string, error) {
	stdout, _, err := gitcmd.NewCommand("show-ref", "--hash").
		AddDynamicArguments(name).
		WithDir(repo.Path).
		RunStdString(repo.Ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not a valid ref") {
			return "", ErrNotExist{name, ""}
		}
		return "", err
	}
	stdout = strings.TrimSpace(stdout)
	if stdout == "" {
		return "", ErrNotExist{name, ""}
	}

	return stdout, nil
}

// GetRefCommitID returns the last commit ID string of given reference (branch or tag).
func (repo *Repository) GetRefCommitID(name string) (string, error) {
	wr, rd, cancel, err := repo.CatFileBatchCheck(repo.Ctx)
	if err != nil {
		return "", err
	}
	defer cancel()
	_, err = wr.Write([]byte(name + "\n"))
	if err != nil {
		return "", err
	}
	shaBs, _, _, err := ReadBatchLine(rd)
	if IsErrNotExist(err) {
		return "", ErrNotExist{name, ""}
	}

	return string(shaBs), nil
}

// IsCommitExist returns true if given commit exists in current repository.
func (repo *Repository) IsCommitExist(name string) bool {
	if err := ensureValidGitRepository(repo.Ctx, repo.Path); err != nil {
		log.Error("IsCommitExist: %v", err)
		return false
	}
	_, _, err := gitcmd.NewCommand("cat-file", "-e").
		AddDynamicArguments(name).
		WithDir(repo.Path).
		RunStdString(repo.Ctx)
	return err == nil
}

func (repo *Repository) getCommit(id ObjectID) (*Commit, error) {
	wr, rd, cancel, err := repo.CatFileBatch(repo.Ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	_, _ = wr.Write([]byte(id.String() + "\n"))

	return repo.getCommitFromBatchReader(wr, rd, id)
}

func (repo *Repository) getCommitFromBatchReader(wr WriteCloserError, rd *bufio.Reader, id ObjectID) (*Commit, error) {
	_, typ, size, err := ReadBatchLine(rd)
	if err != nil {
		if errors.Is(err, io.EOF) || IsErrNotExist(err) {
			return nil, ErrNotExist{ID: id.String()}
		}
		return nil, err
	}

	switch typ {
	case "missing":
		return nil, ErrNotExist{ID: id.String()}
	case "tag":
		// then we need to parse the tag
		// and load the commit
		data, err := io.ReadAll(io.LimitReader(rd, size))
		if err != nil {
			return nil, err
		}
		_, err = rd.Discard(1)
		if err != nil {
			return nil, err
		}
		tag, err := parseTagData(id.Type(), data)
		if err != nil {
			return nil, err
		}

		if _, err := wr.Write([]byte(tag.Object.String() + "\n")); err != nil {
			return nil, err
		}

		commit, err := repo.getCommitFromBatchReader(wr, rd, tag.Object)
		if err != nil {
			return nil, err
		}

		return commit, nil
	case "commit":
		commit, err := CommitFromReader(repo, id, io.LimitReader(rd, size))
		if err != nil {
			return nil, err
		}
		_, err = rd.Discard(1)
		if err != nil {
			return nil, err
		}

		return commit, nil
	default:
		log.Debug("Unknown typ: %s", typ)
		if err := DiscardFull(rd, size+1); err != nil {
			return nil, err
		}
		return nil, ErrNotExist{
			ID: id.String(),
		}
	}
}

// ConvertToGitID returns a GitHash object from a potential ID string
func (repo *Repository) ConvertToGitID(commitID string) (ObjectID, error) {
	objectFormat, err := repo.GetObjectFormat()
	if err != nil {
		return nil, err
	}
	if len(commitID) == objectFormat.FullLength() && objectFormat.IsValid(commitID) {
		ID, err := NewIDFromString(commitID)
		if err == nil {
			return ID, nil
		}
	}

	wr, rd, cancel, err := repo.CatFileBatchCheck(repo.Ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()
	_, err = wr.Write([]byte(commitID + "\n"))
	if err != nil {
		return nil, err
	}
	sha, _, _, err := ReadBatchLine(rd)
	if err != nil {
		if IsErrNotExist(err) {
			return nil, ErrNotExist{commitID, ""}
		}
		return nil, err
	}

	return MustIDFromString(string(sha)), nil
}
