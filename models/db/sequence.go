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

package db

import (
	"context"
	"fmt"
	"regexp"

	"github.com/kumose/kmup/modules/setting"
)

// CountBadSequences looks for broken sequences from recreate-table mistakes
func CountBadSequences(_ context.Context) (int64, error) {
	if !setting.Database.Type.IsPostgreSQL() {
		return 0, nil
	}

	sess := xormEngine.NewSession()
	defer sess.Close()

	var sequences []string
	schema := xormEngine.Dialect().URI().Schema

	sess.Engine().SetSchema("")
	if err := sess.Table("information_schema.sequences").Cols("sequence_name").Where("sequence_name LIKE 'tmp_recreate__%_id_seq%' AND sequence_catalog = ?", setting.Database.Name).Find(&sequences); err != nil {
		return 0, err
	}
	sess.Engine().SetSchema(schema)

	return int64(len(sequences)), nil
}

// FixBadSequences fixes for broken sequences from recreate-table mistakes
func FixBadSequences(_ context.Context) error {
	if !setting.Database.Type.IsPostgreSQL() {
		return nil
	}

	sess := xormEngine.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	var sequences []string
	schema := sess.Engine().Dialect().URI().Schema

	sess.Engine().SetSchema("")
	if err := sess.Table("information_schema.sequences").Cols("sequence_name").Where("sequence_name LIKE 'tmp_recreate__%_id_seq%' AND sequence_catalog = ?", setting.Database.Name).Find(&sequences); err != nil {
		return err
	}
	sess.Engine().SetSchema(schema)

	sequenceRegexp := regexp.MustCompile(`tmp_recreate__(\w+)_id_seq.*`)

	for _, sequence := range sequences {
		tableName := sequenceRegexp.FindStringSubmatch(sequence)[1]
		newSequenceName := tableName + "_id_seq"
		if _, err := sess.Exec(fmt.Sprintf("ALTER SEQUENCE `%s` RENAME TO `%s`", sequence, newSequenceName)); err != nil {
			return err
		}
		if _, err := sess.Exec(fmt.Sprintf("SELECT setval('%s', COALESCE((SELECT MAX(id)+1 FROM `%s`), 1), false)", newSequenceName, tableName)); err != nil {
			return err
		}
	}

	return sess.Commit()
}
