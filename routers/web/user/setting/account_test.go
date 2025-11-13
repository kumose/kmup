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

package setting

import (
	"net/http"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/contexttest"
	"github.com/kumose/kmup/services/forms"

	"github.com/stretchr/testify/assert"
)

func TestChangePassword(t *testing.T) {
	oldPassword := "password"
	setting.MinPasswordLength = 6
	pcALL := []string{"lower", "upper", "digit", "spec"}
	pcLUN := []string{"lower", "upper", "digit"}
	pcLU := []string{"lower", "upper"}

	for _, req := range []struct {
		OldPassword        string
		NewPassword        string
		Retype             string
		Message            string
		PasswordComplexity []string
	}{
		{
			OldPassword:        oldPassword,
			NewPassword:        "Qwerty123456-",
			Retype:             "Qwerty123456-",
			Message:            "",
			PasswordComplexity: pcALL,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "12345",
			Retype:             "12345",
			Message:            "auth.password_too_short",
			PasswordComplexity: pcALL,
		},
		{
			OldPassword:        "12334",
			NewPassword:        "123456",
			Retype:             "123456",
			Message:            "settings.password_incorrect",
			PasswordComplexity: pcALL,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "123456",
			Retype:             "12345",
			Message:            "form.password_not_match",
			PasswordComplexity: pcALL,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "Qwerty",
			Retype:             "Qwerty",
			Message:            "form.password_complexity",
			PasswordComplexity: pcALL,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "Qwerty",
			Retype:             "Qwerty",
			Message:            "form.password_complexity",
			PasswordComplexity: pcLUN,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "QWERTY",
			Retype:             "QWERTY",
			Message:            "form.password_complexity",
			PasswordComplexity: pcLU,
		},
	} {
		t.Run(req.OldPassword+"__"+req.NewPassword, func(t *testing.T) {
			unittest.PrepareTestEnv(t)
			setting.PasswordComplexity = req.PasswordComplexity
			ctx, _ := contexttest.MockContext(t, "user/settings/security")
			contexttest.LoadUser(t, ctx, 2)
			contexttest.LoadRepo(t, ctx, 1)

			web.SetForm(ctx, &forms.ChangePasswordForm{
				OldPassword: req.OldPassword,
				Password:    req.NewPassword,
				Retype:      req.Retype,
			})
			AccountPost(ctx)

			assert.Contains(t, ctx.Flash.ErrorMsg, req.Message)
			assert.Equal(t, http.StatusSeeOther, ctx.Resp.WrittenStatus())
		})
	}
}
