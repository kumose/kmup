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

package secrets

import (
	"regexp"
	"strings"
	"sync"

	"github.com/kumose/kmup/modules/util"
)

// https://docs.github.com/en/actions/learn-github-actions/variables#naming-conventions-for-configuration-variables
// https://docs.github.com/en/actions/security-guides/encrypted-secrets#naming-your-secrets
var globalVars = sync.OnceValue(func() (ret struct {
	namePattern, forbiddenPrefixPattern *regexp.Regexp
},
) {
	ret.namePattern = regexp.MustCompile("(?i)^[A-Z_][A-Z0-9_]*$")
	ret.forbiddenPrefixPattern = regexp.MustCompile("(?i)^GIT(EA|HUB)_")
	return ret
})

func ValidateName(name string) error {
	vars := globalVars()
	if !vars.namePattern.MatchString(name) ||
		vars.forbiddenPrefixPattern.MatchString(name) ||
		strings.EqualFold(name, "CI") /* CI is always set to true in GitHub Actions*/ {
		return util.NewInvalidArgumentErrorf("invalid variable or secret name")
	}
	return nil
}
