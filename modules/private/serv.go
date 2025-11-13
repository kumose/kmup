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
	"net/url"

	asymkey_model "github.com/kumose/kmup/models/asymkey"
	"github.com/kumose/kmup/models/perm"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
)

// KeyAndOwner is the response from ServNoCommand
type KeyAndOwner struct {
	Key   *asymkey_model.PublicKey `json:"key"`
	Owner *user_model.User         `json:"user"`
}

// ServNoCommand returns information about the provided key
func ServNoCommand(ctx context.Context, keyID int64) (*asymkey_model.PublicKey, *user_model.User, error) {
	reqURL := setting.LocalURL + fmt.Sprintf("api/internal/serv/none/%d", keyID)
	req := newInternalRequestAPI(ctx, reqURL, "GET")
	keyAndOwner, extra := requestJSONResp(req, &KeyAndOwner{})
	if extra.HasError() {
		return nil, nil, extra.Error
	}
	return keyAndOwner.Key, keyAndOwner.Owner, nil
}

// ServCommandResults are the results of a call to the private route serv
type ServCommandResults struct {
	IsWiki      bool
	DeployKeyID int64
	KeyID       int64  // public key
	KeyName     string // this field is ambiguous, it can be the name of DeployKey, or the name of the PublicKey
	UserName    string
	UserEmail   string
	UserID      int64
	OwnerName   string
	RepoName    string
	RepoID      int64
}

// ServCommand preps for a serv call
func ServCommand(ctx context.Context, keyID int64, ownerName, repoName string, mode perm.AccessMode, verb, lfsVerb string) (*ServCommandResults, ResponseExtra) {
	reqURL := setting.LocalURL + fmt.Sprintf("api/internal/serv/command/%d/%s/%s?mode=%d",
		keyID,
		url.PathEscape(ownerName),
		url.PathEscape(repoName),
		mode,
	)
	reqURL += "&verb=" + url.QueryEscape(verb)
	// reqURL += "&lfs_verb=" + url.QueryEscape(lfsVerb) // TODO: actually there is no use of this parameter. In the future, the URL construction should be more flexible
	_ = lfsVerb
	req := newInternalRequestAPI(ctx, reqURL, "GET")
	return requestJSONResp(req, &ServCommandResults{})
}
