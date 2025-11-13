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

package asymkey

import (
	"context"
	"strings"

	"github.com/kumose/kmup/models/db"
	"github.com/kumose/kmup/modules/log"

	"github.com/ProtonMail/go-crypto/openpgp"
)

//   __________________  ________   ____  __.
//  /  _____/\______   \/  _____/  |    |/ _|____ ___.__.
// /   \  ___ |     ___/   \  ___  |      <_/ __ <   |  |
// \    \_\  \|    |   \    \_\  \ |    |  \  ___/\___  |
//  \______  /|____|    \______  / |____|__ \___  > ____|
//         \/                  \/          \/   \/\/
//    _____       .___  .___
//   /  _  \    __| _/__| _/
//	/  /_\  \  / __ |/ __ |
// /    |    \/ /_/ / /_/ |
// \____|__  /\____ \____ |
//         \/      \/    \/

// This file contains functions relating to adding GPG Keys

// addGPGKey add key, import and subkeys to database
func addGPGKey(ctx context.Context, key *GPGKey, content string) (err error) {
	// Add GPGKeyImport
	if err = db.Insert(ctx, &GPGKeyImport{
		KeyID:   key.KeyID,
		Content: content,
	}); err != nil {
		return err
	}
	// Save GPG primary key.
	if err = db.Insert(ctx, key); err != nil {
		return err
	}
	// Save GPG subs key.
	for _, subkey := range key.SubsKey {
		if err := addGPGSubKey(ctx, subkey); err != nil {
			return err
		}
	}
	return nil
}

// addGPGSubKey add subkeys to database
func addGPGSubKey(ctx context.Context, key *GPGKey) (err error) {
	// Save GPG primary key.
	if err = db.Insert(ctx, key); err != nil {
		return err
	}
	// Save GPG subs key.
	for _, subkey := range key.SubsKey {
		if err := addGPGSubKey(ctx, subkey); err != nil {
			return err
		}
	}
	return nil
}

// AddGPGKey adds new public key to database.
func AddGPGKey(ctx context.Context, ownerID int64, content, token, signature string) ([]*GPGKey, error) {
	ekeys, err := CheckArmoredGPGKeyString(content)
	if err != nil {
		return nil, err
	}

	return db.WithTx2(ctx, func(ctx context.Context) ([]*GPGKey, error) {
		keys := make([]*GPGKey, 0, len(ekeys))

		verified := false
		// Handle provided signature
		if signature != "" {
			signer, err := openpgp.CheckArmoredDetachedSignature(ekeys, strings.NewReader(token), strings.NewReader(signature), nil)
			if err != nil {
				signer, err = openpgp.CheckArmoredDetachedSignature(ekeys, strings.NewReader(token+"\n"), strings.NewReader(signature), nil)
			}
			if err != nil {
				signer, err = openpgp.CheckArmoredDetachedSignature(ekeys, strings.NewReader(token+"\r\n"), strings.NewReader(signature), nil)
			}
			if err != nil {
				log.Debug("AddGPGKey CheckArmoredDetachedSignature failed: %v", err)
				return nil, ErrGPGInvalidTokenSignature{
					ID:      ekeys[0].PrimaryKey.KeyIdString(),
					Wrapped: err,
				}
			}
			ekeys = []*openpgp.Entity{signer}
			verified = true
		}

		if len(ekeys) > 1 {
			id2key := map[string]*openpgp.Entity{}
			newEKeys := make([]*openpgp.Entity, 0, len(ekeys))
			for _, ekey := range ekeys {
				id := ekey.PrimaryKey.KeyIdString()
				if original, has := id2key[id]; has {
					// Coalesce this with the other one
					for _, subkey := range ekey.Subkeys {
						if subkey.PublicKey == nil {
							continue
						}
						found := false

						for _, originalSubkey := range original.Subkeys {
							if originalSubkey.PublicKey == nil {
								continue
							}
							if originalSubkey.PublicKey.KeyId == subkey.PublicKey.KeyId {
								found = true
								break
							}
						}
						if !found {
							original.Subkeys = append(original.Subkeys, subkey)
						}
					}
					for name, identity := range ekey.Identities {
						if _, has := original.Identities[name]; has {
							continue
						}
						original.Identities[name] = identity
					}
					continue
				}
				id2key[id] = ekey
				newEKeys = append(newEKeys, ekey)
			}
			ekeys = newEKeys
		}

		for _, ekey := range ekeys {
			// Key ID cannot be duplicated.
			has, err := db.GetEngine(ctx).Where("key_id=?", ekey.PrimaryKey.KeyIdString()).
				Get(new(GPGKey))
			if err != nil {
				return nil, err
			} else if has {
				return nil, ErrGPGKeyIDAlreadyUsed{ekey.PrimaryKey.KeyIdString()}
			}

			key, err := parseGPGKey(ctx, ownerID, ekey, verified)
			if err != nil {
				return nil, err
			}

			if err = addGPGKey(ctx, key, content); err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		return keys, nil
	})
}
