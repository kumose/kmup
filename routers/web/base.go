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

package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/kumose/kmup/modules/httpcache"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/storage"
	"github.com/kumose/kmup/modules/util"
	"github.com/kumose/kmup/modules/web/routing"
)

func avatarStorageHandler(storageSetting *setting.Storage, prefix string, objStore storage.ObjectStorage) http.HandlerFunc {
	prefix = strings.Trim(prefix, "/")
	funcInfo := routing.GetFuncInfo(avatarStorageHandler, prefix)

	if storageSetting.ServeDirect() {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.Method != http.MethodGet && req.Method != http.MethodHead {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}

			if !strings.HasPrefix(req.URL.Path, "/"+prefix+"/") {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			defer routing.RecordFuncInfo(req.Context(), funcInfo)()

			rPath := strings.TrimPrefix(req.URL.Path, "/"+prefix+"/")
			rPath = util.PathJoinRelX(rPath)

			u, err := objStore.URL(rPath, path.Base(rPath), req.Method, nil)
			if err != nil {
				if os.IsNotExist(err) || errors.Is(err, os.ErrNotExist) {
					log.Warn("Unable to find %s %s", prefix, rPath)
					http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
				log.Error("Error whilst getting URL for %s %s. Error: %v", prefix, rPath, err)
				http.Error(w, fmt.Sprintf("Error whilst getting URL for %s %s", prefix, rPath), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, req, u.String(), http.StatusTemporaryRedirect)
		}
	}

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if !strings.HasPrefix(req.URL.Path, "/"+prefix+"/") {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer routing.RecordFuncInfo(req.Context(), funcInfo)()

		rPath := strings.TrimPrefix(req.URL.Path, "/"+prefix+"/")
		rPath = util.PathJoinRelX(rPath)
		if rPath == "" || rPath == "." {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		fi, err := objStore.Stat(rPath)
		if err != nil {
			if os.IsNotExist(err) || errors.Is(err, os.ErrNotExist) {
				log.Warn("Unable to find %s %s", prefix, rPath)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Error("Error whilst opening %s %s. Error: %v", prefix, rPath, err)
			http.Error(w, fmt.Sprintf("Error whilst opening %s %s", prefix, rPath), http.StatusInternalServerError)
			return
		}

		fr, err := objStore.Open(rPath)
		if err != nil {
			log.Error("Error whilst opening %s %s. Error: %v", prefix, rPath, err)
			http.Error(w, fmt.Sprintf("Error whilst opening %s %s", prefix, rPath), http.StatusInternalServerError)
			return
		}
		defer fr.Close()

		httpcache.SetCacheControlInHeader(w.Header(), httpcache.CacheControlForPublicStatic())
		http.ServeContent(w, req, path.Base(rPath), fi.ModTime(), fr)
	}
}
