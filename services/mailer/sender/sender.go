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
	"io"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
)

type Sender interface {
	Send(from string, to []string, msg io.WriterTo) error
}

var Send = send

func send(sender Sender, msgs ...*Message) error {
	if setting.MailService == nil {
		log.Error("Mailer: Send is being invoked but mail service hasn't been initialized")
		return nil
	}
	for _, msg := range msgs {
		m := msg.ToMessage()
		froms := m.GetFrom()
		to, err := m.GetRecipients()
		if err != nil {
			return err
		}

		// TODO: implement sending from multiple addresses
		if err := sender.Send(froms[0].Address, to, m); err != nil {
			return err
		}
	}
	return nil
}
