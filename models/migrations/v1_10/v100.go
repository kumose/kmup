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

package v1_10

import (
	"net/url"
	"strings"
	"time"

	"xorm.io/xorm"
)

func UpdateMigrationServiceTypes(x *xorm.Engine) error {
	type Repository struct {
		ID                  int64
		OriginalServiceType int    `xorm:"index default(0)"`
		OriginalURL         string `xorm:"VARCHAR(2048)"`
	}

	if err := x.Sync(new(Repository)); err != nil {
		return err
	}

	var last int
	const batchSize = 50
	for {
		results := make([]Repository, 0, batchSize)
		err := x.Where("original_url <> '' AND original_url IS NOT NULL").
			And("original_service_type = 0 OR original_service_type IS NULL").
			OrderBy("id").
			Limit(batchSize, last).
			Find(&results)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			break
		}
		last += len(results)

		const PlainGitService = 1 // 1 plain git service
		const GithubService = 2   // 2 github.com

		for _, res := range results {
			u, err := url.Parse(res.OriginalURL)
			if err != nil {
				return err
			}
			serviceType := PlainGitService
			if strings.EqualFold(u.Host, "github.com") {
				serviceType = GithubService
			}
			_, err = x.Exec("UPDATE repository SET original_service_type = ? WHERE id = ?", serviceType, res.ID)
			if err != nil {
				return err
			}
		}
	}

	type ExternalLoginUser struct {
		ExternalID        string         `xorm:"pk NOT NULL"`
		UserID            int64          `xorm:"INDEX NOT NULL"`
		LoginSourceID     int64          `xorm:"pk NOT NULL"`
		RawData           map[string]any `xorm:"TEXT JSON"`
		Provider          string         `xorm:"index VARCHAR(25)"`
		Email             string
		Name              string
		FirstName         string
		LastName          string
		NickName          string
		Description       string
		AvatarURL         string
		Location          string
		AccessToken       string
		AccessTokenSecret string
		RefreshToken      string
		ExpiresAt         time.Time
	}

	return x.Sync(new(ExternalLoginUser))
}
