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

package v1_17

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func RenameCredentialIDBytes(x *xorm.Engine) error {
	// This migration maybe rerun so that we should check if it has been run
	credentialIDExist, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "webauthn_credential", "credential_id")
	if err != nil {
		return err
	}
	if credentialIDExist {
		credentialIDBytesExists, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "webauthn_credential", "credential_id_bytes")
		if err != nil {
			return err
		}
		if !credentialIDBytesExists {
			return nil
		}
	}

	err = func() error {
		// webauthnCredential table
		type webauthnCredential struct {
			ID        int64 `xorm:"pk autoincr"`
			Name      string
			LowerName string `xorm:"unique(s)"`
			UserID    int64  `xorm:"INDEX unique(s)"`
			// Note the lack of INDEX here
			CredentialIDBytes []byte `xorm:"VARBINARY(1024)"` // CredentialID is at most 1023 bytes as per spec released 20 July 2022
			PublicKey         []byte
			AttestationType   string
			AAGUID            []byte
			SignCount         uint32 `xorm:"BIGINT"`
			CloneWarning      bool
			CreatedUnix       timeutil.TimeStamp `xorm:"INDEX created"`
			UpdatedUnix       timeutil.TimeStamp `xorm:"INDEX updated"`
		}
		sess := x.NewSession()
		defer sess.Close()
		if err := sess.Begin(); err != nil {
			return err
		}

		if err := sess.Sync(new(webauthnCredential)); err != nil {
			return fmt.Errorf("error on Sync: %w", err)
		}

		if credentialIDExist {
			// if both errors and message exist, drop message at first
			if err := base.DropTableColumns(sess, "webauthn_credential", "credential_id"); err != nil {
				return err
			}
		}

		switch {
		case setting.Database.Type.IsMySQL():
			if _, err := sess.Exec("ALTER TABLE `webauthn_credential` CHANGE credential_id_bytes credential_id VARBINARY(1024)"); err != nil {
				return err
			}
		case setting.Database.Type.IsMSSQL():
			if _, err := sess.Exec("sp_rename 'webauthn_credential.credential_id_bytes', 'credential_id', 'COLUMN'"); err != nil {
				return err
			}
		default:
			if _, err := sess.Exec("ALTER TABLE `webauthn_credential` RENAME COLUMN credential_id_bytes TO credential_id"); err != nil {
				return err
			}
		}
		return sess.Commit()
	}()
	if err != nil {
		return err
	}

	// Create webauthnCredential table
	type webauthnCredential struct {
		ID              int64 `xorm:"pk autoincr"`
		Name            string
		LowerName       string `xorm:"unique(s)"`
		UserID          int64  `xorm:"INDEX unique(s)"`
		CredentialID    []byte `xorm:"INDEX VARBINARY(1024)"` // CredentialID is at most 1023 bytes as per spec released 20 July 2022
		PublicKey       []byte
		AttestationType string
		AAGUID          []byte
		SignCount       uint32 `xorm:"BIGINT"`
		CloneWarning    bool
		CreatedUnix     timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix     timeutil.TimeStamp `xorm:"INDEX updated"`
	}
	return x.Sync(&webauthnCredential{})
}
