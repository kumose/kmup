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

package attribute

import (
	"strings"

	"github.com/kumose/kmup/modules/optional"
)

type Attribute string

const (
	LinguistVendored      = "linguist-vendored"
	LinguistGenerated     = "linguist-generated"
	LinguistDocumentation = "linguist-documentation"
	LinguistDetectable    = "linguist-detectable"
	LinguistLanguage      = "linguist-language"
	GitlabLanguage        = "gitlab-language"
	Lockable              = "lockable"
	Filter                = "filter"
	Diff                  = "diff"
)

var LinguistAttributes = []string{
	LinguistVendored,
	LinguistGenerated,
	LinguistDocumentation,
	LinguistDetectable,
	LinguistLanguage,
	GitlabLanguage,
}

func (a Attribute) IsUnspecified() bool {
	return a == "" || a == "unspecified"
}

func (a Attribute) ToString() optional.Option[string] {
	if !a.IsUnspecified() {
		return optional.Some(string(a))
	}
	return optional.None[string]()
}

// ToBool converts the attribute value to optional boolean: true if "set"/"true", false if "unset"/"false", none otherwise
func (a Attribute) ToBool() optional.Option[bool] {
	switch a {
	case "set", "true":
		return optional.Some(true)
	case "unset", "false":
		return optional.Some(false)
	}
	return optional.None[bool]()
}

type Attributes struct {
	m map[string]Attribute
}

func NewAttributes() *Attributes {
	return &Attributes{m: make(map[string]Attribute)}
}

func (attrs *Attributes) Get(name string) Attribute {
	if value, has := attrs.m[name]; has {
		return value
	}
	return ""
}

func (attrs *Attributes) GetVendored() optional.Option[bool] {
	return attrs.Get(LinguistVendored).ToBool()
}

func (attrs *Attributes) GetGenerated() optional.Option[bool] {
	return attrs.Get(LinguistGenerated).ToBool()
}

func (attrs *Attributes) GetDocumentation() optional.Option[bool] {
	return attrs.Get(LinguistDocumentation).ToBool()
}

func (attrs *Attributes) GetDetectable() optional.Option[bool] {
	return attrs.Get(LinguistDetectable).ToBool()
}

func (attrs *Attributes) GetLinguistLanguage() optional.Option[string] {
	return attrs.Get(LinguistLanguage).ToString()
}

func (attrs *Attributes) GetGitlabLanguage() optional.Option[string] {
	attrStr := attrs.Get(GitlabLanguage).ToString()
	if attrStr.Has() {
		raw := attrStr.Value()
		// gitlab-language may have additional parameters after the language
		// ignore them and just use the master language
		// https://docs.gitlab.com/ee/user/project/highlighting.html#override-syntax-highlighting-for-a-file-type
		if idx := strings.IndexByte(raw, '?'); idx >= 0 {
			return optional.Some(raw[:idx])
		}
	}
	return attrStr
}

func (attrs *Attributes) GetLanguage() optional.Option[string] {
	// prefer linguist-language over gitlab-language
	// if linguist-language is not set, use gitlab-language
	// if both are not set, return none
	language := attrs.GetLinguistLanguage()
	if language.Value() == "" {
		language = attrs.GetGitlabLanguage()
	}
	return language
}
