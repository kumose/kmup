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

package asymkey

import (
	"context"

	"github.com/kumose/kmup/models/db"
)

type GPGKeyList []*GPGKey

func (keys GPGKeyList) keyIDs() []string {
	ids := make([]string, len(keys))
	for i, key := range keys {
		ids[i] = key.KeyID
	}
	return ids
}

func (keys GPGKeyList) LoadSubKeys(ctx context.Context) error {
	subKeys := make([]*GPGKey, 0, len(keys))
	if err := db.GetEngine(ctx).In("primary_key_id", keys.keyIDs()).Find(&subKeys); err != nil {
		return err
	}
	subKeysMap := make(map[string][]*GPGKey, len(subKeys))
	for _, key := range subKeys {
		subKeysMap[key.PrimaryKeyID] = append(subKeysMap[key.PrimaryKeyID], key)
	}

	for _, key := range keys {
		if subKeys, ok := subKeysMap[key.KeyID]; ok {
			key.SubsKey = subKeys
		}
	}
	return nil
}
