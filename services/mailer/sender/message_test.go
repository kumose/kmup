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

package sender

import (
	"strings"
	"testing"
	"time"

	"github.com/kumose/kmup/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMessageID(t *testing.T) {
	mailService := setting.Mailer{
		From: "test@kmup.com",
	}

	setting.MailService = &mailService
	setting.Domain = "localhost"

	date := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
	m := NewMessageFrom("", "display-name", "from-address", "subject", "body")
	m.Date = date
	gm := m.ToMessage()
	assert.Equal(t, "<autogen-946782245000-41e8fc54a8ad3a3f@localhost>", gm.GetGenHeader("Message-ID")[0])

	m = NewMessageFrom("a@b.com", "display-name", "from-address", "subject", "body")
	m.Date = date
	gm = m.ToMessage()
	assert.Equal(t, "<autogen-946782245000-cc88ce3cfe9bd04f@localhost>", gm.GetGenHeader("Message-ID")[0])

	m = NewMessageFrom("a@b.com", "display-name", "from-address", "subject", "body")
	m.SetHeader("Message-ID", "<msg-d@domain.com>")
	gm = m.ToMessage()
	assert.Equal(t, "<msg-d@domain.com>", gm.GetGenHeader("Message-ID")[0])
}

func TestToMessage(t *testing.T) {
	oldConf := setting.MailService
	defer func() {
		setting.MailService = oldConf
	}()
	setting.MailService = &setting.Mailer{
		From: "test@kmup.com",
	}

	m1 := Message{
		Info:            "info",
		FromAddress:     "test@kmup.com",
		FromDisplayName: "Test Kmup",
		To:              "a@b.com",
		Subject:         "Issue X Closed",
		Body:            "Some Issue got closed by Y-Man",
	}

	assertHeaders := func(t *testing.T, expected, header map[string]string) {
		for k, v := range expected {
			assert.Equal(t, v, header[k], "Header %s should be %s but got %s", k, v, header[k])
		}
	}

	buf := &strings.Builder{}
	_, err := m1.ToMessage().WriteTo(buf)
	assert.NoError(t, err)
	header, _ := extractMailHeaderAndContent(t, buf.String())
	assertHeaders(t, map[string]string{
		"Content-Type":             "multipart/alternative;",
		"Date":                     "Mon, 01 Jan 0001 00:00:00 +0000",
		"From":                     "\"Test Kmup\" <test@kmup.com>",
		"Message-ID":               "<autogen--6795364578871-69c000786adc60dc@localhost>",
		"MIME-Version":             "1.0",
		"Subject":                  "Issue X Closed",
		"To":                       "<a@b.com>",
		"X-Auto-Response-Suppress": "All",
	}, header)

	setting.MailService.OverrideHeader = map[string][]string{
		"Message-ID":     {""},               // delete message id
		"Auto-Submitted": {"auto-generated"}, // suppress auto replay
	}

	buf = &strings.Builder{}
	_, err = m1.ToMessage().WriteTo(buf)
	assert.NoError(t, err)
	header, _ = extractMailHeaderAndContent(t, buf.String())
	assertHeaders(t, map[string]string{
		"Content-Type":             "multipart/alternative;",
		"Date":                     "Mon, 01 Jan 0001 00:00:00 +0000",
		"From":                     "\"Test Kmup\" <test@kmup.com>",
		"Message-ID":               "",
		"MIME-Version":             "1.0",
		"Subject":                  "Issue X Closed",
		"To":                       "<a@b.com>",
		"X-Auto-Response-Suppress": "All",
		"Auto-Submitted":           "auto-generated",
	}, header)
}

func extractMailHeaderAndContent(t *testing.T, mail string) (map[string]string, string) {
	header := make(map[string]string)

	parts := strings.SplitN(mail, "boundary=", 2)
	if !assert.Len(t, parts, 2) {
		return nil, ""
	}
	content := strings.TrimSpace("boundary=" + parts[1])

	hParts := strings.SplitSeq(parts[0], "\n")

	for hPart := range hParts {
		parts := strings.SplitN(hPart, ":", 2)
		hk := strings.TrimSpace(parts[0])
		if hk != "" {
			header[hk] = strings.TrimSpace(parts[1])
		}
	}

	return header, content
}
