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

package markup

import (
	"regexp"
	"strings"

	"github.com/kumose/kmup/modules/setting"

	"github.com/microcosm-cc/bluemonday"
)

func (st *Sanitizer) addSanitizerRules(policy *bluemonday.Policy, rules []setting.MarkupSanitizerRule) {
	for _, rule := range rules {
		if rule.AllowDataURIImages {
			policy.AllowDataURIImages()
		}
		if rule.Element != "" {
			if rule.Regexp != "" {
				if !strings.HasPrefix(rule.Regexp, "^") || !strings.HasSuffix(rule.Regexp, "$") {
					panic("Markup sanitizer rule regexp must start with ^ and end with $ to be strict")
				}
				policy.AllowAttrs(rule.AllowAttr).Matching(regexp.MustCompile(rule.Regexp)).OnElements(rule.Element)
			} else {
				policy.AllowAttrs(rule.AllowAttr).OnElements(rule.Element)
			}
		}
	}
}
