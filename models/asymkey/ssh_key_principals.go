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
	"strings"

	"github.com/kumose/kmup/models/db"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

// CheckPrincipalKeyString strips spaces and returns an error if the given principal contains newlines
func CheckPrincipalKeyString(ctx context.Context, user *user_model.User, content string) (_ string, err error) {
	if setting.SSH.Disabled {
		return "", db.ErrSSHDisabled{}
	}

	content = strings.TrimSpace(content)
	if strings.ContainsAny(content, "\r\n") {
		return "", util.NewInvalidArgumentErrorf("only a single line with a single principal please")
	}

	// check all the allowed principals, email, username or anything
	// if any matches, return ok
	for _, v := range setting.SSH.AuthorizedPrincipalsAllow {
		switch v {
		case "anything":
			return content, nil
		case "email":
			emails, err := user_model.GetEmailAddresses(ctx, user.ID)
			if err != nil {
				return "", err
			}
			for _, email := range emails {
				if !email.IsActivated {
					continue
				}
				if content == email.Email {
					return content, nil
				}
			}

		case "username":
			if content == user.Name {
				return content, nil
			}
		}
	}

	return "", fmt.Errorf("didn't match allowed principals: %s", setting.SSH.AuthorizedPrincipalsAllow)
}
