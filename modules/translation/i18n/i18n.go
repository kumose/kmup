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

package i18n

import (
	"html/template"
	"io"
)

var DefaultLocales = NewLocaleStore()

type Locale interface {
	// TrString translates a given key and arguments for a language
	TrString(trKey string, trArgs ...any) string
	// TrHTML translates a given key and arguments for a language, string arguments are escaped to HTML
	TrHTML(trKey string, trArgs ...any) template.HTML
	// HasKey reports if a locale has a translation for a given key
	HasKey(trKey string) bool
}

// LocaleStore provides the functions common to all locale stores
type LocaleStore interface {
	io.Closer

	// SetDefaultLang sets the default language to fall back to
	SetDefaultLang(lang string)
	// ListLangNameDesc provides paired slices of language names to descriptors
	ListLangNameDesc() (names, desc []string)
	// Locale return the locale for the provided language or the default language if not found
	Locale(langName string) (Locale, bool)
	// HasLang returns whether a given language is present in the store
	HasLang(langName string) bool
	// AddLocaleByIni adds a new language to the store
	AddLocaleByIni(langName, langDesc string, source, moreSource []byte) error
}

// ResetDefaultLocales resets the current default locales
// NOTE: this is not synchronized
func ResetDefaultLocales() {
	if DefaultLocales != nil {
		_ = DefaultLocales.Close()
	}
	DefaultLocales = NewLocaleStore()
}

// GetLocale returns the locale from the default locales
func GetLocale(lang string) (Locale, bool) {
	return DefaultLocales.Locale(lang)
}
