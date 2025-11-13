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

package v1_17

import (
	packages_model "github.com/kumose/kmup/models/packages"
	container_module "github.com/kumose/kmup/modules/packages/container"

	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

func AddContainerRepositoryProperty(x *xorm.Engine) (err error) {
	switch x.Dialect().URI().DBType {
	case schemas.SQLITE:
		_, err = x.Exec("INSERT INTO package_property (ref_type, ref_id, name, value) SELECT ?, p.id, ?, u.lower_name || '/' || p.lower_name FROM package p JOIN `user` u ON p.owner_id = u.id WHERE p.type = ?",
			packages_model.PropertyTypePackage, container_module.PropertyRepository, packages_model.TypeContainer)
	case schemas.MSSQL:
		_, err = x.Exec("INSERT INTO package_property (ref_type, ref_id, name, value) SELECT ?, p.id, ?, u.lower_name + '/' + p.lower_name FROM package p JOIN `user` u ON p.owner_id = u.id WHERE p.type = ?",
			packages_model.PropertyTypePackage, container_module.PropertyRepository, packages_model.TypeContainer)
	default:
		_, err = x.Exec("INSERT INTO package_property (ref_type, ref_id, name, value) SELECT ?, p.id, ?, CONCAT(u.lower_name, '/', p.lower_name) FROM package p JOIN `user` u ON p.owner_id = u.id WHERE p.type = ?",
			packages_model.PropertyTypePackage, container_module.PropertyRepository, packages_model.TypeContainer)
	}
	return err
}
