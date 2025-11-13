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

package lfs

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/util"
)

// DetermineEndpoint determines an endpoint from the clone url or uses the specified LFS url.
func DetermineEndpoint(cloneurl, lfsurl string) *url.URL {
	if len(lfsurl) > 0 {
		return endpointFromURL(lfsurl)
	}
	return endpointFromCloneURL(cloneurl)
}

func endpointFromCloneURL(rawurl string) *url.URL {
	ep := endpointFromURL(rawurl)
	if ep == nil {
		return ep
	}

	ep.Path = strings.TrimSuffix(ep.Path, "/")

	if ep.Scheme == "file" {
		return ep
	}

	if path.Ext(ep.Path) == ".git" {
		ep.Path += "/info/lfs"
	} else {
		ep.Path += ".git/info/lfs"
	}

	return ep
}

func endpointFromURL(rawurl string) *url.URL {
	if strings.HasPrefix(rawurl, "/") {
		return endpointFromLocalPath(rawurl)
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		log.Error("lfs.endpointFromUrl: %v", err)
		return nil
	}

	switch u.Scheme {
	case "http", "https":
		return u
	case "git":
		u.Scheme = "https"
		return u
	case "file":
		return u
	default:
		if _, err := os.Stat(rawurl); err == nil {
			return endpointFromLocalPath(rawurl)
		}

		log.Error("lfs.endpointFromUrl: unknown url")
		return nil
	}
}

func endpointFromLocalPath(path string) *url.URL {
	var slash string
	if abs, err := filepath.Abs(path); err == nil {
		if !strings.HasPrefix(abs, "/") {
			slash = "/"
		}
		path = abs
	}

	var gitpath string
	if filepath.Base(path) == ".git" {
		gitpath = path
		path = filepath.Dir(path)
	} else {
		gitpath = filepath.Join(path, ".git")
	}

	if _, err := os.Stat(gitpath); err == nil {
		path = gitpath
	} else if _, err := os.Stat(path); err != nil {
		return nil
	}

	path = "file://" + slash + util.PathEscapeSegments(filepath.ToSlash(path))

	u, _ := url.Parse(path)

	return u
}
