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

package v1_21

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kumose/kmup/modules/git"
	giturl "github.com/kumose/kmup/modules/git/url"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"

	"xorm.io/xorm"
)

func AddRemoteAddressToMirrors(x *xorm.Engine) error {
	type Mirror struct {
		RemoteAddress string `xorm:"VARCHAR(2048)"`
	}

	type PushMirror struct {
		RemoteAddress string `xorm:"VARCHAR(2048)"`
	}

	if err := x.Sync(new(Mirror), new(PushMirror)); err != nil {
		return err
	}

	if err := migratePullMirrors(x); err != nil {
		return err
	}

	return migratePushMirrors(x)
}

func migratePullMirrors(x *xorm.Engine) error {
	type Mirror struct {
		ID            int64  `xorm:"pk autoincr"`
		RepoID        int64  `xorm:"INDEX"`
		RemoteAddress string `xorm:"VARCHAR(2048)"`
		RepoOwner     string
		RepoName      string
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	limit := setting.Database.IterateBufferSize
	if limit <= 0 {
		limit = 50
	}

	start := 0

	for {
		var mirrors []Mirror
		if err := sess.Select("mirror.id, mirror.repo_id, mirror.remote_address, repository.owner_name as repo_owner, repository.name as repo_name").
			Join("INNER", "repository", "repository.id = mirror.repo_id").
			Limit(limit, start).Find(&mirrors); err != nil {
			return err
		}

		if len(mirrors) == 0 {
			break
		}
		start += len(mirrors)

		for _, m := range mirrors {
			remoteAddress, err := getRemoteAddress(m.RepoOwner, m.RepoName, "origin")
			if err != nil {
				return err
			}

			m.RemoteAddress = remoteAddress

			if _, err = sess.ID(m.ID).Cols("remote_address").Update(m); err != nil {
				return err
			}
		}

		if start%1000 == 0 { // avoid a too big transaction
			if err := sess.Commit(); err != nil {
				return err
			}
			if err := sess.Begin(); err != nil {
				return err
			}
		}
	}

	return sess.Commit()
}

func migratePushMirrors(x *xorm.Engine) error {
	type PushMirror struct {
		ID            int64 `xorm:"pk autoincr"`
		RepoID        int64 `xorm:"INDEX"`
		RemoteName    string
		RemoteAddress string `xorm:"VARCHAR(2048)"`
		RepoOwner     string
		RepoName      string
	}

	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	limit := setting.Database.IterateBufferSize
	if limit <= 0 {
		limit = 50
	}

	start := 0

	for {
		var mirrors []PushMirror
		if err := sess.Select("push_mirror.id, push_mirror.repo_id, push_mirror.remote_name, push_mirror.remote_address, repository.owner_name as repo_owner, repository.name as repo_name").
			Join("INNER", "repository", "repository.id = push_mirror.repo_id").
			Limit(limit, start).Find(&mirrors); err != nil {
			return err
		}

		if len(mirrors) == 0 {
			break
		}
		start += len(mirrors)

		for _, m := range mirrors {
			remoteAddress, err := getRemoteAddress(m.RepoOwner, m.RepoName, m.RemoteName)
			if err != nil {
				return err
			}

			m.RemoteAddress = remoteAddress

			if _, err = sess.ID(m.ID).Cols("remote_address").Update(m); err != nil {
				return err
			}
		}

		if start%1000 == 0 { // avoid a too big transaction
			if err := sess.Commit(); err != nil {
				return err
			}
			if err := sess.Begin(); err != nil {
				return err
			}
		}
	}

	return sess.Commit()
}

func getRemoteAddress(ownerName, repoName, remoteName string) (string, error) {
	repoPath := filepath.Join(setting.RepoRootPath, strings.ToLower(ownerName), strings.ToLower(repoName)+".git")
	if exist, _ := util.IsExist(repoPath); !exist {
		return "", nil
	}
	remoteURL, err := git.GetRemoteAddress(context.Background(), repoPath, remoteName)
	if err != nil {
		return "", fmt.Errorf("get remote %s's address of %s/%s failed: %v", remoteName, ownerName, repoName, err)
	}

	u, err := giturl.ParseGitURL(remoteURL)
	if err != nil {
		return "", err
	}
	u.User = nil

	return u.String(), nil
}
