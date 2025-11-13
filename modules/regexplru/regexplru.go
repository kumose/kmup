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

package regexplru

import (
	"regexp"

	"github.com/kumose/kmup/modules/log"

	lru "github.com/hashicorp/golang-lru/v2"
)

var lruCache *lru.Cache[string, any]

func init() {
	var err error
	lruCache, err = lru.New[string, any](1000)
	if err != nil {
		log.Fatal("failed to new LRU cache, err: %v", err)
	}
}

// GetCompiled works like regexp.Compile, the compiled expr or error is stored in LRU cache
func GetCompiled(expr string) (r *regexp.Regexp, err error) {
	v, ok := lruCache.Get(expr)
	if !ok {
		r, err = regexp.Compile(expr)
		if err != nil {
			lruCache.Add(expr, err)
			return nil, err
		}
		lruCache.Add(expr, r)
	} else {
		r, ok = v.(*regexp.Regexp)
		if !ok {
			if err, ok = v.(error); ok {
				return nil, err
			}
			panic("impossible")
		}
	}
	return r, nil
}
