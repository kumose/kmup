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

package setting

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/util"
)

// SessionConfig defines Session settings
var SessionConfig = struct {
	OriginalProvider string
	Provider         string
	// Provider configuration, it's corresponding to provider.
	ProviderConfig string
	// Cookie name to save session ID. Default is "MacaronSession".
	CookieName string
	// Cookie path to store. Default is "/".
	CookiePath string
	// GC interval time in seconds. Default is 3600.
	Gclifetime int64
	// Max life time in seconds. Default is whatever GC interval time is.
	Maxlifetime int64
	// Use HTTPS only. Default is false.
	Secure bool
	// Cookie domain name. Default is empty.
	Domain string
	// SameSite declares if your cookie should be restricted to a first-party or same-site context. Valid strings are "none", "lax", "strict". Default is "lax"
	SameSite http.SameSite
}{
	CookieName:  "i_like_kmup",
	Gclifetime:  86400,
	Maxlifetime: 86400,
	SameSite:    http.SameSiteLaxMode,
}

func loadSessionFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("session")
	SessionConfig.Provider = sec.Key("PROVIDER").In("memory",
		[]string{"memory", "file", "redis", "mysql", "postgres", "couchbase", "memcache", "db"})
	SessionConfig.ProviderConfig = strings.Trim(sec.Key("PROVIDER_CONFIG").MustString(filepath.Join(AppDataPath, "sessions")), "\" ")
	if SessionConfig.Provider == "file" && !filepath.IsAbs(SessionConfig.ProviderConfig) {
		SessionConfig.ProviderConfig = filepath.Join(AppWorkPath, SessionConfig.ProviderConfig)
		checkOverlappedPath("[session].PROVIDER_CONFIG", SessionConfig.ProviderConfig)
	}
	SessionConfig.CookieName = sec.Key("COOKIE_NAME").MustString("i_like_kmup")
	// HINT: INSTALL-PAGE-COOKIE-INIT: the cookie system is not properly initialized on the Install page, so there is no CookiePath
	SessionConfig.CookiePath = util.IfZero(AppSubURL, "/")
	SessionConfig.Secure = sec.Key("COOKIE_SECURE").MustBool(strings.HasPrefix(strings.ToLower(AppURL), "https://"))
	SessionConfig.Gclifetime = sec.Key("GC_INTERVAL_TIME").MustInt64(86400)
	SessionConfig.Maxlifetime = sec.Key("SESSION_LIFE_TIME").MustInt64(86400)
	SessionConfig.Domain = sec.Key("DOMAIN").String()
	samesiteString := sec.Key("SAME_SITE").In("lax", []string{"none", "lax", "strict"})
	switch strings.ToLower(samesiteString) {
	case "none":
		SessionConfig.SameSite = http.SameSiteNoneMode
	case "strict":
		SessionConfig.SameSite = http.SameSiteStrictMode
	default:
		SessionConfig.SameSite = http.SameSiteLaxMode
	}
	shadowConfig, err := json.Marshal(SessionConfig)
	if err != nil {
		log.Fatal("Can't shadow session config: %v", err)
	}
	SessionConfig.ProviderConfig = string(shadowConfig)
	SessionConfig.OriginalProvider = SessionConfig.Provider
	SessionConfig.Provider = "VirtualSession"
}
