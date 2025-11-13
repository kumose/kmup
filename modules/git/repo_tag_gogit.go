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
	"github.com/kumose/kmup/modules/log"

	"github.com/go-git/go-git/v5/plumbing"
)

// IsTagExist returns true if given tag exists in the repository.
func (repo *Repository) IsTagExist(name string) bool {
	_, err := repo.gogitRepo.Reference(plumbing.ReferenceName(TagPrefix+name), true)
	return err == nil
}

// GetTagType gets the type of the tag, either commit (simple) or tag (annotated)
func (repo *Repository) GetTagType(id ObjectID) (string, error) {
	// Get tag type
	obj, err := repo.gogitRepo.Object(plumbing.AnyObject, plumbing.Hash(id.RawValue()))
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return "", &ErrNotExist{ID: id.String()}
		}
		return "", err
	}

	return obj.Type().String(), nil
}

func (repo *Repository) getTag(tagID ObjectID, name string) (*Tag, error) {
	t, ok := repo.tagCache.Get(tagID.String())
	if ok {
		log.Debug("Hit cache: %s", tagID)
		tagClone := *t
		tagClone.Name = name // This is necessary because lightweight tags may have same id
		return &tagClone, nil
	}

	tp, err := repo.GetTagType(tagID)
	if err != nil {
		return nil, err
	}

	// Get the commit ID and tag ID (may be different for annotated tag) for the returned tag object
	commitIDStr, err := repo.GetTagCommitID(name)
	if err != nil {
		// every tag should have a commit ID so return all errors
		return nil, err
	}
	commitID, err := NewIDFromString(commitIDStr)
	if err != nil {
		return nil, err
	}

	// If type is "commit, the tag is a lightweight tag
	if ObjectType(tp) == ObjectCommit {
		commit, err := repo.GetCommit(commitIDStr)
		if err != nil {
			return nil, err
		}
		tag := &Tag{
			Name:    name,
			ID:      tagID,
			Object:  commitID,
			Type:    tp,
			Tagger:  commit.Committer,
			Message: commit.Message(),
		}

		repo.tagCache.Set(tagID.String(), tag)
		return tag, nil
	}

	gogitTag, err := repo.gogitRepo.TagObject(plumbing.Hash(tagID.RawValue()))
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return nil, &ErrNotExist{ID: tagID.String()}
		}

		return nil, err
	}

	tag := &Tag{
		Name:    name,
		ID:      tagID,
		Object:  commitID.Type().MustID(gogitTag.Target[:]),
		Type:    tp,
		Tagger:  &gogitTag.Tagger,
		Message: gogitTag.Message,
	}

	repo.tagCache.Set(tagID.String(), tag)
	return tag, nil
}
