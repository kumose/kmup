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

package auth

import (
	"net/http"
	"strings"

	user_model "github.com/kumose/kmup/models/user"
)

// Ensure the struct implements the interface.
var (
	_ Method = &Group{}
)

// Group implements the Auth interface with serval Auth.
type Group struct {
	methods []Method
}

// NewGroup creates a new auth group
func NewGroup(methods ...Method) *Group {
	return &Group{
		methods: methods,
	}
}

// Add adds a new method to group
func (b *Group) Add(method Method) {
	b.methods = append(b.methods, method)
}

// Name returns group's methods name
func (b *Group) Name() string {
	names := make([]string, 0, len(b.methods))
	for _, m := range b.methods {
		names = append(names, m.Name())
	}
	return strings.Join(names, ",")
}

func (b *Group) Verify(req *http.Request, w http.ResponseWriter, store DataStore, sess SessionStore) (*user_model.User, error) {
	// Try to sign in with each of the enabled plugins
	var retErr error
	for _, m := range b.methods {
		user, err := m.Verify(req, w, store, sess)
		if err != nil {
			if retErr == nil {
				retErr = err
			}
			// Try other methods if this one failed.
			// Some methods may share the same protocol to detect if they are matched.
			// For example, OAuth2 and conan.Auth both read token from "Authorization: Bearer <token>" header,
			// If OAuth2 returns error, we should give conan.Auth a chance to try.
			continue
		}

		// If any method returns a user, we can stop trying.
		// Return the user and ignore any error returned by previous methods.
		if user != nil {
			if store.GetData()["AuthedMethod"] == nil {
				store.GetData()["AuthedMethod"] = m.Name()
			}
			return user, nil
		}
	}

	// If no method returns a user, return the error returned by the first method.
	return nil, retErr
}
