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

package httpcache

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

type CacheControlOptions struct {
	IsPublic    bool
	MaxAge      time.Duration
	NoTransform bool
}

// SetCacheControlInHeader sets suitable cache-control headers in the response
func SetCacheControlInHeader(h http.Header, opts *CacheControlOptions) {
	directives := make([]string, 0, 4)

	// "max-age=0 + must-revalidate" (aka "no-cache") is preferred instead of "no-store"
	// because browsers may restore some input fields after navigate-back / reload a page.
	publicPrivate := util.Iif(opts.IsPublic, "public", "private")
	if setting.IsProd {
		if opts.MaxAge == 0 {
			directives = append(directives, "max-age=0", "private", "must-revalidate")
		} else {
			directives = append(directives, publicPrivate, "max-age="+strconv.Itoa(int(opts.MaxAge.Seconds())))
		}
	} else {
		// use dev-related controls, and remind users they are using non-prod setting.
		directives = append(directives, "max-age=0", publicPrivate, "must-revalidate")
		h.Set("X-Kmup-Debug", fmt.Sprintf("RUN_MODE=%v, MaxAge=%s", setting.RunMode, opts.MaxAge))
	}

	if opts.NoTransform {
		directives = append(directives, "no-transform")
	}
	h.Set("Cache-Control", strings.Join(directives, ", "))
}

func CacheControlForPublicStatic() *CacheControlOptions {
	return &CacheControlOptions{
		IsPublic:    true,
		MaxAge:      setting.StaticCacheTime,
		NoTransform: true,
	}
}

func CacheControlForPrivateStatic() *CacheControlOptions {
	return &CacheControlOptions{
		MaxAge:      setting.StaticCacheTime,
		NoTransform: true,
	}
}

// HandleGenericETagCache handles ETag-based caching for a HTTP request.
// It returns true if the request was handled.
func HandleGenericETagCache(req *http.Request, w http.ResponseWriter, etag string) (handled bool) {
	if len(etag) > 0 {
		w.Header().Set("Etag", etag)
		if checkIfNoneMatchIsValid(req, etag) {
			w.WriteHeader(http.StatusNotModified)
			return true
		}
	}
	// not sure whether it is a public content, so just use "private" (old behavior)
	SetCacheControlInHeader(w.Header(), CacheControlForPrivateStatic())
	return false
}

// checkIfNoneMatchIsValid tests if the header If-None-Match matches the ETag
func checkIfNoneMatchIsValid(req *http.Request, etag string) bool {
	ifNoneMatch := req.Header.Get("If-None-Match")
	if len(ifNoneMatch) > 0 {
		for item := range strings.SplitSeq(ifNoneMatch, ",") {
			item = strings.TrimPrefix(strings.TrimSpace(item), "W/") // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag#directives
			if item == etag {
				return true
			}
		}
	}
	return false
}

// HandleGenericETagTimeCache handles ETag-based caching with Last-Modified caching for a HTTP request.
// It returns true if the request was handled.
func HandleGenericETagTimeCache(req *http.Request, w http.ResponseWriter, etag string, lastModified *time.Time) (handled bool) {
	if len(etag) > 0 {
		w.Header().Set("Etag", etag)
	}
	if lastModified != nil && !lastModified.IsZero() {
		// http.TimeFormat required a UTC time, refer to https://pkg.go.dev/net/http#TimeFormat
		w.Header().Set("Last-Modified", lastModified.UTC().Format(http.TimeFormat))
	}

	if len(etag) > 0 {
		if checkIfNoneMatchIsValid(req, etag) {
			w.WriteHeader(http.StatusNotModified)
			return true
		}
	}
	if lastModified != nil && !lastModified.IsZero() {
		ifModifiedSince := req.Header.Get("If-Modified-Since")
		if ifModifiedSince != "" {
			t, err := time.Parse(http.TimeFormat, ifModifiedSince)
			if err == nil && lastModified.Unix() <= t.Unix() {
				w.WriteHeader(http.StatusNotModified)
				return true
			}
		}
	}

	// not sure whether it is a public content, so just use "private" (old behavior)
	SetCacheControlInHeader(w.Header(), CacheControlForPrivateStatic())
	return false
}
