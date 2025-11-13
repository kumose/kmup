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

package recaptcha

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/kumose/kmup/modules/json"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"
)

// Response is the structure of JSON returned from API
type Response struct {
	Success     bool        `json:"success"`
	ChallengeTS string      `json:"challenge_ts"`
	Hostname    string      `json:"hostname"`
	ErrorCodes  []ErrorCode `json:"error-codes"`
}

const apiURL = "api/siteverify"

// Verify calls Google Recaptcha API to verify token
func Verify(ctx context.Context, response string) (bool, error) {
	post := url.Values{
		"secret":   {setting.Service.RecaptchaSecret},
		"response": {response},
	}
	// Basically a copy of http.PostForm, but with a context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		util.URLJoin(setting.Service.RecaptchaURL, apiURL), strings.NewReader(post.Encode()))
	if err != nil {
		return false, fmt.Errorf("Failed to create CAPTCHA request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("Failed to send CAPTCHA response: %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("Failed to read CAPTCHA response: %s", err)
	}

	var jsonResponse Response
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return false, fmt.Errorf("Failed to parse CAPTCHA response: %s", err)
	}
	var respErr error
	if len(jsonResponse.ErrorCodes) > 0 {
		respErr = jsonResponse.ErrorCodes[0]
	}
	return jsonResponse.Success, respErr
}

// ErrorCode is a reCaptcha error
type ErrorCode string

// String fulfills the Stringer interface
func (e ErrorCode) String() string {
	switch e {
	case "missing-input-secret":
		return "The secret parameter is missing."
	case "invalid-input-secret":
		return "The secret parameter is invalid or malformed."
	case "missing-input-response":
		return "The response parameter is missing."
	case "invalid-input-response":
		return "The response parameter is invalid or malformed."
	case "bad-request":
		return "The request is invalid or malformed."
	case "timeout-or-duplicate":
		return "The response is no longer valid: either is too old or has been used previously."
	}
	return string(e)
}

// Error fulfills the error interface
func (e ErrorCode) Error() string {
	return e.String()
}
