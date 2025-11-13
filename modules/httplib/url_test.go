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

package httplib

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestIsRelativeURL(t *testing.T) {
	defer test.MockVariableValue(&setting.AppURL, "http://localhost:3000/sub/")()
	defer test.MockVariableValue(&setting.AppSubURL, "/sub")()
	rel := []string{
		"",
		"foo",
		"/",
		"/foo?k=%20#abc",
	}
	for _, s := range rel {
		assert.True(t, IsRelativeURL(s), "rel = %q", s)
	}
	abs := []string{
		"//",
		"\\\\",
		"/\\",
		"\\/",
		"mailto:a@b.com",
		"https://test.com",
	}
	for _, s := range abs {
		assert.False(t, IsRelativeURL(s), "abs = %q", s)
	}
}

func TestGuessCurrentHostURL(t *testing.T) {
	defer test.MockVariableValue(&setting.AppURL, "http://cfg-host/sub/")()
	defer test.MockVariableValue(&setting.AppSubURL, "/sub")()
	headersWithProto := http.Header{"X-Forwarded-Proto": {"https"}}

	t.Run("Legacy", func(t *testing.T) {
		defer test.MockVariableValue(&setting.PublicURLDetection, setting.PublicURLLegacy)()

		assert.Equal(t, "http://cfg-host", GuessCurrentHostURL(t.Context()))

		// legacy: "Host" is not used when there is no "X-Forwarded-Proto" header
		ctx := context.WithValue(t.Context(), RequestContextKey, &http.Request{Host: "req-host:3000"})
		assert.Equal(t, "http://cfg-host", GuessCurrentHostURL(ctx))

		// if "X-Forwarded-Proto" exists, then use it and "Host" header
		ctx = context.WithValue(t.Context(), RequestContextKey, &http.Request{Host: "req-host:3000", Header: headersWithProto})
		assert.Equal(t, "https://req-host:3000", GuessCurrentHostURL(ctx))
	})

	t.Run("Auto", func(t *testing.T) {
		defer test.MockVariableValue(&setting.PublicURLDetection, setting.PublicURLAuto)()

		assert.Equal(t, "http://cfg-host", GuessCurrentHostURL(t.Context()))

		// auto: always use "Host" header, the scheme is determined by "X-Forwarded-Proto" header, or TLS config if no "X-Forwarded-Proto" header
		ctx := context.WithValue(t.Context(), RequestContextKey, &http.Request{Host: "req-host:3000"})
		assert.Equal(t, "http://req-host:3000", GuessCurrentHostURL(ctx))

		ctx = context.WithValue(t.Context(), RequestContextKey, &http.Request{Host: "req-host", TLS: &tls.ConnectionState{}})
		assert.Equal(t, "https://req-host", GuessCurrentHostURL(ctx))

		ctx = context.WithValue(t.Context(), RequestContextKey, &http.Request{Host: "req-host:3000", Header: headersWithProto})
		assert.Equal(t, "https://req-host:3000", GuessCurrentHostURL(ctx))
	})
}

func TestMakeAbsoluteURL(t *testing.T) {
	defer test.MockVariableValue(&setting.Protocol, "http")()
	defer test.MockVariableValue(&setting.AppURL, "http://cfg-host/sub/")()
	defer test.MockVariableValue(&setting.AppSubURL, "/sub")()

	ctx := t.Context()
	assert.Equal(t, "http://cfg-host/sub/", MakeAbsoluteURL(ctx, ""))
	assert.Equal(t, "http://cfg-host/foo", MakeAbsoluteURL(ctx, "foo"))
	assert.Equal(t, "http://cfg-host/foo", MakeAbsoluteURL(ctx, "/foo"))
	assert.Equal(t, "http://other/foo", MakeAbsoluteURL(ctx, "http://other/foo"))

	ctx = context.WithValue(ctx, RequestContextKey, &http.Request{
		Host: "user-host",
	})
	assert.Equal(t, "http://cfg-host/foo", MakeAbsoluteURL(ctx, "/foo"))

	ctx = context.WithValue(ctx, RequestContextKey, &http.Request{
		Host: "user-host",
		Header: map[string][]string{
			"X-Forwarded-Host": {"forwarded-host"},
		},
	})
	assert.Equal(t, "http://cfg-host/foo", MakeAbsoluteURL(ctx, "/foo"))

	ctx = context.WithValue(ctx, RequestContextKey, &http.Request{
		Host: "user-host",
		Header: map[string][]string{
			"X-Forwarded-Host":  {"forwarded-host"},
			"X-Forwarded-Proto": {"https"},
		},
	})
	assert.Equal(t, "https://user-host/foo", MakeAbsoluteURL(ctx, "/foo"))
}

func TestIsCurrentKmupSiteURL(t *testing.T) {
	defer test.MockVariableValue(&setting.AppURL, "http://localhost:3000/sub/")()
	defer test.MockVariableValue(&setting.AppSubURL, "/sub")()
	ctx := t.Context()
	good := []string{
		"?key=val",
		"/sub",
		"/sub/",
		"/sub/foo",
		"/sub/foo/",
		"http://localhost:3000/sub?key=val",
		"http://localhost:3000/sub/",
	}
	for _, s := range good {
		assert.True(t, IsCurrentKmupSiteURL(ctx, s), "good = %q", s)
	}
	bad := []string{
		".",
		"foo",
		"/",
		"//",
		"\\\\",
		"/foo",
		"http://localhost:3000/sub/..",
		"http://localhost:3000/other",
		"http://other/",
	}
	for _, s := range bad {
		assert.False(t, IsCurrentKmupSiteURL(ctx, s), "bad = %q", s)
	}

	setting.AppURL = "http://localhost:3000/"
	setting.AppSubURL = ""
	assert.False(t, IsCurrentKmupSiteURL(ctx, "//"))
	assert.False(t, IsCurrentKmupSiteURL(ctx, "\\\\"))
	assert.False(t, IsCurrentKmupSiteURL(ctx, "http://localhost"))
	assert.True(t, IsCurrentKmupSiteURL(ctx, "http://localhost:3000?key=val"))

	ctx = context.WithValue(ctx, RequestContextKey, &http.Request{
		Host: "user-host",
		Header: map[string][]string{
			"X-Forwarded-Host":  {"forwarded-host"},
			"X-Forwarded-Proto": {"https"},
		},
	})
	assert.True(t, IsCurrentKmupSiteURL(ctx, "http://localhost:3000"))
	assert.True(t, IsCurrentKmupSiteURL(ctx, "https://user-host"))
	assert.False(t, IsCurrentKmupSiteURL(ctx, "https://forwarded-host"))
}

func TestParseKmupSiteURL(t *testing.T) {
	defer test.MockVariableValue(&setting.AppURL, "http://localhost:3000/sub/")()
	defer test.MockVariableValue(&setting.AppSubURL, "/sub")()
	ctx := t.Context()
	tests := []struct {
		url string
		exp *KmupSiteURL
	}{
		{"http://localhost:3000/sub?k=v", &KmupSiteURL{RoutePath: ""}},
		{"http://localhost:3000/sub/", &KmupSiteURL{RoutePath: ""}},
		{"http://localhost:3000/sub/foo", &KmupSiteURL{RoutePath: "/foo"}},
		{"http://localhost:3000/sub/foo/bar", &KmupSiteURL{RoutePath: "/foo/bar", OwnerName: "foo", RepoName: "bar"}},
		{"http://localhost:3000/sub/foo/bar/", &KmupSiteURL{RoutePath: "/foo/bar", OwnerName: "foo", RepoName: "bar"}},
		{"http://localhost:3000/sub/attachments/bar", &KmupSiteURL{RoutePath: "/attachments/bar"}},
		{"http://localhost:3000/other", nil},
		{"http://other/", nil},
	}
	for _, test := range tests {
		su := ParseKmupSiteURL(ctx, test.url)
		assert.Equal(t, test.exp, su, "URL = %s", test.url)
	}
}
