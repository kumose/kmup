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

package mcaptcha

import (
	"context"
	"fmt"

	"github.com/kumose/kmup/modules/setting"

	"codeberg.org/gusted/mcaptcha"
)

func Verify(ctx context.Context, token string) (bool, error) {
	valid, err := mcaptcha.Verify(ctx, &mcaptcha.VerifyOpts{
		InstanceURL: setting.Service.McaptchaURL,
		Sitekey:     setting.Service.McaptchaSitekey,
		Secret:      setting.Service.McaptchaSecret,
		Token:       token,
	})
	if err != nil {
		return false, fmt.Errorf("wasn't able to verify mCaptcha: %w", err)
	}
	return valid, nil
}
