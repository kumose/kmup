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
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/nosql"
)

// GlobalLock represents configuration of global lock
var GlobalLock = struct {
	ServiceType    string
	ServiceConnStr string
}{
	ServiceType: "memory",
}

func loadGlobalLockFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("global_lock")
	GlobalLock.ServiceType = sec.Key("SERVICE_TYPE").MustString("memory")
	switch GlobalLock.ServiceType {
	case "memory":
	case "redis":
		connStr := sec.Key("SERVICE_CONN_STR").String()
		if connStr == "" {
			log.Fatal("SERVICE_CONN_STR is empty for redis")
		}
		u := nosql.ToRedisURI(connStr)
		if u == nil {
			log.Fatal("SERVICE_CONN_STR %s is not a valid redis connection string", connStr)
		}
		GlobalLock.ServiceConnStr = connStr
	default:
		log.Fatal("Unknown sync lock service type: %s", GlobalLock.ServiceType)
	}
}
