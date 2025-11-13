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

package v1_13

import (
	"github.com/kumose/kmup/modules/timeutil"

	"xorm.io/xorm"
)

func AddTimeStamps(x *xorm.Engine) error {
	// this will add timestamps where it is useful to have

	// Star represents a starred repo by an user.
	type Star struct {
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	}
	if err := x.Sync(new(Star)); err != nil {
		return err
	}

	// Label represents a label of repository for issues.
	type Label struct {
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}
	if err := x.Sync(new(Label)); err != nil {
		return err
	}

	// Follow represents relations of user and their followers.
	type Follow struct {
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	}
	if err := x.Sync(new(Follow)); err != nil {
		return err
	}

	// Watch is connection request for receiving repository notification.
	type Watch struct {
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}
	if err := x.Sync(new(Watch)); err != nil {
		return err
	}

	// Collaboration represent the relation between an individual and a repository.
	type Collaboration struct {
		CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
		UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	}
	return x.Sync(new(Collaboration))
}
