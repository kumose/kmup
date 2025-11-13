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

package user_test

import (
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"

	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	keyName := "test_user_setting"
	assert.NoError(t, unittest.PrepareTestDatabase())

	newSetting := &user_model.Setting{UserID: 99, SettingKey: keyName, SettingValue: "Kmup User Setting Test"}

	// create setting
	err := user_model.SetUserSetting(t.Context(), newSetting.UserID, newSetting.SettingKey, newSetting.SettingValue)
	assert.NoError(t, err)
	// test about saving unchanged values
	err = user_model.SetUserSetting(t.Context(), newSetting.UserID, newSetting.SettingKey, newSetting.SettingValue)
	assert.NoError(t, err)

	// get specific setting
	settings, err := user_model.GetSettings(t.Context(), 99, []string{keyName})
	assert.NoError(t, err)
	assert.Len(t, settings, 1)
	assert.Equal(t, newSetting.SettingValue, settings[keyName].SettingValue)

	settingValue, err := user_model.GetUserSetting(t.Context(), 99, keyName)
	assert.NoError(t, err)
	assert.Equal(t, newSetting.SettingValue, settingValue)

	settingValue, err = user_model.GetUserSetting(t.Context(), 99, "no_such")
	assert.NoError(t, err)
	assert.Empty(t, settingValue)

	// updated setting
	updatedSetting := &user_model.Setting{UserID: 99, SettingKey: keyName, SettingValue: "Updated"}
	err = user_model.SetUserSetting(t.Context(), updatedSetting.UserID, updatedSetting.SettingKey, updatedSetting.SettingValue)
	assert.NoError(t, err)

	// get all settings
	settings, err = user_model.GetUserAllSettings(t.Context(), 99)
	assert.NoError(t, err)
	assert.Len(t, settings, 1)
	assert.Equal(t, updatedSetting.SettingValue, settings[updatedSetting.SettingKey].SettingValue)

	// delete setting
	err = user_model.DeleteUserSetting(t.Context(), 99, keyName)
	assert.NoError(t, err)
	settings, err = user_model.GetUserAllSettings(t.Context(), 99)
	assert.NoError(t, err)
	assert.Empty(t, settings)
}
