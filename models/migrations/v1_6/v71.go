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

package v1_6

import (
	"fmt"

	"github.com/kumose/kmup/models/migrations/base"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"

	"xorm.io/xorm"
)

func AddScratchHash(x *xorm.Engine) error {
	// TwoFactor see models/twofactor.go
	type TwoFactor struct {
		ID               int64 `xorm:"pk autoincr"`
		UID              int64 `xorm:"UNIQUE"`
		Secret           string
		ScratchToken     string
		ScratchSalt      string
		ScratchHash      string
		LastUsedPasscode string             `xorm:"VARCHAR(10)"`
		CreatedUnix      timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix      timeutil.TimeStamp `xorm:"INDEX updated"`
	}

	if err := x.Sync(new(TwoFactor)); err != nil {
		return fmt.Errorf("Sync: %w", err)
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	// transform all tokens to hashes
	const batchSize = 100
	for start := 0; ; start += batchSize {
		tfas := make([]*TwoFactor, 0, batchSize)
		if err := sess.Limit(batchSize, start).Find(&tfas); err != nil {
			return err
		}
		if len(tfas) == 0 {
			break
		}

		for _, tfa := range tfas {
			// generate salt
			salt, err := util.CryptoRandomString(10)
			if err != nil {
				return err
			}
			tfa.ScratchSalt = salt
			tfa.ScratchHash = base.HashToken(tfa.ScratchToken, salt)

			if _, err := sess.ID(tfa.ID).Cols("scratch_salt, scratch_hash").Update(tfa); err != nil {
				return fmt.Errorf("couldn't add in scratch_hash and scratch_salt: %w", err)
			}
		}
	}

	// Commit and begin new transaction for dropping columns
	if err := sess.Commit(); err != nil {
		return err
	}
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := base.DropTableColumns(sess, "two_factor", "scratch_token"); err != nil {
		return err
	}
	return sess.Commit()
}
