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
	"bytes"
	"net/http"
	"strconv"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/services/context"
)

type languageResponse []*repo_model.LanguageStat

func (l languageResponse) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.WriteString("{"); err != nil {
		return nil, err
	}
	for i, lang := range l {
		if i > 0 {
			if _, err := buf.WriteString(","); err != nil {
				return nil, err
			}
		}
		if _, err := buf.WriteString(strconv.Quote(lang.Language)); err != nil {
			return nil, err
		}
		if _, err := buf.WriteString(":"); err != nil {
			return nil, err
		}
		if _, err := buf.WriteString(strconv.FormatInt(lang.Size, 10)); err != nil {
			return nil, err
		}
	}
	if _, err := buf.WriteString("}"); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GetLanguages returns languages and number of bytes of code written
func GetLanguages(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/languages repository repoGetLanguages
	// ---
	// summary: Get languages and number of bytes of code written
	// produces:
	//   - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// responses:
	//   "404":
	//     "$ref": "#/responses/notFound"
	//   "200":
	//     "$ref": "#/responses/LanguageStatistics"

	langs, err := repo_model.GetLanguageStats(ctx, ctx.Repo.Repository)
	if err != nil {
		log.Error("GetLanguageStats failed: %v", err)
		ctx.APIErrorInternal(err)
		return
	}

	resp := make(languageResponse, len(langs))
	copy(resp, langs)

	ctx.JSON(http.StatusOK, resp)
}
