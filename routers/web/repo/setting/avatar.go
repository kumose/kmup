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

package setting

import (
	"errors"
	"fmt"
	"io"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/typesniffer"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/services/context"
	"github.com/kumose/kmup/services/forms"
	repo_service "github.com/kumose/kmup/services/repository"
)

// UpdateAvatarSetting update repo's avatar
func UpdateAvatarSetting(ctx *context.Context, form forms.AvatarForm) error {
	ctxRepo := ctx.Repo.Repository

	if form.Avatar == nil {
		// No avatar is uploaded and we not removing it here.
		// No random avatar generated here.
		// Just exit, no action.
		if ctxRepo.CustomAvatarRelativePath() == "" {
			log.Trace("No avatar was uploaded for repo: %d. Default icon will appear instead.", ctxRepo.ID)
		}
		return nil
	}

	r, err := form.Avatar.Open()
	if err != nil {
		return fmt.Errorf("Avatar.Open: %w", err)
	}
	defer r.Close()

	if form.Avatar.Size > setting.Avatar.MaxFileSize {
		return errors.New(ctx.Locale.TrString("settings.uploaded_avatar_is_too_big", form.Avatar.Size/1024, setting.Avatar.MaxFileSize/1024))
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}
	st := typesniffer.DetectContentType(data)
	if !(st.IsImage() && !st.IsSvgImage()) {
		return errors.New(ctx.Locale.TrString("settings.uploaded_avatar_not_a_image"))
	}
	if err = repo_service.UploadAvatar(ctx, ctxRepo, data); err != nil {
		return fmt.Errorf("UploadAvatar: %w", err)
	}
	return nil
}

// SettingsAvatar save new POSTed repository avatar
func SettingsAvatar(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.AvatarForm)
	form.Source = forms.AvatarLocal
	if err := UpdateAvatarSetting(ctx, *form); err != nil {
		ctx.Flash.Error(err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("repo.settings.update_avatar_success"))
	}
	ctx.Redirect(ctx.Repo.RepoLink + "/settings")
}

// SettingsDeleteAvatar delete repository avatar
func SettingsDeleteAvatar(ctx *context.Context) {
	if err := repo_service.DeleteAvatar(ctx, ctx.Repo.Repository); err != nil {
		ctx.Flash.Error(fmt.Sprintf("DeleteAvatar: %v", err))
	}
	ctx.JSONRedirect(ctx.Repo.RepoLink + "/settings")
}
