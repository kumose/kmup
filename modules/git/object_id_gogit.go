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
	"github.com/go-git/go-git/v5/plumbing/hash"
)

func ParseGogitHash(h plumbing.Hash) ObjectID {
	switch hash.Size {
	case 20:
		return Sha1ObjectFormat.MustID(h[:])
	case 32:
		return Sha256ObjectFormat.MustID(h[:])
	}

	return nil
}

func ParseGogitHashArray(objectIDs []plumbing.Hash) []ObjectID {
	ret := make([]ObjectID, len(objectIDs))
	for i, h := range objectIDs {
		ret[i] = ParseGogitHash(h)
	}

	return ret
}
