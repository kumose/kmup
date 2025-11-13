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

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/util"
)

type CfgSecKey struct {
	Sec, Key string
}

type Value[T any] struct {
	mu sync.RWMutex

	cfgSecKey CfgSecKey
	dynKey    string

	def, value T
	revision   int
}

func (value *Value[T]) parse(key, valStr string) (v T) {
	v = value.def
	if valStr != "" {
		if err := json.Unmarshal(util.UnsafeStringToBytes(valStr), &v); err != nil {
			log.Error("Unable to unmarshal json config for key %q, err: %v", key, err)
		}
	}
	return v
}

func (value *Value[T]) Value(ctx context.Context) (v T) {
	dg := GetDynGetter()
	if dg == nil {
		// this is an edge case: the database is not initialized but the system setting is going to be used
		// it should panic to avoid inconsistent config values (from config / system setting) and fix the code
		panic("no config dyn value getter")
	}

	rev := dg.GetRevision(ctx)

	// if the revision in the database doesn't change, use the last value
	value.mu.RLock()
	if rev == value.revision {
		v = value.value
		value.mu.RUnlock()
		return v
	}
	value.mu.RUnlock()

	// try to parse the config and cache it
	var valStr *string
	if dynVal, has := dg.GetValue(ctx, value.dynKey); has {
		valStr = &dynVal
	} else if cfgVal, has := GetCfgSecKeyGetter().GetValue(value.cfgSecKey.Sec, value.cfgSecKey.Key); has {
		valStr = &cfgVal
	}
	if valStr == nil {
		v = value.def
	} else {
		v = value.parse(value.dynKey, *valStr)
	}

	value.mu.Lock()
	value.value = v
	value.revision = rev
	value.mu.Unlock()
	return v
}

func (value *Value[T]) DynKey() string {
	return value.dynKey
}

func (value *Value[T]) WithDefault(def T) *Value[T] {
	value.def = def
	return value
}

func (value *Value[T]) DefaultValue() T {
	return value.def
}

func (value *Value[T]) WithFileConfig(cfgSecKey CfgSecKey) *Value[T] {
	value.cfgSecKey = cfgSecKey
	return value
}

func ValueJSON[T any](dynKey string) *Value[T] {
	return &Value[T]{dynKey: dynKey}
}
