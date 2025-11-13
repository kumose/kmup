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

// Webhook settings
var Webhook = struct {
	QueueLength     int
	DeliverTimeout  int
	SkipTLSVerify   bool
	AllowedHostList string
	Types           []string
	PagingNum       int
	ProxyURL        string
	ProxyURLFixed   *url.URL
	ProxyHosts      []string
}{
	QueueLength:    1000,
	DeliverTimeout: 5,
	SkipTLSVerify:  false,
	PagingNum:      10,
	ProxyURL:       "",
	ProxyHosts:     []string{},
}

func loadWebhookFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("webhook")
	Webhook.QueueLength = sec.Key("QUEUE_LENGTH").MustInt(1000)
	Webhook.DeliverTimeout = sec.Key("DELIVER_TIMEOUT").MustInt(5)
	Webhook.SkipTLSVerify = sec.Key("SKIP_TLS_VERIFY").MustBool()
	Webhook.AllowedHostList = sec.Key("ALLOWED_HOST_LIST").MustString("")
	Webhook.Types = []string{"kmup", "gogs", "slack", "discord", "dingtalk", "telegram", "msteams", "feishu", "matrix", "wechatwork", "packagist"}
	Webhook.PagingNum = sec.Key("PAGING_NUM").MustInt(10)
	Webhook.ProxyURL = sec.Key("PROXY_URL").MustString("")
	if Webhook.ProxyURL != "" {
		var err error
		Webhook.ProxyURLFixed, err = url.Parse(Webhook.ProxyURL)
		if err != nil {
			log.Error("Webhook PROXY_URL is not valid")
			Webhook.ProxyURL = ""
		}
	}
	Webhook.ProxyHosts = sec.Key("PROXY_HOSTS").Strings(",")
}
