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

package queue

import (
	"github.com/kumose/kmup/modules/setting"
)

type BaseConfig struct {
	ManagedName string
	DataFullDir string // the caller must prepare an absolute path

	ConnStr string
	Length  int

	QueueFullName, SetFullName string
}

func toBaseConfig(managedName string, queueSetting setting.QueueSettings) *BaseConfig {
	baseConfig := &BaseConfig{
		ManagedName: managedName,
		DataFullDir: queueSetting.Datadir,

		ConnStr: queueSetting.ConnStr,
		Length:  queueSetting.Length,
	}

	// queue name and set name
	baseConfig.QueueFullName = managedName + queueSetting.QueueName
	baseConfig.SetFullName = baseConfig.QueueFullName + queueSetting.SetName
	if baseConfig.SetFullName == baseConfig.QueueFullName {
		baseConfig.SetFullName += "_unique"
	}
	return baseConfig
}
