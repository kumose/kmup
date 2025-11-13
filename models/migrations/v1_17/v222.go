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
	"errors"
	"fmt"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func DropOldCredentialIDColumn(x *xorm.Engine) error {
	// This migration maybe rerun so that we should check if it has been run
	credentialIDExist, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "webauthn_credential", "credential_id")
	if err != nil {
		return err
	}
	if !credentialIDExist {
		// Column is already non-extant
		return nil
	}
	credentialIDBytesExists, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "webauthn_credential", "credential_id_bytes")
	if err != nil {
		return err
	}
	if !credentialIDBytesExists {
		// looks like 221 hasn't properly run
		return errors.New("webauthn_credential does not have a credential_id_bytes column... it is not safe to run this migration")
	}

	// Create webauthnCredential table
	type webauthnCredential struct {
		ID           int64 `xorm:"pk autoincr"`
		Name         string
		LowerName    string `xorm:"unique(s)"`
		UserID       int64  `xorm:"INDEX unique(s)"`
		CredentialID string `xorm:"INDEX VARCHAR(410)"`
		// Note the lack of the INDEX on CredentialIDBytes - we will add this in v223.go
		CredentialIDBytes []byte `xorm:"VARBINARY(1024)"` // CredentialID is at most 1023 bytes as per spec released 20 July 2022
		PublicKey         []byte
		AttestationType   string
		AAGUID            []byte
		SignCount         uint32 `xorm:"BIGINT"`
		CloneWarning      bool
		CreatedUnix       timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix       timeutil.TimeStamp `xorm:"INDEX updated"`
	}
	if err := x.Sync(&webauthnCredential{}); err != nil {
		return err
	}

	// Drop the old credential ID
	sess := x.NewSession()
	defer sess.Close()

	if err := base.DropTableColumns(sess, "webauthn_credential", "credential_id"); err != nil {
		return fmt.Errorf("unable to drop old credentialID column: %w", err)
	}
	return sess.Commit()
}
