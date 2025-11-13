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

package hcaptcha

const (
	ErrMissingInputSecret           ErrorCode = "missing-input-secret"
	ErrInvalidInputSecret           ErrorCode = "invalid-input-secret"
	ErrMissingInputResponse         ErrorCode = "missing-input-response"
	ErrInvalidInputResponse         ErrorCode = "invalid-input-response"
	ErrBadRequest                   ErrorCode = "bad-request"
	ErrInvalidOrAlreadySeenResponse ErrorCode = "invalid-or-already-seen-response"
	ErrNotUsingDummyPasscode        ErrorCode = "not-using-dummy-passcode"
	ErrSitekeySecretMismatch        ErrorCode = "sitekey-secret-mismatch"
)

// ErrorCode is any possible error from hCaptcha
type ErrorCode string

// String fulfills the Stringer interface
func (err ErrorCode) String() string {
	switch err {
	case ErrMissingInputSecret:
		return "Your secret key is missing."
	case ErrInvalidInputSecret:
		return "Your secret key is invalid or malformed."
	case ErrMissingInputResponse:
		return "The response parameter (verification token) is missing."
	case ErrInvalidInputResponse:
		return "The response parameter (verification token) is invalid or malformed."
	case ErrBadRequest:
		return "The request is invalid or malformed."
	case ErrInvalidOrAlreadySeenResponse:
		return "The response parameter has already been checked, or has another issue."
	case ErrNotUsingDummyPasscode:
		return "You have used a testing sitekey but have not used its matching secret."
	case ErrSitekeySecretMismatch:
		return "The sitekey is not registered with the provided secret."
	default:
		return ""
	}
}

// Error fulfills the error interface
func (err ErrorCode) Error() string {
	return err.String()
}
