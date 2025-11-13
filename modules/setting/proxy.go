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
	"net/url"

	"github.com/kumose/kmup/modules/log"
)

// Proxy settings
var Proxy = struct {
	Enabled       bool
	ProxyURL      string
	ProxyURLFixed *url.URL
	ProxyHosts    []string
}{
	Enabled:    false,
	ProxyURL:   "",
	ProxyHosts: []string{},
}

func loadProxyFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("proxy")
	Proxy.Enabled = sec.Key("PROXY_ENABLED").MustBool(false)
	Proxy.ProxyURL = sec.Key("PROXY_URL").MustString("")
	if Proxy.ProxyURL != "" {
		var err error
		Proxy.ProxyURLFixed, err = url.Parse(Proxy.ProxyURL)
		if err != nil {
			log.Error("Global PROXY_URL is not valid")
			Proxy.ProxyURL = ""
		}
	}
	Proxy.ProxyHosts = sec.Key("PROXY_HOSTS").Strings(",")
}
