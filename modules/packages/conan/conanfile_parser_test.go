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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	name             = "ConanPackage"
	version          = "1.2"
	license          = "MIT"
	author           = "Kmup <info@kmup.io>"
	homepage         = "https://kmup.io/"
	url              = "https://kmup.com/"
	description      = "Description of ConanPackage"
	topic1           = "kmup"
	topic2           = "conan"
	contentConanfile = `from conans import ConanFile, CMake, tools

class ConanPackageConan(ConanFile):
    name = "` + name + `"
    version = "` + version + `"
    license = "` + license + `"
    author = "` + author + `"
    homepage = "` + homepage + `"
    url = "` + url + `"
    description = "` + description + `"
    topics = ("` + topic1 + `", "` + topic2 + `")
    settings = "os", "compiler", "build_type", "arch"
    options = {"shared": [True, False], "fPIC": [True, False]}
    default_options = {"shared": False, "fPIC": True}
    generators = "cmake"
`
)

func TestParseConanfile(t *testing.T) {
	metadata, err := ParseConanfile(strings.NewReader(contentConanfile))
	assert.NoError(t, err)
	assert.Equal(t, license, metadata.License)
	assert.Equal(t, author, metadata.Author)
	assert.Equal(t, homepage, metadata.ProjectURL)
	assert.Equal(t, url, metadata.RepositoryURL)
	assert.Equal(t, description, metadata.Description)
	assert.Equal(t, []string{topic1, topic2}, metadata.Keywords)
}
