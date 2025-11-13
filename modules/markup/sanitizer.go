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
	"sync"

	"github.com/microcosm-cc/bluemonday"
)

// Sanitizer is a protection wrapper of *bluemonday.Policy which does not allow
// any modification to the underlying policies once it's been created.
type Sanitizer struct {
	defaultPolicy     *bluemonday.Policy
	descriptionPolicy *bluemonday.Policy
	rendererPolicies  map[string]*bluemonday.Policy
	allowAllRegex     *regexp.Regexp
}

var (
	defaultSanitizer     *Sanitizer
	defaultSanitizerOnce sync.Once
)

func GetDefaultSanitizer() *Sanitizer {
	defaultSanitizerOnce.Do(func() {
		defaultSanitizer = &Sanitizer{
			rendererPolicies: map[string]*bluemonday.Policy{},
			allowAllRegex:    regexp.MustCompile(".+"),
		}
		for name, renderer := range renderers {
			sanitizerRules := renderer.SanitizerRules()
			if len(sanitizerRules) > 0 {
				policy := defaultSanitizer.createDefaultPolicy()
				defaultSanitizer.addSanitizerRules(policy, sanitizerRules)
				defaultSanitizer.rendererPolicies[name] = policy
			}
		}
		defaultSanitizer.defaultPolicy = defaultSanitizer.createDefaultPolicy()
		defaultSanitizer.descriptionPolicy = defaultSanitizer.createRepoDescriptionPolicy()
	})
	return defaultSanitizer
}

func ResetDefaultSanitizerForTesting() {
	defaultSanitizer = nil
	defaultSanitizerOnce = sync.Once{}
}
