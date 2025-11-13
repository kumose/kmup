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

package cmd

import (
	"net"
	"net/http"
	"net/http/fcgi"
	"strings"

	"github.com/kumose/kmup/modules/graceful"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"
)

func runHTTP(network, listenAddr, name string, m http.Handler, useProxyProtocol bool) error {
	return graceful.HTTPListenAndServe(network, listenAddr, name, m, useProxyProtocol)
}

// NoHTTPRedirector tells our cleanup routine that we will not be using a fallback http redirector
func NoHTTPRedirector() {
	graceful.GetManager().InformCleanup()
}

// NoInstallListener tells our cleanup routine that we will not be using a possibly provided listener
// for our install HTTP/HTTPS service
func NoInstallListener() {
	graceful.GetManager().InformCleanup()
}

func runFCGI(network, listenAddr, name string, m http.Handler, useProxyProtocol bool) error {
	// This needs to handle stdin as fcgi point
	fcgiServer := graceful.NewServer(network, listenAddr, name)

	err := fcgiServer.ListenAndServe(func(listener net.Listener) error {
		return fcgi.Serve(listener, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if setting.AppSubURL != "" {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, setting.AppSubURL)
			}
			m.ServeHTTP(resp, req)
		}))
	}, useProxyProtocol)
	if err != nil {
		log.Fatal("Failed to start FCGI main server: %v", err)
	}
	log.Info("FCGI Listener: %s Closed", listenAddr)
	return err
}
