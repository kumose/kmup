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

	"github.com/kumose/kmup/modules/translation"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestViewBranches(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	req := NewRequest(t, "GET", "/user2/repo1/branches")
	resp := MakeRequest(t, req, http.StatusOK)

	htmlDoc := NewHTMLParser(t, resp.Body)
	_, exists := htmlDoc.doc.Find(".delete-branch-button").Attr("data-url")
	assert.False(t, exists, "The template has changed")
}

func TestDeleteBranch(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	deleteBranch(t)
}

func TestUndoDeleteBranch(t *testing.T) {
	onKmupRun(t, func(t *testing.T, u *url.URL) {
		deleteBranch(t)
		htmlDoc, name := branchAction(t, ".restore-branch-button")
		assert.Contains(t,
			htmlDoc.doc.Find(".ui.positive.message").Text(),
			translation.NewLocale("en-US").TrString("repo.branch.restore_success", name),
		)
	})
}

func deleteBranch(t *testing.T) {
	htmlDoc, name := branchAction(t, ".delete-branch-button")
	assert.Contains(t,
		htmlDoc.doc.Find(".ui.positive.message").Text(),
		translation.NewLocale("en-US").TrString("repo.branch.deletion_success", name),
	)
}

func branchAction(t *testing.T, button string) (*HTMLDoc, string) {
	session := loginUser(t, "user2")
	req := NewRequest(t, "GET", "/user2/repo1/branches")
	resp := session.MakeRequest(t, req, http.StatusOK)

	htmlDoc := NewHTMLParser(t, resp.Body)
	link, exists := htmlDoc.doc.Find(button).Attr("data-url")
	if !assert.True(t, exists, "The template has changed") {
		t.Skip()
	}

	req = NewRequestWithValues(t, "POST", link, map[string]string{
		"_csrf": htmlDoc.GetCSRF(),
	})
	session.MakeRequest(t, req, http.StatusOK)

	url, err := url.Parse(link)
	assert.NoError(t, err)
	req = NewRequest(t, "GET", "/user2/repo1/branches")
	resp = session.MakeRequest(t, req, http.StatusOK)

	return NewHTMLParser(t, resp.Body), url.Query().Get("name")
}
