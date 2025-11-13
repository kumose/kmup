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

package private

import (
	stdCtx "context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/private"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/mailer"
	sender_service "github.com/kumose/kmup/services/mailer/sender"
)

// SendEmail pushes messages to mail queue
//
// It doesn't wait before each message will be processed
func SendEmail(ctx *context.PrivateContext) {
	if setting.MailService == nil {
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: "Mail service is not enabled.",
		})
		return
	}

	var mail private.Email
	rd := ctx.Req.Body
	defer rd.Close()

	if err := json.NewDecoder(rd).Decode(&mail); err != nil {
		log.Error("JSON Decode failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: err.Error(),
		})
		return
	}

	var emails []string
	if len(mail.To) > 0 {
		for _, uname := range mail.To {
			user, err := user_model.GetUserByName(ctx, uname)
			if err != nil {
				err := fmt.Sprintf("Failed to get user information: %v", err)
				log.Error(err)
				ctx.JSON(http.StatusInternalServerError, private.Response{
					Err: err,
				})
				return
			}

			if user != nil && len(user.Email) > 0 {
				emails = append(emails, user.Email)
			}
		}
	} else {
		err := db.Iterate(ctx, nil, func(ctx stdCtx.Context, user *user_model.User) error {
			if len(user.Email) > 0 && user.IsActive {
				emails = append(emails, user.Email)
			}
			return nil
		})
		if err != nil {
			err := fmt.Sprintf("Failed to find users: %v", err)
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, private.Response{
				Err: err,
			})
			return
		}
	}

	sendEmail(ctx, mail.Subject, mail.Message, emails)
}

func sendEmail(ctx *context.PrivateContext, subject, message string, to []string) {
	for _, email := range to {
		msg := sender_service.NewMessage(email, subject, message)
		mailer.SendAsync(msg)
	}

	wasSent := strconv.Itoa(len(to))

	ctx.PlainText(http.StatusOK, wasSent)
}
