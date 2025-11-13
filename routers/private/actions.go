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
	"errors"
	"fmt"
	"net/http"
	"strings"

	actions_model "github.com/kumose/kmup/models/actions"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/private"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/services/context"
)

// GenerateActionsRunnerToken generates a new runner token for a given scope
func GenerateActionsRunnerToken(ctx *context.PrivateContext) {
	var genRequest private.GenerateTokenRequest
	rd := ctx.Req.Body
	defer rd.Close()

	if err := json.NewDecoder(rd).Decode(&genRequest); err != nil {
		log.Error("JSON Decode failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: err.Error(),
		})
		return
	}

	owner, repo, err := parseScope(ctx, genRequest.Scope)
	if err != nil {
		log.Error("parseScope failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: err.Error(),
		})
	}

	token, err := actions_model.GetLatestRunnerToken(ctx, owner, repo)
	if errors.Is(err, util.ErrNotExist) || (token != nil && !token.IsActive) {
		token, err = actions_model.NewRunnerToken(ctx, owner, repo)
		if err != nil {
			errMsg := fmt.Sprintf("error while creating runner token: %v", err)
			log.Error("NewRunnerToken failed: %v", errMsg)
			ctx.JSON(http.StatusInternalServerError, private.Response{
				Err: errMsg,
			})
			return
		}
	} else if err != nil {
		errMsg := fmt.Sprintf("could not get unactivated runner token: %v", err)
		log.Error("GetLatestRunnerToken failed: %v", errMsg)
		ctx.JSON(http.StatusInternalServerError, private.Response{
			Err: errMsg,
		})
		return
	}

	ctx.PlainText(http.StatusOK, token.Token)
}

func parseScope(ctx *context.PrivateContext, scope string) (ownerID, repoID int64, err error) {
	ownerID = 0
	repoID = 0
	if scope == "" {
		return ownerID, repoID, nil
	}

	ownerName, repoName, found := strings.Cut(scope, "/")

	u, err := user_model.GetUserByName(ctx, ownerName)
	if err != nil {
		return ownerID, repoID, err
	}
	ownerID = u.ID

	if !found {
		return ownerID, repoID, nil
	}

	r, err := repo_model.GetRepositoryByName(ctx, u.ID, repoName)
	if err != nil {
		return ownerID, repoID, err
	}
	repoID = r.ID
	return ownerID, repoID, nil
}
