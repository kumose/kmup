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

package activitypub

import (
	"context"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/util"
)

const rsaBits = 3072

// GetKeyPair function returns a user's private and public keys
func GetKeyPair(ctx context.Context, user *user_model.User) (pub, priv string, err error) {
	var settings map[string]*user_model.Setting
	settings, err = user_model.GetSettings(ctx, user.ID, []string{user_model.UserActivityPubPrivPem, user_model.UserActivityPubPubPem})
	if err != nil {
		return pub, priv, err
	} else if len(settings) == 0 {
		if priv, pub, err = util.GenerateKeyPair(rsaBits); err != nil {
			return pub, priv, err
		}
		if err = user_model.SetUserSetting(ctx, user.ID, user_model.UserActivityPubPrivPem, priv); err != nil {
			return pub, priv, err
		}
		if err = user_model.SetUserSetting(ctx, user.ID, user_model.UserActivityPubPubPem, pub); err != nil {
			return pub, priv, err
		}
		return pub, priv, err
	}
	priv = settings[user_model.UserActivityPubPrivPem].SettingValue
	pub = settings[user_model.UserActivityPubPubPem].SettingValue
	return pub, priv, err
}

// GetPublicKey function returns a user's public key
func GetPublicKey(ctx context.Context, user *user_model.User) (pub string, err error) {
	pub, _, err = GetKeyPair(ctx, user)
	return pub, err
}

// GetPrivateKey function returns a user's private key
func GetPrivateKey(ctx context.Context, user *user_model.User) (priv string, err error) {
	_, priv, err = GetKeyPair(ctx, user)
	return priv, err
}
