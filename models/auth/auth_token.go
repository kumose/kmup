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

package auth

import (
	"context"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"

	"xorm.io/builder"
)

var ErrAuthTokenNotExist = util.NewNotExistErrorf("auth token does not exist")

type AuthToken struct { //nolint:revive // export stutter
	ID          string `xorm:"pk"`
	TokenHash   string
	UserID      int64              `xorm:"INDEX"`
	ExpiresUnix timeutil.TimeStamp `xorm:"INDEX"`
}

func init() {
	db.RegisterModel(new(AuthToken))
}

func InsertAuthToken(ctx context.Context, t *AuthToken) error {
	_, err := db.GetEngine(ctx).Insert(t)
	return err
}

func GetAuthTokenByID(ctx context.Context, id string) (*AuthToken, error) {
	at := &AuthToken{}

	has, err := db.GetEngine(ctx).ID(id).Get(at)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrAuthTokenNotExist
	}
	return at, nil
}

func UpdateAuthTokenByID(ctx context.Context, t *AuthToken) error {
	_, err := db.GetEngine(ctx).ID(t.ID).Cols("token_hash", "expires_unix").Update(t)
	return err
}

func DeleteAuthTokenByID(ctx context.Context, id string) error {
	_, err := db.GetEngine(ctx).ID(id).Delete(&AuthToken{})
	return err
}

func DeleteAuthTokensByUserID(ctx context.Context, uid int64) error {
	_, err := db.GetEngine(ctx).Where(builder.Eq{"user_id": uid}).Delete(&AuthToken{})
	return err
}

func DeleteExpiredAuthTokens(ctx context.Context) error {
	_, err := db.GetEngine(ctx).Where(builder.Lt{"expires_unix": timeutil.TimeStampNow()}).Delete(&AuthToken{})
	return err
}
