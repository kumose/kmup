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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

// CamoEncode encodes a lnk to fit with the go-camo and camo proxy links. The purposes of camo-proxy are:
// 1. Allow accessing "http://" images on a HTTPS site by using the "https://" URLs provided by camo-proxy.
// 2. Hide the visitor's real IP (protect privacy) when accessing external images.
func CamoEncode(link string) string {
	if strings.HasPrefix(link, setting.Camo.ServerURL) {
		return link
	}

	mac := hmac.New(sha1.New, []byte(setting.Camo.HMACKey))
	_, _ = mac.Write([]byte(link)) // hmac does not return errors
	macSum := b64encode(mac.Sum(nil))
	encodedURL := b64encode([]byte(link))

	return util.URLJoin(setting.Camo.ServerURL, macSum, encodedURL)
}

func b64encode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func camoHandleLink(link string) string {
	if setting.Camo.Enabled {
		lnkURL, err := url.Parse(link)
		if err == nil && lnkURL.IsAbs() && !strings.HasPrefix(link, setting.AppURL) &&
			(setting.Camo.Always || lnkURL.Scheme != "https") {
			return CamoEncode(link)
		}
	}
	return link
}
