// Copyright 2015 The Gogs Authors. All rights reserved.
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

//go:build !gogit

package git

import (
	"fmt"
	"time"

	"github.com/kumose/kmup/modules/util"
)

// Signature represents the Author, Committer or Tagger information.
type Signature struct {
	Name  string    // the committer name, it can be anything
	Email string    // the committer email, it can be anything
	When  time.Time // the timestamp of the signature
}

func (s *Signature) String() string {
	return fmt.Sprintf("%s <%s>", s.Name, s.Email)
}

// Decode decodes a byte array representing a signature to signature
func (s *Signature) Decode(b []byte) {
	*s = *parseSignatureFromCommitLine(util.UnsafeBytesToString(b))
}
