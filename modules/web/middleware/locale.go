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

package middleware

import (
	"net/http"

	"github.com/kumose/kmup/modules/translation"
	"github.com/kumose/kmup/modules/translation/i18n"

	"golang.org/x/text/language"
)

// Locale handle locale
func Locale(resp http.ResponseWriter, req *http.Request) translation.Locale {
	// 1. Check URL arguments.
	lang := req.URL.Query().Get("lang")
	changeLang := lang != ""

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		ck, _ := req.Cookie("lang")
		if ck != nil {
			lang = ck.Value
		}
	}

	// Check again in case someone changes the supported language list.
	if lang != "" && !i18n.DefaultLocales.HasLang(lang) {
		lang = ""
		changeLang = false
	}

	// 3. Get language information from 'Accept-Language'.
	// The first element in the list is chosen to be the default language automatically.
	if len(lang) == 0 {
		tags, _, _ := language.ParseAcceptLanguage(req.Header.Get("Accept-Language"))
		tag := translation.Match(tags...)
		lang = tag.String()
	}

	if changeLang {
		SetLocaleCookie(resp, lang, 1<<31-1)
	}

	return translation.NewLocale(lang)
}

// SetLocaleCookie convenience function to set the locale cookie consistently
func SetLocaleCookie(resp http.ResponseWriter, lang string, maxAge int) {
	SetSiteCookie(resp, "lang", lang, maxAge)
}

// DeleteLocaleCookie convenience function to delete the locale cookie consistently
// Setting the lang cookie will trigger the middleware to reset the language to previous state.
func DeleteLocaleCookie(resp http.ResponseWriter) {
	SetSiteCookie(resp, "lang", "", -1)
}
