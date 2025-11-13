// Copyright 2020 The Kmup Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_14

import (
	"context"

	"github.com/kumose/kmup/models/migrations/base"

	"xorm.io/xorm"
)

func ConvertTaskTypeToString(x *xorm.Engine) error {
	const (
		GOGS int = iota + 1
		SLACK
		KMUP
		DISCORD
		DINGTALK
		TELEGRAM
		MSTEAMS
		FEISHU
		MATRIX
		WECHATWORK
	)

	hookTaskTypes := map[int]string{
		KMUP:       "kmup",
		GOGS:       "gogs",
		SLACK:      "slack",
		DISCORD:    "discord",
		DINGTALK:   "dingtalk",
		TELEGRAM:   "telegram",
		MSTEAMS:    "msteams",
		FEISHU:     "feishu",
		MATRIX:     "matrix",
		WECHATWORK: "wechatwork",
	}

	type HookTask struct {
		Typ string `xorm:"VARCHAR(16) index"`
	}
	if err := x.Sync(new(HookTask)); err != nil {
		return err
	}

	// to keep the migration could be rerun
	exist, err := x.Dialect().IsColumnExist(x.DB(), context.Background(), "hook_task", "type")
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	for i, s := range hookTaskTypes {
		if _, err := x.Exec("UPDATE hook_task set typ = ? where `type`=?", s, i); err != nil {
			return err
		}
	}

	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if err := base.DropTableColumns(sess, "hook_task", "type"); err != nil {
		return err
	}

	return sess.Commit()
}
