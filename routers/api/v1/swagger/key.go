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

package swagger

import (
	api "github.com/kumose/kmup/modules/structs"
)

// PublicKey
// swagger:response PublicKey
type swaggerResponsePublicKey struct {
	// in:body
	Body api.PublicKey `json:"body"`
}

// PublicKeyList
// swagger:response PublicKeyList
type swaggerResponsePublicKeyList struct {
	// in:body
	Body []api.PublicKey `json:"body"`
}

// GPGKey
// swagger:response GPGKey
type swaggerResponseGPGKey struct {
	// in:body
	Body api.GPGKey `json:"body"`
}

// GPGKeyList
// swagger:response GPGKeyList
type swaggerResponseGPGKeyList struct {
	// in:body
	Body []api.GPGKey `json:"body"`
}

// DeployKey
// swagger:response DeployKey
type swaggerResponseDeployKey struct {
	// in:body
	Body api.DeployKey `json:"body"`
}

// DeployKeyList
// swagger:response DeployKeyList
type swaggerResponseDeployKeyList struct {
	// in:body
	Body []api.DeployKey `json:"body"`
}
