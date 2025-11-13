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
	"fmt"

	"github.com/kumose/kmup/modules/setting"
)

// UpdatePublicKeyInRepo update public key and if necessary deploy key updates
func UpdatePublicKeyInRepo(ctx context.Context, keyID, repoID int64) error {
	// Ask for running deliver hook and test pull request tasks.
	reqURL := setting.LocalURL + fmt.Sprintf("api/internal/ssh/%d/update/%d", keyID, repoID)
	req := newInternalRequestAPI(ctx, reqURL, "POST")
	_, extra := requestJSONResp(req, &ResponseText{})
	return extra.Error
}

// AuthorizedPublicKeyByContent searches content as prefix (leak e-mail part)
// and returns public key found.
func AuthorizedPublicKeyByContent(ctx context.Context, content string) (*ResponseText, ResponseExtra) {
	// Ask for running deliver hook and test pull request tasks.
	reqURL := setting.LocalURL + "api/internal/ssh/authorized_keys"
	req := newInternalRequestAPI(ctx, reqURL, "POST")
	req.Param("content", content)
	return requestJSONResp(req, &ResponseText{})
}
