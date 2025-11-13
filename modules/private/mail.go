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
	"context"

	"github.com/kumose/kmup/modules/setting"
)

// Email structure holds a data for sending general emails
type Email struct {
	Subject string
	Message string
	To      []string
}

// SendEmail calls the internal SendEmail function
// It accepts a list of usernames.
// If DB contains these users it will send the email to them.
// If to list == nil, it's supposed to send emails to every user present in DB
func SendEmail(ctx context.Context, subject, message string, to []string) (*ResponseText, ResponseExtra) {
	reqURL := setting.LocalURL + "api/internal/mail/send"

	req := newInternalRequestAPI(ctx, reqURL, "POST", Email{
		Subject: subject,
		Message: message,
		To:      to,
	})

	return requestJSONResp(req, &ResponseText{})
}
