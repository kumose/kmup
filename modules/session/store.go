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

package session

import (
	"net/http"

	"github.com/kumose/kmup/modules/setting"

	"github.com/kumose-go/chi/session"
)

type RawStore = session.RawStore

type Store interface {
	RawStore
	Destroy(http.ResponseWriter, *http.Request) error
}

type mockStoreContextKeyStruct struct{}

var MockStoreContextKey = mockStoreContextKeyStruct{}

// RegenerateSession regenerates the underlying session and returns the new store
func RegenerateSession(resp http.ResponseWriter, req *http.Request) (Store, error) {
	for _, f := range BeforeRegenerateSession {
		f(resp, req)
	}
	if setting.IsInTesting {
		if store := req.Context().Value(MockStoreContextKey); store != nil {
			return store.(Store), nil
		}
	}
	return session.RegenerateSession(resp, req)
}

func GetContextSession(req *http.Request) Store {
	if setting.IsInTesting {
		if store := req.Context().Value(MockStoreContextKey); store != nil {
			return store.(Store)
		}
	}
	return session.GetSession(req)
}

// BeforeRegenerateSession is a list of functions that are called before a session is regenerated.
var BeforeRegenerateSession []func(http.ResponseWriter, *http.Request)
