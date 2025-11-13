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
	"testing"

	"github.com/kumose/kmup/modules/options"
	repo_module "github.com/kumose/kmup/modules/repository"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIListLicenseTemplates(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	req := NewRequest(t, "GET", "/api/v1/licenses")
	resp := MakeRequest(t, req, http.StatusOK)

	// This tests if the API returns a list of strings
	var licenseList []api.LicensesTemplateListEntry
	DecodeJSON(t, resp, &licenseList)
}

func TestAPIGetLicenseTemplateInfo(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// If Kmup has for some reason no License templates, we need to skip this test
	if len(repo_module.Licenses) == 0 {
		return
	}

	// Use the first template for the test
	licenseName := repo_module.Licenses[0]

	urlStr := "/api/v1/licenses/" + url.PathEscape(licenseName)
	req := NewRequest(t, "GET", urlStr)
	resp := MakeRequest(t, req, http.StatusOK)

	var licenseInfo api.LicenseTemplateInfo
	DecodeJSON(t, resp, &licenseInfo)

	// We get the text of the template here
	text, _ := options.License(licenseName)

	assert.Equal(t, licenseInfo.Key, licenseName)
	assert.Equal(t, licenseInfo.Name, licenseName)
	assert.Equal(t, licenseInfo.Body, string(text))
}
