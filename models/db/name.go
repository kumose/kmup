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
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/kumose/kmup/modules/util"
)

// ErrNameReserved represents a "reserved name" error.
type ErrNameReserved struct {
	Name string
}

// IsErrNameReserved checks if an error is a ErrNameReserved.
func IsErrNameReserved(err error) bool {
	_, ok := err.(ErrNameReserved)
	return ok
}

func (err ErrNameReserved) Error() string {
	return fmt.Sprintf("name is reserved [name: %s]", err.Name)
}

// Unwrap unwraps this as a ErrInvalid err
func (err ErrNameReserved) Unwrap() error {
	return util.ErrInvalidArgument
}

// ErrNamePatternNotAllowed represents a "pattern not allowed" error.
type ErrNamePatternNotAllowed struct {
	Pattern string
}

// IsErrNamePatternNotAllowed checks if an error is an ErrNamePatternNotAllowed.
func IsErrNamePatternNotAllowed(err error) bool {
	_, ok := err.(ErrNamePatternNotAllowed)
	return ok
}

func (err ErrNamePatternNotAllowed) Error() string {
	return fmt.Sprintf("name pattern is not allowed [pattern: %s]", err.Pattern)
}

// Unwrap unwraps this as a ErrInvalid err
func (err ErrNamePatternNotAllowed) Unwrap() error {
	return util.ErrInvalidArgument
}

// ErrNameCharsNotAllowed represents a "character not allowed in name" error.
type ErrNameCharsNotAllowed struct {
	Name string
}

// IsErrNameCharsNotAllowed checks if an error is an ErrNameCharsNotAllowed.
func IsErrNameCharsNotAllowed(err error) bool {
	_, ok := err.(ErrNameCharsNotAllowed)
	return ok
}

func (err ErrNameCharsNotAllowed) Error() string {
	return fmt.Sprintf("name is invalid [%s]: must be valid alpha or numeric or dash(-_) or dot characters", err.Name)
}

// Unwrap unwraps this as a ErrInvalid err
func (err ErrNameCharsNotAllowed) Unwrap() error {
	return util.ErrInvalidArgument
}

// IsUsableName checks if name is reserved or pattern of name is not allowed
// based on given reserved names and patterns.
// Names are exact match, patterns can be a prefix or suffix match with placeholder '*'.
func IsUsableName(reservedNames, reservedPatterns []string, name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return util.NewInvalidArgumentErrorf("name is empty")
	}

	if slices.Contains(reservedNames, name) {
		return ErrNameReserved{name}
	}

	for _, pat := range reservedPatterns {
		if pat[0] == '*' && strings.HasSuffix(name, pat[1:]) ||
			(pat[len(pat)-1] == '*' && strings.HasPrefix(name, pat[:len(pat)-1])) {
			return ErrNamePatternNotAllowed{pat}
		}
	}

	return nil
}
