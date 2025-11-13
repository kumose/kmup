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

	"github.com/kumose/kmup/models/db"

	"golang.org/x/crypto/ssh"
	"xorm.io/builder"
)

// The database is used in checkKeyFingerprint. However, most of these functions probably belong in a module

// checkKeyFingerprint only checks if key fingerprint has been used as a public key,
// it is OK to use same key as deploy key for multiple repositories/users.
func checkKeyFingerprint(ctx context.Context, fingerprint string) error {
	has, err := db.Exist[PublicKey](ctx, builder.Eq{"fingerprint": fingerprint})
	if err != nil {
		return err
	} else if has {
		return ErrKeyAlreadyExist{0, fingerprint, ""}
	}
	return nil
}

func calcFingerprintNative(publicKeyContent string) (string, error) {
	// Calculate fingerprint.
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKeyContent))
	if err != nil {
		return "", err
	}
	return ssh.FingerprintSHA256(pk), nil
}

// CalcFingerprint calculate public key's fingerprint
func CalcFingerprint(publicKeyContent string) (string, error) {
	fp, err := calcFingerprintNative(publicKeyContent)
	if err != nil {
		if IsErrKeyUnableVerify(err) {
			return "", err
		}
		return "", fmt.Errorf("CalcFingerprint: %w", err)
	}
	return fp, nil
}
