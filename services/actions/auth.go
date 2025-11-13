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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	"github.com/golang-jwt/jwt/v5"
)

type actionsClaims struct {
	jwt.RegisteredClaims
	Scp    string `json:"scp"`
	TaskID int64
	RunID  int64
	JobID  int64
	Ac     string `json:"ac"`
}

type actionsCacheScope struct {
	Scope      string
	Permission actionsCachePermission
}

type actionsCachePermission int

const (
	actionsCachePermissionRead = 1 << iota
	actionsCachePermissionWrite
)

func CreateAuthorizationToken(taskID, runID, jobID int64) (string, error) {
	now := time.Now()

	ac, err := json.Marshal(&[]actionsCacheScope{
		{
			Scope:      "",
			Permission: actionsCachePermissionWrite,
		},
	})
	if err != nil {
		return "", err
	}

	claims := actionsClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(1*time.Hour + setting.Actions.EndlessTaskTimeout)),
			NotBefore: jwt.NewNumericDate(now),
		},
		Scp:    fmt.Sprintf("Actions.Results:%d:%d", runID, jobID),
		Ac:     string(ac),
		TaskID: taskID,
		RunID:  runID,
		JobID:  jobID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(setting.GetGeneralTokenSigningSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAuthorizationToken(req *http.Request) (int64, error) {
	h := req.Header.Get("Authorization")
	if h == "" {
		return 0, nil
	}

	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 {
		log.Error("split token failed: %s", h)
		return 0, errors.New("split token failed")
	}

	return TokenToTaskID(parts[1])
}

// TokenToTaskID returns the TaskID associated with the provided JWT token
func TokenToTaskID(token string) (int64, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &actionsClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return setting.GetGeneralTokenSigningSecret(), nil
	})
	if err != nil {
		return 0, err
	}

	c, ok := parsedToken.Claims.(*actionsClaims)
	if !parsedToken.Valid || !ok {
		return 0, errors.New("invalid token claim")
	}

	return c.TaskID, nil
}
