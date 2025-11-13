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

package externalaccount

import (
	"context"
	"strconv"
	"strings"

	issues_model "github.com/kumose/kmup/models/issues"
	repo_model "github.com/kumose/kmup/models/repo"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/structs"

	"github.com/markbates/goth"
)

func toExternalLoginUser(authSourceID int64, user *user_model.User, gothUser goth.User) *user_model.ExternalLoginUser {
	return &user_model.ExternalLoginUser{
		ExternalID:        gothUser.UserID,
		UserID:            user.ID,
		LoginSourceID:     authSourceID,
		RawData:           gothUser.RawData,
		Provider:          gothUser.Provider,
		Email:             gothUser.Email,
		Name:              gothUser.Name,
		FirstName:         gothUser.FirstName,
		LastName:          gothUser.LastName,
		NickName:          gothUser.NickName,
		Description:       gothUser.Description,
		AvatarURL:         gothUser.AvatarURL,
		Location:          gothUser.Location,
		AccessToken:       gothUser.AccessToken,
		AccessTokenSecret: gothUser.AccessTokenSecret,
		RefreshToken:      gothUser.RefreshToken,
		ExpiresAt:         gothUser.ExpiresAt,
	}
}

// LinkAccountToUser link the gothUser to the user
func LinkAccountToUser(ctx context.Context, authSourceID int64, user *user_model.User, gothUser goth.User) error {
	externalLoginUser := toExternalLoginUser(authSourceID, user, gothUser)

	if err := user_model.LinkExternalToUser(ctx, user, externalLoginUser); err != nil {
		return err
	}

	externalID := externalLoginUser.ExternalID

	var tp structs.GitServiceType
	for _, s := range structs.SupportedFullGitService {
		if strings.EqualFold(s.Name(), gothUser.Provider) {
			tp = s
			break
		}
	}

	if tp.Name() != "" {
		return UpdateMigrationsByType(ctx, tp, externalID, user.ID)
	}

	return nil
}

// EnsureLinkExternalToUser link the gothUser to the user
func EnsureLinkExternalToUser(ctx context.Context, authSourceID int64, user *user_model.User, gothUser goth.User) error {
	externalLoginUser := toExternalLoginUser(authSourceID, user, gothUser)
	return user_model.EnsureLinkExternalToUser(ctx, externalLoginUser)
}

// UpdateMigrationsByType updates all migrated repositories' posterid from gitServiceType to replace originalAuthorID to posterID
func UpdateMigrationsByType(ctx context.Context, tp structs.GitServiceType, externalUserID string, userID int64) error {
	// Skip update if externalUserID is not a valid numeric ID or exceeds int64
	if _, err := strconv.ParseInt(externalUserID, 10, 64); err != nil {
		return nil
	}

	if err := issues_model.UpdateIssuesMigrationsByType(ctx, tp, externalUserID, userID); err != nil {
		return err
	}

	if err := issues_model.UpdateCommentsMigrationsByType(ctx, tp, externalUserID, userID); err != nil {
		return err
	}

	if err := repo_model.UpdateReleasesMigrationsByType(ctx, tp, externalUserID, userID); err != nil {
		return err
	}

	if err := issues_model.UpdateReactionsMigrationsByType(ctx, tp, externalUserID, userID); err != nil {
		return err
	}
	return issues_model.UpdateReviewsMigrationsByType(ctx, tp, externalUserID, userID)
}
