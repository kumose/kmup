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

// DEPRECATED should not be removed because users maybe upgrade from lower version to the latest version
// if these are removed, the warning will not be shown
// - will need to set default for [queue.task] LENGTH to 1000 though
func loadTaskFrom(rootCfg ConfigProvider) {
	taskSec := rootCfg.Section("task")
	queueTaskSec := rootCfg.Section("queue.task")

	deprecatedSetting(rootCfg, "task", "QUEUE_TYPE", "queue.task", "TYPE", "v1.19.0")
	deprecatedSetting(rootCfg, "task", "QUEUE_CONN_STR", "queue.task", "CONN_STR", "v1.19.0")
	deprecatedSetting(rootCfg, "task", "QUEUE_LENGTH", "queue.task", "LENGTH", "v1.19.0")

	switch taskSec.Key("QUEUE_TYPE").MustString("channel") {
	case "channel":
		queueTaskSec.Key("TYPE").MustString("persistable-channel")
		queueTaskSec.Key("CONN_STR").MustString(taskSec.Key("QUEUE_CONN_STR").MustString(""))
	case "redis":
		queueTaskSec.Key("TYPE").MustString("redis")
		queueTaskSec.Key("CONN_STR").MustString(taskSec.Key("QUEUE_CONN_STR").MustString("addrs=127.0.0.1:6379 db=0"))
	}
	queueTaskSec.Key("LENGTH").MustInt(taskSec.Key("QUEUE_LENGTH").MustInt(1000))
}
