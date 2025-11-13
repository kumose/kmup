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

	"github.com/microcosm-cc/bluemonday"
)

// createRepoDescriptionPolicy returns a minimal more strict policy that is used for
// repository descriptions.
func (st *Sanitizer) createRepoDescriptionPolicy() *bluemonday.Policy {
	policy := bluemonday.NewPolicy()
	policy.AllowStandardURLs()

	// Allow italics and bold.
	policy.AllowElements("i", "b", "em", "strong")

	// Allow code.
	policy.AllowElements("code")

	// Allow links
	policy.AllowAttrs("href", "target", "rel").OnElements("a")

	// Allow classes for emojis
	policy.AllowAttrs("class").Matching(regexp.MustCompile(`^emoji$`)).OnElements("img", "span")
	policy.AllowAttrs("aria-label").OnElements("span")

	return policy
}

// SanitizeDescription sanitizes the HTML generated for a repository description.
func SanitizeDescription(s string) string {
	return GetDefaultSanitizer().descriptionPolicy.Sanitize(s)
}
