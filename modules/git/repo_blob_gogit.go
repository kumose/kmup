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
	"github.com/go-git/go-git/v5/plumbing"
)

func (repo *Repository) getBlob(id ObjectID) (*Blob, error) {
	encodedObj, err := repo.gogitRepo.Storer.EncodedObject(plumbing.AnyObject, plumbing.Hash(id.RawValue()))
	if err != nil {
		return nil, ErrNotExist{id.String(), ""}
	}

	return &Blob{
		ID:              id,
		gogitEncodedObj: encodedObj,
	}, nil
}
