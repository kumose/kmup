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

package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kumose/kmup/modules/auth"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/glob"
	"github.com/kumose/kmup/modules/util"

	"github.com/kumose-go/chi/binding"
)

const (
	// ErrGitRefName is git reference name error
	ErrGitRefName = "GitRefNameError"
	// ErrGlobPattern is returned when glob pattern is invalid
	ErrGlobPattern = "GlobPattern"
	// ErrRegexPattern is returned when a regex pattern is invalid
	ErrRegexPattern = "RegexPattern"
	// ErrUsername is username error
	ErrUsername = "UsernameError"
	// ErrInvalidGroupTeamMap is returned when a group team mapping is invalid
	ErrInvalidGroupTeamMap = "InvalidGroupTeamMap"
)

// AddBindingRules adds additional binding rules
func AddBindingRules() {
	addGitRefNameBindingRule()
	addValidURLListBindingRule()
	addValidURLBindingRule()
	addValidSiteURLBindingRule()
	addGlobPatternRule()
	addRegexPatternRule()
	addGlobOrRegexPatternRule()
	addUsernamePatternRule()
	addValidGroupTeamMapRule()
}

func addGitRefNameBindingRule() {
	// Git refname validation rule
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "GitRefName"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := fmt.Sprintf("%v", val)

			if !git.IsValidRefPattern(str) {
				errs.Add([]string{name}, ErrGitRefName, "GitRefName")
				return false, errs
			}
			return true, errs
		},
	})
}

func addValidURLListBindingRule() {
	// URL validation rule
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "ValidUrlList"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := fmt.Sprintf("%v", val)
			if len(str) == 0 {
				errs.Add([]string{name}, binding.ERR_URL, "Url")
				return false, errs
			}

			ok := true
			urls := util.SplitTrimSpace(str, "\n")
			for _, u := range urls {
				if !IsValidURL(u) {
					ok = false
					errs.Add([]string{name}, binding.ERR_URL, u)
				}
			}

			return ok, errs
		},
	})
}

func addValidURLBindingRule() {
	// URL validation rule
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "ValidUrl"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := fmt.Sprintf("%v", val)
			if len(str) != 0 && !IsValidURL(str) {
				errs.Add([]string{name}, binding.ERR_URL, "Url")
				return false, errs
			}

			return true, errs
		},
	})
}

func addValidSiteURLBindingRule() {
	// URL validation rule
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "ValidSiteUrl"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := fmt.Sprintf("%v", val)
			if len(str) != 0 && !IsValidSiteURL(str) {
				errs.Add([]string{name}, binding.ERR_URL, "Url")
				return false, errs
			}

			return true, errs
		},
	})
}

func addGlobPatternRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "GlobPattern"
		},
		IsValid: globPatternValidator,
	})
}

func globPatternValidator(errs binding.Errors, name string, val any) (bool, binding.Errors) {
	str := fmt.Sprintf("%v", val)

	if len(str) != 0 {
		if _, err := glob.Compile(str); err != nil {
			errs.Add([]string{name}, ErrGlobPattern, err.Error())
			return false, errs
		}
	}

	return true, errs
}

func addRegexPatternRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "RegexPattern"
		},
		IsValid: regexPatternValidator,
	})
}

func regexPatternValidator(errs binding.Errors, name string, val any) (bool, binding.Errors) {
	str := fmt.Sprintf("%v", val)

	if _, err := regexp.Compile(str); err != nil {
		errs.Add([]string{name}, ErrRegexPattern, err.Error())
		return false, errs
	}

	return true, errs
}

func addGlobOrRegexPatternRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "GlobOrRegexPattern"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := strings.TrimSpace(fmt.Sprintf("%v", val))

			if len(str) >= 2 && strings.HasPrefix(str, "/") && strings.HasSuffix(str, "/") {
				return regexPatternValidator(errs, name, str[1:len(str)-1])
			}
			return globPatternValidator(errs, name, val)
		},
	})
}

func addUsernamePatternRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "Username"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := fmt.Sprintf("%v", val)
			if !IsValidUsername(str) {
				errs.Add([]string{name}, ErrUsername, "invalid username")
				return false, errs
			}
			return true, errs
		},
	})
}

func addValidGroupTeamMapRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "ValidGroupTeamMap"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			_, err := auth.UnmarshalGroupTeamMapping(fmt.Sprintf("%v", val))
			if err != nil {
				errs.Add([]string{name}, ErrInvalidGroupTeamMap, err.Error())
				return false, errs
			}

			return true, errs
		},
	})
}

func portOnly(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return ""
	}
	if i := strings.Index(hostport, "]:"); i != -1 {
		return hostport[i+len("]:"):]
	}
	if strings.Contains(hostport, "]") {
		return ""
	}
	return hostport[colon+len(":"):]
}

func validPort(p string) bool {
	for _, r := range []byte(p) {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
