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

package public

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/kumose/kmup/modules/assetfs"
	"github.com/kumose/kmup/modules/container"
	"github.com/kumose/kmup/modules/httpcache"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

func CustomAssets() *assetfs.Layer {
	return assetfs.Local("custom", setting.CustomPath, "public")
}

func AssetFS() *assetfs.LayeredFS {
	return assetfs.Layered(CustomAssets(), BuiltinAssets())
}

// FileHandlerFunc implements the static handler for serving files in "public" assets
func FileHandlerFunc() http.HandlerFunc {
	assetFS := AssetFS()
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" && req.Method != "HEAD" {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handleRequest(resp, req, assetFS, req.URL.Path)
	}
}

// parseAcceptEncoding parse Accept-Encoding: deflate, gzip;q=1.0, *;q=0.5 as compress methods
func parseAcceptEncoding(val string) container.Set[string] {
	parts := strings.Split(val, ";")
	types := make(container.Set[string])
	for v := range strings.SplitSeq(parts[0], ",") {
		types.Add(strings.TrimSpace(v))
	}
	return types
}

// setWellKnownContentType will set the Content-Type if the file is a well-known type.
// See the comments of detectWellKnownMimeType
func setWellKnownContentType(w http.ResponseWriter, file string) {
	mimeType := detectWellKnownMimeType(path.Ext(file))
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}
}

func handleRequest(w http.ResponseWriter, req *http.Request, fs http.FileSystem, file string) {
	// actually, fs (http.FileSystem) is designed to be a safe interface, relative paths won't bypass its parent directory, it's also fine to do a clean here
	f, err := fs.Open(util.PathJoinRelX(file))
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("[Static] Open %q failed: %v", file, err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("[Static] %q exists, but fails to open: %v", file, err)
		return
	}

	// need to serve index file? (no at the moment)
	if fi.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	servePublicAsset(w, req, fi, fi.ModTime(), f)
}

// servePublicAsset serve http content
func servePublicAsset(w http.ResponseWriter, req *http.Request, fi os.FileInfo, modtime time.Time, content io.ReadSeeker) {
	setWellKnownContentType(w, fi.Name())
	httpcache.SetCacheControlInHeader(w.Header(), httpcache.CacheControlForPublicStatic())
	encodings := parseAcceptEncoding(req.Header.Get("Accept-Encoding"))
	fiEmbedded, _ := fi.(assetfs.EmbeddedFileInfo)
	if encodings.Contains("gzip") && fiEmbedded != nil {
		// try to provide gzip content directly from bindata
		if gzipBytes, ok := fiEmbedded.GetGzipContent(); ok {
			rdGzip := bytes.NewReader(gzipBytes)
			// all gzipped static files (from bindata) are managed by Kmup, so we can make sure every file has the correct ext name
			// then we can get the correct Content-Type, we do not need to do http.DetectContentType on the decompressed data
			if w.Header().Get("Content-Type") == "" {
				w.Header().Set("Content-Type", "application/octet-stream")
			}
			w.Header().Set("Content-Encoding", "gzip")
			http.ServeContent(w, req, fi.Name(), modtime, rdGzip)
			return
		}
	}
	http.ServeContent(w, req, fi.Name(), modtime, content)
}
