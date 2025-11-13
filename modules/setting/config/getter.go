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

package config

import (
	"context"
	"sync"
)

var getterMu sync.RWMutex

type CfgSecKeyGetter interface {
	GetValue(sec, key string) (v string, has bool)
}

var cfgSecKeyGetterInternal CfgSecKeyGetter

func SetCfgSecKeyGetter(p CfgSecKeyGetter) {
	getterMu.Lock()
	cfgSecKeyGetterInternal = p
	getterMu.Unlock()
}

func GetCfgSecKeyGetter() CfgSecKeyGetter {
	getterMu.RLock()
	defer getterMu.RUnlock()
	return cfgSecKeyGetterInternal
}

type DynKeyGetter interface {
	GetValue(ctx context.Context, key string) (v string, has bool)
	GetRevision(ctx context.Context) int
	InvalidateCache()
}

var dynKeyGetterInternal DynKeyGetter

func SetDynGetter(p DynKeyGetter) {
	getterMu.Lock()
	dynKeyGetterInternal = p
	getterMu.Unlock()
}

func GetDynGetter() DynKeyGetter {
	getterMu.RLock()
	defer getterMu.RUnlock()
	return dynKeyGetterInternal
}
