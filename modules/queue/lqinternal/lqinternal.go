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

package lqinternal

import (
	"bytes"
	"encoding/binary"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func QueueItemIDBytes(id int64) []byte {
	buf := make([]byte, 8)
	binary.PutVarint(buf, id)
	return buf
}

func QueueItemKeyBytes(prefix []byte, id int64) []byte {
	key := make([]byte, len(prefix), len(prefix)+1+8)
	copy(key, prefix)
	key = append(key, '-')
	return append(key, QueueItemIDBytes(id)...)
}

func RemoveLevelQueueKeys(db *leveldb.DB, namePrefix []byte) {
	keyPrefix := make([]byte, len(namePrefix)+1)
	copy(keyPrefix, namePrefix)
	keyPrefix[len(namePrefix)] = '-'

	it := db.NewIterator(nil, &opt.ReadOptions{Strict: opt.NoStrict})
	defer it.Release()
	for it.Next() {
		if bytes.HasPrefix(it.Key(), keyPrefix) {
			_ = db.Delete(it.Key(), nil)
		}
	}
}

func ListLevelQueueKeys(db *leveldb.DB) (res [][]byte) {
	it := db.NewIterator(nil, &opt.ReadOptions{Strict: opt.NoStrict})
	defer it.Release()
	for it.Next() {
		res = append(res, it.Key())
	}
	return res
}
