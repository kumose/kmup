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

package backend

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kumose/kmup/modules/json"
	lfslock "github.com/kumose/kmup/modules/structs"

	"github.com/charmbracelet/git-lfs-transfer/transfer"
)

var _ transfer.LockBackend = &kmupLockBackend{}

type kmupLockBackend struct {
	ctx          context.Context
	g            *KmupBackend
	server       *url.URL
	authToken    string
	internalAuth string
	logger       transfer.Logger
}

func newKmupLockBackend(g *KmupBackend) transfer.LockBackend {
	server := g.server.JoinPath("locks")
	return &kmupLockBackend{ctx: g.ctx, g: g, server: server, authToken: g.authToken, internalAuth: g.internalAuth, logger: g.logger}
}

// Create implements transfer.LockBackend
func (g *kmupLockBackend) Create(path, refname string) (transfer.Lock, error) {
	reqBody := lfslock.LFSLockRequest{Path: path}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		g.logger.Log("json marshal error", err)
		return nil, err
	}
	headers := map[string]string{
		headerAuthorization:    g.authToken,
		headerKmupInternalAuth: g.internalAuth,
		headerAccept:           mimeGitLFS,
		headerContentType:      mimeGitLFS,
	}
	req := newInternalRequestLFS(g.ctx, g.server.String(), http.MethodPost, headers, bodyBytes)
	resp, err := req.Response()
	if err != nil {
		g.logger.Log("http request error", err)
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		g.logger.Log("http read error", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		g.logger.Log("http statuscode error", resp.StatusCode, statusCodeToErr(resp.StatusCode))
		return nil, statusCodeToErr(resp.StatusCode)
	}
	var respBody lfslock.LFSLockResponse
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		g.logger.Log("json umarshal error", err)
		return nil, err
	}

	if respBody.Lock == nil {
		g.logger.Log("api returned nil lock")
		return nil, errors.New("api returned nil lock")
	}
	respLock := respBody.Lock
	owner := userUnknown
	if respLock.Owner != nil {
		owner = respLock.Owner.Name
	}
	lock := newKmupLock(g, respLock.ID, respLock.Path, respLock.LockedAt, owner)
	return lock, nil
}

// Unlock implements transfer.LockBackend
func (g *kmupLockBackend) Unlock(lock transfer.Lock) error {
	reqBody := lfslock.LFSLockDeleteRequest{}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		g.logger.Log("json marshal error", err)
		return err
	}
	headers := map[string]string{
		headerAuthorization:    g.authToken,
		headerKmupInternalAuth: g.internalAuth,
		headerAccept:           mimeGitLFS,
		headerContentType:      mimeGitLFS,
	}
	req := newInternalRequestLFS(g.ctx, g.server.JoinPath(lock.ID(), "unlock").String(), http.MethodPost, headers, bodyBytes)
	resp, err := req.Response()
	if err != nil {
		g.logger.Log("http request error", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		g.logger.Log("http statuscode error", resp.StatusCode, statusCodeToErr(resp.StatusCode))
		return statusCodeToErr(resp.StatusCode)
	}
	// no need to read response

	return nil
}

// FromPath implements transfer.LockBackend
func (g *kmupLockBackend) FromPath(path string) (transfer.Lock, error) {
	v := url.Values{
		argPath: []string{path},
	}

	respLocks, _, err := g.queryLocks(v)
	if err != nil {
		return nil, err
	}

	if len(respLocks) == 0 {
		return nil, transfer.ErrNotFound
	}
	return respLocks[0], nil
}

// FromID implements transfer.LockBackend
func (g *kmupLockBackend) FromID(id string) (transfer.Lock, error) {
	v := url.Values{
		argID: []string{id},
	}

	respLocks, _, err := g.queryLocks(v)
	if err != nil {
		return nil, err
	}

	if len(respLocks) == 0 {
		return nil, transfer.ErrNotFound
	}
	return respLocks[0], nil
}

// Range implements transfer.LockBackend
func (g *kmupLockBackend) Range(cursor string, limit int, iter func(transfer.Lock) error) (string, error) {
	v := url.Values{
		argLimit: []string{strconv.FormatInt(int64(limit), 10)},
	}
	if cursor != "" {
		v[argCursor] = []string{cursor}
	}

	respLocks, cursor, err := g.queryLocks(v)
	if err != nil {
		return "", err
	}

	for _, lock := range respLocks {
		err := iter(lock)
		if err != nil {
			return "", err
		}
	}
	return cursor, nil
}

func (g *kmupLockBackend) queryLocks(v url.Values) ([]transfer.Lock, string, error) {
	serverURLWithQuery := g.server.JoinPath() // get a copy
	serverURLWithQuery.RawQuery = v.Encode()
	headers := map[string]string{
		headerAuthorization:    g.authToken,
		headerKmupInternalAuth: g.internalAuth,
		headerAccept:           mimeGitLFS,
		headerContentType:      mimeGitLFS,
	}
	req := newInternalRequestLFS(g.ctx, serverURLWithQuery.String(), http.MethodGet, headers, nil)
	resp, err := req.Response()
	if err != nil {
		g.logger.Log("http request error", err)
		return nil, "", err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		g.logger.Log("http read error", err)
		return nil, "", err
	}
	if resp.StatusCode != http.StatusOK {
		g.logger.Log("http statuscode error", resp.StatusCode, statusCodeToErr(resp.StatusCode))
		return nil, "", statusCodeToErr(resp.StatusCode)
	}
	var respBody lfslock.LFSLockList
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		g.logger.Log("json umarshal error", err)
		return nil, "", err
	}

	respLocks := make([]transfer.Lock, 0, len(respBody.Locks))
	for _, respLock := range respBody.Locks {
		owner := userUnknown
		if respLock.Owner != nil {
			owner = respLock.Owner.Name
		}
		lock := newKmupLock(g, respLock.ID, respLock.Path, respLock.LockedAt, owner)
		respLocks = append(respLocks, lock)
	}
	return respLocks, respBody.Next, nil
}

var _ transfer.Lock = &kmupLock{}

type kmupLock struct {
	g        *kmupLockBackend
	id       string
	path     string
	lockedAt time.Time
	owner    string
}

func newKmupLock(g *kmupLockBackend, id, path string, lockedAt time.Time, owner string) transfer.Lock {
	return &kmupLock{g: g, id: id, path: path, lockedAt: lockedAt, owner: owner}
}

// Unlock implements transfer.Lock
func (g *kmupLock) Unlock() error {
	return g.g.Unlock(g)
}

// ID implements transfer.Lock
func (g *kmupLock) ID() string {
	return g.id
}

// Path implements transfer.Lock
func (g *kmupLock) Path() string {
	return g.path
}

// FormattedTimestamp implements transfer.Lock
func (g *kmupLock) FormattedTimestamp() string {
	return g.lockedAt.UTC().Format(time.RFC3339)
}

// OwnerName implements transfer.Lock
func (g *kmupLock) OwnerName() string {
	return g.owner
}

func (g *kmupLock) CurrentUser() (string, error) {
	return userSelf, nil
}

// AsLockSpec implements transfer.Lock
func (g *kmupLock) AsLockSpec(ownerID bool) ([]string, error) {
	msgs := []string{
		"lock " + g.ID(),
		fmt.Sprintf("path %s %s", g.ID(), g.Path()),
		fmt.Sprintf("locked-at %s %s", g.ID(), g.FormattedTimestamp()),
		fmt.Sprintf("ownername %s %s", g.ID(), g.OwnerName()),
	}
	if ownerID {
		user, err := g.CurrentUser()
		if err != nil {
			return nil, fmt.Errorf("error getting current user: %w", err)
		}
		who := "theirs"
		if user == g.OwnerName() {
			who = "ours"
		}
		msgs = append(msgs, fmt.Sprintf("owner %s %s", g.ID(), who))
	}
	return msgs, nil
}

// AsArguments implements transfer.Lock
func (g *kmupLock) AsArguments() []string {
	return []string{
		"id=" + g.ID(),
		"path=" + g.Path(),
		"locked-at=" + g.FormattedTimestamp(),
		"ownername=" + g.OwnerName(),
	}
}
