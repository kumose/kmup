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
	"fmt"
	"os"
	"path/filepath"
	"time"

	asymkey_model "github.com/kumose/kmup/models/asymkey"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

// RewriteAllPublicKeys removes any authorized key and rewrite all keys from database again.
// Note: db.GetEngine(ctx).Iterate does not get latest data after insert/delete, so we have to call this function
// outside any session scope independently.
func RewriteAllPublicKeys(ctx context.Context) error {
	// Don't rewrite key if internal server
	if setting.SSH.StartBuiltinServer || !setting.SSH.CreateAuthorizedKeysFile {
		return nil
	}

	return asymkey_model.WithSSHOpLocker(func() error {
		return rewriteAllPublicKeys(ctx)
	})
}

func rewriteAllPublicKeys(ctx context.Context) error {
	if setting.SSH.RootPath != "" {
		// First of ensure that the RootPath is present, and if not make it with 0700 permissions
		// This of course doesn't guarantee that this is the right directory for authorized_keys
		// but at least if it's supposed to be this directory and it doesn't exist and we're the
		// right user it will at least be created properly.
		err := os.MkdirAll(setting.SSH.RootPath, 0o700)
		if err != nil {
			log.Error("Unable to MkdirAll(%s): %v", setting.SSH.RootPath, err)
			return err
		}
	}

	fPath := filepath.Join(setting.SSH.RootPath, "authorized_keys")
	tmpPath := fPath + ".tmp"
	t, err := os.OpenFile(tmpPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		t.Close()
		if err := util.Remove(tmpPath); err != nil {
			log.Warn("Unable to remove temporary authorized keys file: %s: Error: %v", tmpPath, err)
		}
	}()

	if setting.SSH.AuthorizedKeysBackup {
		isExist, err := util.IsExist(fPath)
		if err != nil {
			log.Error("Unable to check if %s exists. Error: %v", fPath, err)
			return err
		}
		if isExist {
			bakPath := fmt.Sprintf("%s_%d.kmup_bak", fPath, time.Now().Unix())
			if err = util.CopyFile(fPath, bakPath); err != nil {
				return err
			}
		}
	}

	if err := asymkey_model.RegeneratePublicKeys(ctx, t); err != nil {
		return err
	}

	t.Close()
	return util.Rename(tmpPath, fPath)
}
