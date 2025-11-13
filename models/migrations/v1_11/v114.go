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

package v1_11

import (
	"net/url"

	"xorm.io/xorm"
)

func SanitizeOriginalURL(x *xorm.Engine) error {
	type Repository struct {
		ID          int64
		OriginalURL string `xorm:"VARCHAR(2048)"`
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

		for _, res := range results {
			u, err := url.Parse(res.OriginalURL)
			if err != nil {
				// it is ok to continue here, we only care about fixing URLs that we can read
				continue
			}
			u.User = nil
			originalURL := u.String()
			_, err = x.Exec("UPDATE repository SET original_url = ? WHERE id = ?", originalURL, res.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
