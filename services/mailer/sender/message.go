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
	"fmt"
	"hash/fnv"
	"net/mail"
	"strings"
	"time"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"

	"github.com/jaytaylor/html2text"
	gomail "github.com/wneessen/go-mail"
)

// Message mail body and log info
type Message struct {
	Info            string // Message information for log purpose.
	FromAddress     string
	FromDisplayName string
	To              string // Use only one recipient to prevent leaking of addresses
	ReplyTo         string
	Subject         string
	Date            time.Time
	Body            string
	Headers         map[string][]string
}

// ToMessage converts a Message to gomail.Message
func (m *Message) ToMessage() *gomail.Msg {
	msg := gomail.NewMsg()
	addr := mail.Address{Name: m.FromDisplayName, Address: m.FromAddress}
	_ = msg.SetAddrHeader("From", addr.String())
	_ = msg.SetAddrHeader("To", m.To)
	if m.ReplyTo != "" {
		msg.SetGenHeader("Reply-To", m.ReplyTo)
	}
	for header := range m.Headers {
		msg.SetGenHeader(gomail.Header(header), m.Headers[header]...)
	}

	if setting.MailService.SubjectPrefix != "" {
		msg.SetGenHeader("Subject", setting.MailService.SubjectPrefix+" "+m.Subject)
	} else {
		msg.SetGenHeader("Subject", m.Subject)
	}
	msg.SetDateWithValue(m.Date)
	msg.SetGenHeader("X-Auto-Response-Suppress", "All")

	plainBody, err := html2text.FromString(m.Body)
	if err != nil || setting.MailService.SendAsPlainText {
		if strings.Contains(util.TruncateRunes(m.Body, 100), "<html>") {
			log.Warn("Mail contains HTML but configured to send as plain text.")
		}
		msg.SetBodyString("text/plain", plainBody)
	} else {
		msg.SetBodyString("text/plain", plainBody)
		msg.AddAlternativeString("text/html", m.Body)
	}

	if len(msg.GetGenHeader("Message-ID")) == 0 {
		msg.SetGenHeader("Message-ID", m.generateAutoMessageID())
	}

	for k, v := range setting.MailService.OverrideHeader {
		if len(msg.GetGenHeader(gomail.Header(k))) != 0 {
			log.Debug("Mailer override header '%s' as per config", k)
		}
		msg.SetGenHeader(gomail.Header(k), v...)
	}

	return msg
}

// SetHeader adds additional headers to a message
func (m *Message) SetHeader(field string, value ...string) {
	m.Headers[field] = value
}

func (m *Message) generateAutoMessageID() string {
	dateMs := m.Date.UnixNano() / 1e6
	h := fnv.New64()
	if len(m.To) > 0 {
		_, _ = h.Write([]byte(m.To))
	}
	_, _ = h.Write([]byte(m.Subject))
	_, _ = h.Write([]byte(m.Body))
	return fmt.Sprintf("<autogen-%d-%016x@%s>", dateMs, h.Sum64(), setting.Domain)
}

// NewMessageFrom creates new mail message object with custom From header.
func NewMessageFrom(to, fromDisplayName, fromAddress, subject, body string) *Message {
	log.Trace("NewMessageFrom (body):\n%s", body)

	return &Message{
		FromAddress:     fromAddress,
		FromDisplayName: fromDisplayName,
		To:              to,
		Subject:         subject,
		Date:            time.Now(),
		Body:            body,
		Headers:         map[string][]string{},
	}
}

// NewMessage creates new mail message object with default From header.
func NewMessage(to, subject, body string) *Message {
	return NewMessageFrom(to, setting.MailService.FromName, setting.MailService.FromEmail, subject, body)
}
