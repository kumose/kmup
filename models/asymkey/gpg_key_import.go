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

//    __________________  ________   ____  __.
//   /  _____/\______   \/  _____/  |    |/ _|____ ___.__.
//  /   \  ___ |     ___/   \  ___  |      <_/ __ <   |  |
//  \    \_\  \|    |   \    \_\  \ |    |  \  ___/\___  |
//   \______  /|____|    \______  / |____|__ \___  > ____|
//          \/                  \/          \/   \/\/
//  .___                              __
//  |   | _____ ______   ____________/  |_
//  |   |/     \\____ \ /  _ \_  __ \   __\
//  |   |  Y Y  \  |_> >  <_> )  | \/|  |
//  |___|__|_|  /   __/ \____/|__|   |__|
//            \/|__|

// This file contains functions related to the original import of a key

// GPGKeyImport the original import of key
type GPGKeyImport struct {
	KeyID   string `xorm:"pk CHAR(16) NOT NULL"`
	Content string `xorm:"MEDIUMTEXT NOT NULL"`
}

func init() {
	db.RegisterModel(new(GPGKeyImport))
}

// GetGPGImportByKeyID returns the import public armored key by given KeyID.
func GetGPGImportByKeyID(ctx context.Context, keyID string) (*GPGKeyImport, error) {
	key := new(GPGKeyImport)
	has, err := db.GetEngine(ctx).ID(keyID).Get(key)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrGPGKeyImportNotExist{keyID}
	}
	return key, nil
}
