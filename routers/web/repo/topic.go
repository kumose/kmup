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

package repo

import (
	"net/http"
	"strings"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/context"
)

// TopicsPost response for creating repository
func TopicsPost(ctx *context.Context) {
	if ctx.Doer == nil {
		ctx.JSON(http.StatusForbidden, map[string]any{
			"message": "Only owners could change the topics.",
		})
		return
	}

	topics := make([]string, 0)
	topicsStr := ctx.FormTrim("topics")
	if len(topicsStr) > 0 {
		topics = strings.Split(topicsStr, ",")
	}

	validTopics, invalidTopics := repo_model.SanitizeAndValidateTopics(topics)

	if len(validTopics) > 25 {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]any{
			"invalidTopics": nil,
			"message":       ctx.Tr("repo.topic.count_prompt"),
		})
		return
	}

	if len(invalidTopics) > 0 {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]any{
			"invalidTopics": invalidTopics,
			"message":       ctx.Tr("repo.topic.format_prompt"),
		})
		return
	}

	err := repo_model.SaveTopics(ctx, ctx.Repo.Repository.ID, validTopics...)
	if err != nil {
		log.Error("SaveTopics failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, map[string]any{
			"message": "Save topics failed.",
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}
