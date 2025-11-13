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

// Package private contains all internal routes. The package name "internal" isn't usable because Golang reserves it for disabling cross-package usage.
package private

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/private"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/web"
	"github.com/kumose/kmup/routers/common"
	"github.com/kumose/kmup/services/context"

	"github.com/kumose-go/chi/binding"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

func authInternal(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if setting.InternalToken == "" {
			log.Warn(`The INTERNAL_TOKEN setting is missing from the configuration file: %q, internal API can't work.`, setting.CustomConf)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		tokens := req.Header.Get("X-Kmup-Internal-Auth") // TODO: use something like JWT or HMAC to avoid passing the token in the clear
		after, found := strings.CutPrefix(tokens, "Bearer ")
		authSucceeded := found && subtle.ConstantTimeCompare([]byte(after), []byte(setting.InternalToken)) == 1
		if !authSucceeded {
			log.Debug("Forbidden attempt to access internal url: Authorization header: %s", tokens)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, req)
	})
}

// bind binding an obj to a handler
func bind[T any](_ T) any {
	return func(ctx *context.PrivateContext) {
		theObj := new(T) // create a new form obj for every request but not use obj directly
		binding.Bind(ctx.Req, theObj)
		web.SetForm(ctx, theObj)
	}
}

// Routes registers all internal APIs routes to web application.
// These APIs will be invoked by internal commands for example `kmup serv` and etc.
func Routes() *web.Router {
	r := web.NewRouter()
	r.Use(context.PrivateContexter())
	r.Use(authInternal)
	// Log the real ip address of the request from SSH is really helpful for diagnosing sometimes.
	// Since internal API will be sent only from Kmup sub commands and it's under control (checked by InternalToken), we can trust the headers.
	r.Use(chi_middleware.RealIP)

	r.Post("/ssh/authorized_keys", AuthorizedPublicKeyByContent)
	r.Post("/ssh/{id}/update/{repoid}", UpdatePublicKeyInRepo)
	r.Post("/ssh/log", bind(private.SSHLogOption{}), SSHLog)
	r.Post("/hook/pre-receive/{owner}/{repo}", RepoAssignment, bind(private.HookOptions{}), HookPreReceive)
	r.Post("/hook/post-receive/{owner}/{repo}", context.OverrideContext(), bind(private.HookOptions{}), HookPostReceive)
	r.Post("/hook/proc-receive/{owner}/{repo}", context.OverrideContext(), RepoAssignment, bind(private.HookOptions{}), HookProcReceive)
	r.Post("/hook/set-default-branch/{owner}/{repo}/{branch}", RepoAssignment, SetDefaultBranch)
	r.Get("/serv/none/{keyid}", ServNoCommand)
	r.Get("/serv/command/{keyid}/{owner}/{repo}", ServCommand)
	r.Post("/manager/shutdown", Shutdown)
	r.Post("/manager/restart", Restart)
	r.Post("/manager/reload-templates", ReloadTemplates)
	r.Post("/manager/flush-queues", bind(private.FlushOptions{}), FlushQueues)
	r.Post("/manager/pause-logging", PauseLogging)
	r.Post("/manager/resume-logging", ResumeLogging)
	r.Post("/manager/release-and-reopen-logging", ReleaseReopenLogging)
	r.Post("/manager/set-log-sql", SetLogSQL)
	r.Post("/manager/add-logger", bind(private.LoggerOptions{}), AddLogger)
	r.Post("/manager/remove-logger/{logger}/{writer}", RemoveLogger)
	r.Get("/manager/processes", Processes)
	r.Post("/mail/send", SendEmail)
	r.Post("/restore_repo", RestoreRepo)
	r.Post("/actions/generate_actions_runner_token", GenerateActionsRunnerToken)

	r.Group("/repo", func() {
		// FIXME: it is not right to use context.Contexter here because all routes here should use PrivateContext
		// Fortunately, the LFS handlers are able to handle requests without a complete web context
		common.AddOwnerRepoGitLFSRoutes(r, func(ctx *context.PrivateContext) {
			webContext := &context.Context{Base: ctx.Base}         // see above, it shouldn't manually construct the web context
			ctx.SetContextValue(context.WebContextKey, webContext) // FIXME: this is not ideal but no other way at the moment
		})
	})

	return r
}
