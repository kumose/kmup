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

package incoming

import (
	"strings"
	"testing"

	"github.com/jhillyerd/enmime"
	"github.com/stretchr/testify/assert"
)

func TestIsAutomaticReply(t *testing.T) {
	cases := []struct {
		Headers  map[string]string
		Expected bool
	}{
		{
			Headers:  map[string]string{},
			Expected: false,
		},
		{
			Headers: map[string]string{
				"Auto-Submitted": "no",
			},
			Expected: false,
		},
		{
			Headers: map[string]string{
				"Auto-Submitted": "yes",
			},
			Expected: true,
		},
		{
			Headers: map[string]string{
				"X-Autoreply": "no",
			},
			Expected: false,
		},
		{
			Headers: map[string]string{
				"X-Autoreply": "yes",
			},
			Expected: true,
		},
		{
			Headers: map[string]string{
				"X-Autorespond": "yes",
			},
			Expected: true,
		},
	}

	for _, c := range cases {
		b := enmime.Builder().
			From("Dummy", "dummy@kmup.io").
			To("Dummy", "dummy@kmup.io")
		for k, v := range c.Headers {
			b = b.Header(k, v)
		}
		root, err := b.Build()
		assert.NoError(t, err)
		env, err := enmime.EnvelopeFromPart(root)
		assert.NoError(t, err)

		assert.Equal(t, c.Expected, isAutomaticReply(env))
	}
}

func TestGetContentFromMailReader(t *testing.T) {
	mailString := "Content-Type: multipart/mixed; boundary=message-boundary\r\n" +
		"\r\n" +
		"--message-boundary\r\n" +
		"Content-Type: multipart/alternative; boundary=text-boundary\r\n" +
		"\r\n" +
		"--text-boundary\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Disposition: inline\r\n" +
		"\r\n" +
		"mail content\r\n" +
		"--text-boundary--\r\n" +
		"--message-boundary\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Disposition: attachment; filename=attachment.txt\r\n" +
		"\r\n" +
		"attachment content\r\n" +
		"--message-boundary--\r\n"

	env, err := enmime.ReadEnvelope(strings.NewReader(mailString))
	assert.NoError(t, err)
	content := getContentFromMailReader(env)
	assert.Equal(t, "mail content", content.Content)
	assert.Len(t, content.Attachments, 1)
	assert.Equal(t, "attachment.txt", content.Attachments[0].Name)
	assert.Equal(t, []byte("attachment content"), content.Attachments[0].Content)

	mailString = "Content-Type: multipart/mixed; boundary=message-boundary\r\n" +
		"\r\n" +
		"--message-boundary\r\n" +
		"Content-Type: multipart/alternative; boundary=text-boundary\r\n" +
		"\r\n" +
		"--text-boundary\r\n" +
		"Content-Type: text/html\r\n" +
		"Content-Disposition: inline\r\n" +
		"\r\n" +
		"<p>mail content</p>\r\n" +
		"--text-boundary--\r\n" +
		"--message-boundary--\r\n"

	env, err = enmime.ReadEnvelope(strings.NewReader(mailString))
	assert.NoError(t, err)
	content = getContentFromMailReader(env)
	assert.Equal(t, "mail content", content.Content)
	assert.Empty(t, content.Attachments)

	mailString = "Content-Type: multipart/mixed; boundary=message-boundary\r\n" +
		"\r\n" +
		"--message-boundary\r\n" +
		"Content-Type: multipart/alternative; boundary=text-boundary\r\n" +
		"\r\n" +
		"--text-boundary\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Disposition: inline\r\n" +
		"\r\n" +
		"mail content without signature\r\n" +
		"--\r\n" +
		"signature\r\n" +
		"--text-boundary--\r\n" +
		"--message-boundary--\r\n"

	env, err = enmime.ReadEnvelope(strings.NewReader(mailString))
	assert.NoError(t, err)
	content = getContentFromMailReader(env)
	assert.NoError(t, err)
	assert.Equal(t, "mail content without signature", content.Content)
	assert.Empty(t, content.Attachments)
}
