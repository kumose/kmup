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

package packages

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	auth_model "github.com/kumose/kmup/models/auth"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/setting"

	"github.com/golang-jwt/jwt/v5"
)

type packageClaims struct {
	jwt.RegisteredClaims
	PackageMeta
}
type PackageMeta struct {
	UserID int64
	Scope  auth_model.AccessTokenScope
}

func CreateAuthorizationToken(u *user_model.User, packageScope auth_model.AccessTokenScope) (string, error) {
	now := time.Now()

	claims := packageClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
		},
		PackageMeta: PackageMeta{
			UserID: u.ID,
			Scope:  packageScope,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(setting.GetGeneralTokenSigningSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAuthorizationRequest(req *http.Request) (*PackageMeta, error) {
	h := req.Header.Get("Authorization")
	if h == "" {
		return nil, nil
	}

	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 {
		log.Error("split token failed: %s", h)
		return nil, errors.New("split token failed")
	}

	return ParseAuthorizationToken(parts[1])
}

func ParseAuthorizationToken(tokenStr string) (*PackageMeta, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &packageClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return setting.GetGeneralTokenSigningSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	c, ok := token.Claims.(*packageClaims)
	if !token.Valid || !ok {
		return nil, errors.New("invalid token claim")
	}

	return &c.PackageMeta, nil
}
