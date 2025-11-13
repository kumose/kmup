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

package rubygems

import (
	"strings"
	"testing"

	rubygems_module "github.com/kumose/kmup/modules/packages/rubygems"

	"github.com/stretchr/testify/assert"
)

func TestWritePackageVersion(t *testing.T) {
	buf := &strings.Builder{}

	writePackageVersionForList(nil, "1.0", " ", buf)
	assert.Equal(t, "1.0 ", buf.String())
	buf.Reset()

	writePackageVersionForList(&rubygems_module.Metadata{Platform: "ruby"}, "1.0", " ", buf)
	assert.Equal(t, "1.0 ", buf.String())
	buf.Reset()

	writePackageVersionForList(&rubygems_module.Metadata{Platform: "linux"}, "1.0", " ", buf)
	assert.Equal(t, "1.0_linux ", buf.String())
	buf.Reset()

	writePackageVersionForDependency("1.0", "", buf)
	assert.Equal(t, "1.0 ", buf.String())
	buf.Reset()

	writePackageVersionForDependency("1.0", "ruby", buf)
	assert.Equal(t, "1.0 ", buf.String())
	buf.Reset()

	writePackageVersionForDependency("1.0", "os", buf)
	assert.Equal(t, "1.0-os ", buf.String())
	buf.Reset()
}
