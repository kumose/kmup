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

package git

import (
	"strconv"
	"strings"
	"time"

	"github.com/kumose/kmup/modules/log"
)

// Helper to get a signature from the commit line, which looks like:
//
//	full name <user@example.com> 1378823654 +0200
//
// Haven't found the official reference for the standard format yet.
// This function never fails, if the "line" can't be parsed, it returns a default Signature with "zero" time.
func parseSignatureFromCommitLine(line string) *Signature {
	sig := &Signature{}
	s1, sx, ok1 := strings.Cut(line, " <")
	s2, s3, ok2 := strings.Cut(sx, "> ")
	if !ok1 || !ok2 {
		sig.Name = line
		return sig
	}
	sig.Name, sig.Email = s1, s2

	if strings.Count(s3, " ") == 1 {
		ts, tz, _ := strings.Cut(s3, " ")
		seconds, _ := strconv.ParseInt(ts, 10, 64)
		if tzTime, err := time.Parse("-0700", tz); err == nil {
			sig.When = time.Unix(seconds, 0).In(tzTime.Location())
		}
	} else {
		// the old kmup code tried to parse the date in a few different formats, but it's not clear why.
		// according to public document, only the standard format "timestamp timezone" could be found, so drop other formats.
		log.Error("suspicious commit line format: %q", line)
		for _, fmt := range []string{ /*"Mon Jan _2 15:04:05 2006 -0700"*/ } {
			if t, err := time.Parse(fmt, s3); err == nil {
				sig.When = t
				break
			}
		}
	}
	return sig
}
