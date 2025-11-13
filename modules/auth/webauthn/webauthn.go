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

package webauthn

import (
	"context"
	"encoding/binary"
	"encoding/gob"

	"github.com/kumose/kmup/models/auth"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/util"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

// WebAuthn represents the global WebAuthn instance
var WebAuthn *webauthn.WebAuthn

// Init initializes the WebAuthn instance from the config.
func Init() {
	gob.Register(&webauthn.SessionData{})

	appURL, _ := protocol.FullyQualifiedOrigin(setting.AppURL)

	WebAuthn = &webauthn.WebAuthn{
		Config: &webauthn.Config{
			RPDisplayName: setting.AppName,
			RPID:          setting.Domain,
			RPOrigins:     []string{appURL},
			AuthenticatorSelection: protocol.AuthenticatorSelection{
				UserVerification: protocol.VerificationDiscouraged,
			},
			AttestationPreference: protocol.PreferDirectAttestation,
		},
	}
}

// user represents an implementation of webauthn.User based on User model
type user struct {
	ctx  context.Context
	User *user_model.User

	defaultAuthFlags protocol.AuthenticatorFlags
}

var _ webauthn.User = (*user)(nil)

func NewWebAuthnUser(ctx context.Context, u *user_model.User, defaultAuthFlags ...protocol.AuthenticatorFlags) webauthn.User {
	return &user{ctx: ctx, User: u, defaultAuthFlags: util.OptionalArg(defaultAuthFlags)}
}

// WebAuthnID implements the webauthn.User interface
func (u *user) WebAuthnID() []byte {
	id := make([]byte, 8)
	binary.PutVarint(id, u.User.ID)
	return id
}

// WebAuthnName implements the webauthn.User interface
func (u *user) WebAuthnName() string {
	return util.IfZero(u.User.LoginName, u.User.Name)
}

// WebAuthnDisplayName implements the webauthn.User interface
func (u *user) WebAuthnDisplayName() string {
	return u.User.DisplayName()
}

// WebAuthnCredentials implements the webauthn.User interface
func (u *user) WebAuthnCredentials() []webauthn.Credential {
	dbCreds, err := auth.GetWebAuthnCredentialsByUID(u.ctx, u.User.ID)
	if err != nil {
		return nil
	}
	return dbCreds.ToCredentials(u.defaultAuthFlags)
}
