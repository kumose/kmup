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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadMailerFrom(t *testing.T) {
	kases := map[string]*Mailer{
		"smtp.mydomain.com": {
			SMTPAddr: "smtp.mydomain.com",
			SMTPPort: "465",
		},
		"smtp.mydomain.com:123": {
			SMTPAddr: "smtp.mydomain.com",
			SMTPPort: "123",
		},
		":123": {
			SMTPAddr: "127.0.0.1",
			SMTPPort: "123",
		},
	}
	for host, kase := range kases {
		t.Run(host, func(t *testing.T) {
			cfg, _ := NewConfigProviderFromData("")
			sec := cfg.Section("mailer")
			sec.NewKey("ENABLED", "true")
			sec.NewKey("HOST", host)

			// Check mailer setting
			loadMailerFrom(cfg)

			assert.Equal(t, kase.SMTPAddr, MailService.SMTPAddr)
			assert.Equal(t, kase.SMTPPort, MailService.SMTPPort)
		})
	}
}
