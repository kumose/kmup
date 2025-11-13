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

package oauth2

import (
	"context"
	"encoding/gob"
	"net/http"
	"sync"

	"github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/optional"
	"github.com/kumose/kmup/modules/setting"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

var gothRWMutex = sync.RWMutex{}

// UsersStoreKey is the key for the store
const UsersStoreKey = "kmup-oauth2-sessions"

// ProviderHeaderKey is the HTTP header key
const ProviderHeaderKey = "kmup-oauth2-provider"

// Init initializes the oauth source
func Init(ctx context.Context) error {
	// Lock our mutex
	gothRWMutex.Lock()

	gob.Register(&sessions.Session{})

	gothic.Store = &SessionsStore{
		maxLength: int64(setting.OAuth2.MaxTokenLength),
	}

	gothic.SetState = func(req *http.Request) string {
		return uuid.New().String()
	}

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return req.Header.Get(ProviderHeaderKey), nil
	}

	// Unlock our mutex
	gothRWMutex.Unlock()

	return initOAuth2Sources(ctx)
}

// ResetOAuth2 clears existing OAuth2 providers and loads them from DB
func ResetOAuth2(ctx context.Context) error {
	ClearProviders()
	return initOAuth2Sources(ctx)
}

// initOAuth2Sources is used to load and register all active OAuth2 providers
func initOAuth2Sources(ctx context.Context) error {
	authSources, err := db.Find[auth.Source](ctx, auth.FindSourcesOptions{
		IsActive:  optional.Some(true),
		LoginType: auth.OAuth2,
	})
	if err != nil {
		return err
	}
	for _, source := range authSources {
		oauth2Source, ok := source.Cfg.(*Source)
		if !ok {
			continue
		}
		err := oauth2Source.RegisterSource()
		if err != nil {
			log.Critical("Unable to register source: %s due to Error: %v.", source.Name, err)
		}
	}
	return nil
}
