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

package integration

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	repo_module "github.com/kumose/kmup/modules/repository"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIListLabelTemplates(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	req := NewRequest(t, "GET", "/api/v1/label/templates")
	resp := MakeRequest(t, req, http.StatusOK)

	var templateList []string
	DecodeJSON(t, resp, &templateList)

	for i := range repo_module.LabelTemplateFiles {
		assert.Equal(t, repo_module.LabelTemplateFiles[i].DisplayName, templateList[i])
	}
}

func TestAPIGetLabelTemplateInfo(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// If Kmup has for some reason no Label templates, we need to skip this test
	if len(repo_module.LabelTemplateFiles) == 0 {
		return
	}

	// Use the first template for the test
	templateName := repo_module.LabelTemplateFiles[0].DisplayName

	urlStr := "/api/v1/label/templates/" + url.PathEscape(templateName)
	req := NewRequest(t, "GET", urlStr)
	resp := MakeRequest(t, req, http.StatusOK)

	var templateInfo []api.LabelTemplate
	DecodeJSON(t, resp, &templateInfo)

	labels, err := repo_module.LoadTemplateLabelsByDisplayName(templateName)
	assert.NoError(t, err)

	for i := range labels {
		assert.Equal(t, strings.TrimLeft(labels[i].Color, "#"), templateInfo[i].Color)
		assert.Equal(t, labels[i].Description, templateInfo[i].Description)
		assert.Equal(t, labels[i].Exclusive, templateInfo[i].Exclusive)
		assert.Equal(t, labels[i].Name, templateInfo[i].Name)
	}
}
