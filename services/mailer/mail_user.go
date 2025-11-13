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

package mailer

import (
	"bytes"
	"fmt"

	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/templates"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/translation"
	sender_service "github.com/kumose/kmup/services/mailer/sender"
)

const (
	mailAuthActivate       templates.TplName = "user/auth/activate"
	mailAuthActivateEmail  templates.TplName = "user/auth/activate_email"
	mailAuthResetPassword  templates.TplName = "user/auth/reset_passwd"
	mailAuthRegisterNotify templates.TplName = "user/auth/register_notify"
)

// sendUserMail sends a mail to the user
func sendUserMail(language string, u *user_model.User, tpl templates.TplName, code, subject, info string) {
	locale := translation.NewLocale(language)
	data := map[string]any{
		"locale":            locale,
		"DisplayName":       u.DisplayName(),
		"ActiveCodeLives":   timeutil.MinutesToFriendly(setting.Service.ActiveCodeLives, locale),
		"ResetPwdCodeLives": timeutil.MinutesToFriendly(setting.Service.ResetPwdCodeLives, locale),
		"Code":              code,
		"Language":          locale.Language(),
	}

	var content bytes.Buffer

	if err := LoadedTemplates().BodyTemplates.ExecuteTemplate(&content, string(tpl), data); err != nil {
		log.Error("Template: %v", err)
		return
	}

	msg := sender_service.NewMessage(u.EmailTo(), subject, content.String())
	msg.Info = fmt.Sprintf("UID: %d, %s", u.ID, info)

	SendAsync(msg)
}

// SendActivateAccountMail sends an activation mail to the user (new user registration)
func SendActivateAccountMail(locale translation.Locale, u *user_model.User) {
	if setting.MailService == nil {
		// No mail service configured
		return
	}
	opts := &user_model.TimeLimitCodeOptions{Purpose: user_model.TimeLimitCodeActivateAccount}
	sendUserMail(locale.Language(), u, mailAuthActivate, user_model.GenerateUserTimeLimitCode(opts, u), locale.TrString("mail.activate_account"), "activate account")
}

// SendResetPasswordMail sends a password reset mail to the user
func SendResetPasswordMail(u *user_model.User) {
	if setting.MailService == nil {
		// No mail service configured
		return
	}
	locale := translation.NewLocale(u.Language)
	opts := &user_model.TimeLimitCodeOptions{Purpose: user_model.TimeLimitCodeResetPassword}
	sendUserMail(u.Language, u, mailAuthResetPassword, user_model.GenerateUserTimeLimitCode(opts, u), locale.TrString("mail.reset_password"), "recover account")
}

// SendActivateEmailMail sends confirmation email to confirm new email address
func SendActivateEmailMail(u *user_model.User, email string) {
	if setting.MailService == nil {
		// No mail service configured
		return
	}
	locale := translation.NewLocale(u.Language)
	opts := &user_model.TimeLimitCodeOptions{Purpose: user_model.TimeLimitCodeActivateEmail, NewEmail: email}
	data := map[string]any{
		"locale":          locale,
		"DisplayName":     u.DisplayName(),
		"ActiveCodeLives": timeutil.MinutesToFriendly(setting.Service.ActiveCodeLives, locale),
		"Code":            user_model.GenerateUserTimeLimitCode(opts, u),
		"Email":           email,
		"Language":        locale.Language(),
	}

	var content bytes.Buffer

	if err := LoadedTemplates().BodyTemplates.ExecuteTemplate(&content, string(mailAuthActivateEmail), data); err != nil {
		log.Error("Template: %v", err)
		return
	}

	msg := sender_service.NewMessage(email, locale.TrString("mail.activate_email"), content.String())
	msg.Info = fmt.Sprintf("UID: %d, activate email", u.ID)

	SendAsync(msg)
}

// SendRegisterNotifyMail triggers a notify e-mail by admin created a account.
func SendRegisterNotifyMail(u *user_model.User) {
	if setting.MailService == nil || !u.IsActive {
		// No mail service configured OR user is inactive
		return
	}
	locale := translation.NewLocale(u.Language)

	data := map[string]any{
		"locale":      locale,
		"DisplayName": u.DisplayName(),
		"Username":    u.Name,
		"Language":    locale.Language(),
	}

	var content bytes.Buffer

	if err := LoadedTemplates().BodyTemplates.ExecuteTemplate(&content, string(mailAuthRegisterNotify), data); err != nil {
		log.Error("Template: %v", err)
		return
	}

	msg := sender_service.NewMessage(u.EmailTo(), locale.TrString("mail.register_notify", setting.AppName), content.String())
	msg.Info = fmt.Sprintf("UID: %d, registration notify", u.ID)

	SendAsync(msg)
}
