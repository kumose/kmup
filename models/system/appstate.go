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

package system

import (
	"context"

	"github.com/kumose/kmup/models/db"
)

// AppState represents a state record in database
// if one day we would make Kmup run as a cluster,
// we can introduce a new field `Scope` here to store different states for different nodes
type AppState struct {
	ID       string `xorm:"pk varchar(200)"`
	Revision int64
	Content  string `xorm:"LONGTEXT"`
}

func init() {
	db.RegisterModel(new(AppState))
}

// SaveAppStateContent saves the app state item to database
func SaveAppStateContent(ctx context.Context, key, content string) error {
	return db.WithTx(ctx, func(ctx context.Context) error {
		eng := db.GetEngine(ctx)
		// try to update existing row
		res, err := eng.Exec("UPDATE app_state SET revision=revision+1, content=? WHERE id=?", content, key)
		if err != nil {
			return err
		}
		rows, _ := res.RowsAffected()
		if rows != 0 {
			// the existing row is updated, so we can return
			return nil
		}
		// if no existing row, insert a new row
		_, err = eng.Insert(&AppState{ID: key, Content: content})
		return err
	})
}

// GetAppStateContent gets an app state from database
func GetAppStateContent(ctx context.Context, key string) (content string, err error) {
	e := db.GetEngine(ctx)
	appState := &AppState{ID: key}
	has, err := e.Get(appState)
	if err != nil {
		return "", err
	} else if !has {
		return "", nil
	}
	return appState.Content, nil
}
