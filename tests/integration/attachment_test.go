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
	"bytes"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	repo_model "github.com/kumose/kmup/models/repo"
	"github.com/kumose/kmup/modules/storage"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/tests"

	"github.com/stretchr/testify/assert"
)

func testGeneratePngBytes() []byte {
	myImage := image.NewRGBA(image.Rect(0, 0, 32, 32))
	var buff bytes.Buffer
	_ = png.Encode(&buff, myImage)
	return buff.Bytes()
}

func testCreateIssueAttachment(t *testing.T, session *TestSession, csrf, repoURL, filename string, content []byte, expectedStatus int) string {
	body := &bytes.Buffer{}

	// Setup multi-part
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	assert.NoError(t, err)
	_, err = part.Write(content)
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	req := NewRequestWithBody(t, "POST", repoURL+"/issues/attachments", body)
	req.Header.Add("X-Csrf-Token", csrf)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp := session.MakeRequest(t, req, expectedStatus)

	if expectedStatus != http.StatusOK {
		return ""
	}
	var obj map[string]string
	DecodeJSON(t, resp, &obj)
	return obj["uuid"]
}

func TestCreateAnonymousAttachment(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	session := emptyTestSession(t)
	testCreateIssueAttachment(t, session, GetAnonymousCSRFToken(t, session), "user2/repo1", "image.png", testGeneratePngBytes(), http.StatusSeeOther)
}

func TestCreateIssueAttachment(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	const repoURL = "user2/repo1"
	session := loginUser(t, "user2")
	uuid := testCreateIssueAttachment(t, session, GetUserCSRFToken(t, session), repoURL, "image.png", testGeneratePngBytes(), http.StatusOK)

	req := NewRequest(t, "GET", repoURL+"/issues/new")
	resp := session.MakeRequest(t, req, http.StatusOK)
	htmlDoc := NewHTMLParser(t, resp.Body)

	link, exists := htmlDoc.doc.Find("form#new-issue").Attr("action")
	assert.True(t, exists, "The template has changed")

	postData := map[string]string{
		"_csrf":   htmlDoc.GetCSRF(),
		"title":   "New Issue With Attachment",
		"content": "some content",
		"files":   uuid,
	}

	req = NewRequestWithValues(t, "POST", link, postData)
	resp = session.MakeRequest(t, req, http.StatusOK)
	test.RedirectURL(resp) // check that redirect URL exists

	// Validate that attachment is available
	req = NewRequest(t, "GET", "/attachments/"+uuid)
	session.MakeRequest(t, req, http.StatusOK)

	// anonymous visit should be allowed because user2/repo1 is a public repository
	MakeRequest(t, req, http.StatusOK)
}

func TestGetAttachment(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	adminSession := loginUser(t, "user1")
	user2Session := loginUser(t, "user2")
	user8Session := loginUser(t, "user8")
	emptySession := emptyTestSession(t)
	testCases := []struct {
		name       string
		uuid       string
		createFile bool
		session    *TestSession
		want       int
	}{
		{"LinkedIssueUUID", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", true, user2Session, http.StatusOK},
		{"LinkedCommentUUID", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a17", true, user2Session, http.StatusOK},
		{"linked_release_uuid", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a19", true, user2Session, http.StatusOK},
		{"NotExistingUUID", "b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a18", false, user2Session, http.StatusNotFound},
		{"FileMissing", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a18", false, user2Session, http.StatusInternalServerError},
		{"NotLinked", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a20", true, user2Session, http.StatusNotFound},
		{"NotLinkedAccessibleByUploader", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a20", true, user8Session, http.StatusOK},
		{"PublicByNonLogged", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", true, emptySession, http.StatusOK},
		{"PrivateByNonLogged", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", true, emptySession, http.StatusNotFound},
		{"PrivateAccessibleByAdmin", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", true, adminSession, http.StatusOK},
		{"PrivateAccessibleByUser", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", true, user2Session, http.StatusOK},
		{"RepoNotAccessibleByUser", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12", true, user8Session, http.StatusNotFound},
		{"OrgNotAccessibleByUser", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a21", true, user8Session, http.StatusNotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write empty file to be available for response
			if tc.createFile {
				_, err := storage.Attachments.Save(repo_model.AttachmentRelativePath(tc.uuid), strings.NewReader("hello world"), -1)
				assert.NoError(t, err)
			}
			// Actual test
			req := NewRequest(t, "GET", "/attachments/"+tc.uuid)
			tc.session.MakeRequest(t, req, tc.want)
		})
	}
}
