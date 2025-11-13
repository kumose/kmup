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

package actions

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/modules/timeutil"
	"github.com/kumose/kmup/modules/util"
)

func generateSaltedToken() (string, string, string, string, error) {
	salt, err := util.CryptoRandomString(10)
	if err != nil {
		return "", "", "", "", err
	}
	buf, err := util.CryptoRandomBytes(20)
	if err != nil {
		return "", "", "", "", err
	}
	token := hex.EncodeToString(buf)
	hash := auth_model.HashToken(token, salt)
	return token, salt, hash, token[len(token)-8:], nil
}

/*
LogIndexes is the index for mapping log line number to buffer offset.
Because it uses varint encoding, it is impossible to predict its size.
But we can make a simple estimate with an assumption that each log line has 200 byte, then:
| lines     | file size           | index size         |
|-----------|---------------------|--------------------|
| 100       | 20 KiB(20000)       | 258 B(258)         |
| 1000      | 195 KiB(200000)     | 2.9 KiB(2958)      |
| 10000     | 1.9 MiB(2000000)    | 34 KiB(34715)      |
| 100000    | 19 MiB(20000000)    | 386 KiB(394715)    |
| 1000000   | 191 MiB(200000000)  | 4.1 MiB(4323626)   |
| 10000000  | 1.9 GiB(2000000000) | 47 MiB(49323626)   |
| 100000000 | 19 GiB(20000000000) | 490 MiB(513424280) |
*/
type LogIndexes []int64

func (indexes *LogIndexes) FromDB(b []byte) error {
	reader := bytes.NewReader(b)
	for {
		v, err := binary.ReadVarint(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("binary ReadVarint: %w", err)
		}
		*indexes = append(*indexes, v)
	}
}

func (indexes *LogIndexes) ToDB() ([]byte, error) {
	buf, i := make([]byte, binary.MaxVarintLen64*len(*indexes)), 0
	for _, v := range *indexes {
		n := binary.PutVarint(buf[i:], v)
		i += n
	}
	return buf[:i], nil
}

var timeSince = time.Since

func calculateDuration(started, stopped timeutil.TimeStamp, status Status) time.Duration {
	if started == 0 {
		return 0
	}
	s := started.AsTime()
	if status.IsDone() {
		return stopped.AsTime().Sub(s)
	}
	return timeSince(s).Truncate(time.Second)
}

// best effort function to convert an action schedule to action run, to be used in GenerateKmupContext
func (s *ActionSchedule) ToActionRun() *ActionRun {
	return &ActionRun{
		Title:         s.Title,
		RepoID:        s.RepoID,
		Repo:          s.Repo,
		OwnerID:       s.OwnerID,
		WorkflowID:    s.WorkflowID,
		TriggerUserID: s.TriggerUserID,
		TriggerUser:   s.TriggerUser,
		Ref:           s.Ref,
		CommitSHA:     s.CommitSHA,
		Event:         s.Event,
		EventPayload:  s.EventPayload,
		Created:       s.Created,
		Updated:       s.Updated,
	}
}
