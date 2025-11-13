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
	"bytes"
	"io"

	"github.com/kumose/kmup/modules/log"
)

// DummySender Sender sendmail mail sender
type DummySender struct{}

var _ Sender = &DummySender{}

// Send send email
func (s *DummySender) Send(from string, to []string, msg io.WriterTo) error {
	buf := bytes.Buffer{}
	if _, err := msg.WriteTo(&buf); err != nil {
		return err
	}
	log.Debug("Mail From: %s To: %v Body: %s", from, to, buf.String())
	return nil
}
