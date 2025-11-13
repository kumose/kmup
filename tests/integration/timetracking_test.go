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
	"path"
	"testing"
	"time"

	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func TestViewTimetrackingControls(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	t.Run("Exist", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		session := loginUser(t, "user2")
		testViewTimetrackingControls(t, session, "user2", "repo1", "1", true)
	})

	t.Run("Non-exist", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		session := loginUser(t, "user5")
		testViewTimetrackingControls(t, session, "user2", "repo1", "1", false)
	})

	t.Run("Disabled", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		session := loginUser(t, "user2")
		testViewTimetrackingControls(t, session, "org3", "repo3", "1", false)
	})
}

func testViewTimetrackingControls(t *testing.T, session *TestSession, user, repo, issue string, canTrackTime bool) {
	req := NewRequest(t, "GET", path.Join(user, repo, "issues", issue))
	resp := session.MakeRequest(t, req, http.StatusOK)

	htmlDoc := NewHTMLParser(t, resp.Body)

	AssertHTMLElement(t, htmlDoc, ".issue-start-time", canTrackTime)
	AssertHTMLElement(t, htmlDoc, ".issue-add-time", canTrackTime)

	issueLink := path.Join(user, repo, "issues", issue)
	reqStart := NewRequestWithValues(t, "POST", path.Join(issueLink, "times", "stopwatch", "start"), map[string]string{
		"_csrf": htmlDoc.GetCSRF(),
	})
	if canTrackTime {
		session.MakeRequest(t, reqStart, http.StatusOK)

		req = NewRequest(t, "GET", issueLink)
		resp = session.MakeRequest(t, req, http.StatusOK)
		htmlDoc = NewHTMLParser(t, resp.Body)

		events := htmlDoc.doc.Find(".event > .comment-text-line")
		assert.Contains(t, events.Last().Text(), "started working")

		AssertHTMLElement(t, htmlDoc, ".issue-stop-time", true)
		AssertHTMLElement(t, htmlDoc, ".issue-cancel-time", true)

		// Sleep for 1 second to not get wrong order for stopping timer
		time.Sleep(time.Second)

		reqStop := NewRequestWithValues(t, "POST", path.Join(issueLink, "times", "stopwatch", "stop"), map[string]string{
			"_csrf": htmlDoc.GetCSRF(),
		})
		session.MakeRequest(t, reqStop, http.StatusOK)

		req = NewRequest(t, "GET", issueLink)
		resp = session.MakeRequest(t, req, http.StatusOK)
		htmlDoc = NewHTMLParser(t, resp.Body)

		events = htmlDoc.doc.Find(".event > .comment-text-line")
		assert.Contains(t, events.Last().Text(), "worked for ")
	} else {
		session.MakeRequest(t, reqStart, http.StatusNotFound)
	}
}
