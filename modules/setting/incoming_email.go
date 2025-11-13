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

package setting

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/kumose/kmup/modules/log"
)

var IncomingEmail = struct {
	Enabled              bool
	ReplyToAddress       string
	TokenPlaceholder     string `ini:"-"`
	Host                 string
	Port                 int
	UseTLS               bool `ini:"USE_TLS"`
	SkipTLSVerify        bool `ini:"SKIP_TLS_VERIFY"`
	Username             string
	Password             string
	Mailbox              string
	DeleteHandledMessage bool
	MaximumMessageSize   uint32
}{
	Mailbox:              "INBOX",
	DeleteHandledMessage: true,
	TokenPlaceholder:     "%{token}",
	MaximumMessageSize:   10485760,
}

func loadIncomingEmailFrom(rootCfg ConfigProvider) {
	mustMapSetting(rootCfg, "email.incoming", &IncomingEmail)

	if !IncomingEmail.Enabled {
		return
	}

	if err := checkReplyToAddress(); err != nil {
		log.Fatal("Invalid incoming_mail.REPLY_TO_ADDRESS (%s): %v", IncomingEmail.ReplyToAddress, err)
	}
}

func checkReplyToAddress() error {
	parsed, err := mail.ParseAddress(IncomingEmail.ReplyToAddress)
	if err != nil {
		return err
	}

	if parsed.Name != "" {
		return errors.New("name must not be set")
	}

	c := strings.Count(IncomingEmail.ReplyToAddress, IncomingEmail.TokenPlaceholder)
	switch c {
	case 0:
		return fmt.Errorf("%s must appear in the user part of the address (before the @)", IncomingEmail.TokenPlaceholder)
	case 1:
	default:
		return fmt.Errorf("%s must appear only once", IncomingEmail.TokenPlaceholder)
	}

	parts := strings.Split(IncomingEmail.ReplyToAddress, "@")
	if !strings.Contains(parts[0], IncomingEmail.TokenPlaceholder) {
		return fmt.Errorf("%s must appear in the user part of the address (before the @)", IncomingEmail.TokenPlaceholder)
	}

	return nil
}
