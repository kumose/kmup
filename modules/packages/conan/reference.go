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

package conan

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kumose/kmup/modules/log"
	"github.com/kumose/kmup/modules/util"
)

const (
	// taken from https://github.com/conan-io/conan/blob/develop/conans/model/ref.py
	minChars = 2
	maxChars = 51

	// DefaultRevision if no revision is specified
	DefaultRevision = "0"
)

var (
	namePattern     = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9_][a-zA-Z0-9_\+\.-]{%d,%d}$`, minChars-1, maxChars-1))
	revisionPattern = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9]{1,%d}$`, maxChars))

	ErrValidation = util.NewInvalidArgumentErrorf("could not validate one or more reference fields")
)

// RecipeReference represents a recipe <Name>/<Version>@<User>/<Channel>#<Revision>
type RecipeReference struct {
	Name     string
	Version  string
	User     string
	Channel  string
	Revision string
}

func NewRecipeReference(name, version, user, channel, revision string) (*RecipeReference, error) {
	log.Trace("Conan Recipe: %s/%s(@%s/%s(#%s))", name, version, user, channel, revision)

	if user == "_" {
		user = ""
	}
	if channel == "_" {
		channel = ""
	}

	if (user != "" && channel == "") || (user == "" && channel != "") {
		return nil, ErrValidation
	}

	if !namePattern.MatchString(name) {
		return nil, ErrValidation
	}

	v := strings.TrimSpace(version)
	if v == "" {
		return nil, ErrValidation
	}
	if user != "" && !namePattern.MatchString(user) {
		return nil, ErrValidation
	}
	if channel != "" && !namePattern.MatchString(channel) {
		return nil, ErrValidation
	}
	if revision != "" && !revisionPattern.MatchString(revision) {
		return nil, ErrValidation
	}

	return &RecipeReference{name, v, user, channel, revision}, nil
}

func (r *RecipeReference) RevisionOrDefault() string {
	if r.Revision == "" {
		return DefaultRevision
	}
	return r.Revision
}

func (r *RecipeReference) String() string {
	rev := ""
	if r.Revision != "" {
		rev = "#" + r.Revision
	}
	if r.User == "" || r.Channel == "" {
		return fmt.Sprintf("%s/%s%s", r.Name, r.Version, rev)
	}
	return fmt.Sprintf("%s/%s@%s/%s%s", r.Name, r.Version, r.User, r.Channel, rev)
}

func (r *RecipeReference) LinkName() string {
	user := r.User
	if user == "" {
		user = "_"
	}
	channel := r.Channel
	if channel == "" {
		channel = "_"
	}
	return fmt.Sprintf("%s/%s/%s/%s/%s", r.Name, r.Version, user, channel, r.RevisionOrDefault())
}

func (r *RecipeReference) WithRevision(revision string) *RecipeReference {
	return &RecipeReference{r.Name, r.Version, r.User, r.Channel, revision}
}

// AsKey builds the additional key for the package file
func (r *RecipeReference) AsKey() string {
	return fmt.Sprintf("%s|%s|%s", r.User, r.Channel, r.RevisionOrDefault())
}

// PackageReference represents a package of a recipe <Name>/<Version>@<User>/<Channel>#<Revision> <Reference>#<Revision>
type PackageReference struct {
	Recipe    *RecipeReference
	Reference string
	Revision  string
}

func NewPackageReference(recipe *RecipeReference, reference, revision string) (*PackageReference, error) {
	log.Trace("Conan Package: %v %s(#%s)", recipe, reference, revision)

	if recipe == nil {
		return nil, ErrValidation
	}
	if reference == "" || !revisionPattern.MatchString(reference) {
		return nil, ErrValidation
	}
	if revision != "" && !revisionPattern.MatchString(revision) {
		return nil, ErrValidation
	}

	return &PackageReference{recipe, reference, revision}, nil
}

func (r *PackageReference) RevisionOrDefault() string {
	if r.Revision == "" {
		return DefaultRevision
	}
	return r.Revision
}

func (r *PackageReference) LinkName() string {
	return fmt.Sprintf("%s/%s", r.Reference, r.RevisionOrDefault())
}

func (r *PackageReference) WithRevision(revision string) *PackageReference {
	return &PackageReference{r.Recipe, r.Reference, revision}
}

// AsKey builds the additional key for the package file
func (r *PackageReference) AsKey() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", r.Recipe.User, r.Recipe.Channel, r.Recipe.RevisionOrDefault(), r.Reference, r.RevisionOrDefault())
}
