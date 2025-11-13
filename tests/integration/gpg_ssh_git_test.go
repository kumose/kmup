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

package integration

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"os"
	"testing"

	auth_model "github.com/kumose/kmup/models/auth"
	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/process"
	"github.com/kumose/kmup/modules/setting"
	api "github.com/kumose/kmup/modules/structs"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/tests"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestGPGGit(t *testing.T) {
	tmpDir := t.TempDir() // use a temp dir to avoid messing with the user's GPG keyring
	err := os.Chmod(tmpDir, 0o700)
	assert.NoError(t, err)

	t.Setenv("GNUPGHOME", tmpDir)

	// Need to create a root key
	rootKeyPair, err := importTestingKey()
	require.NoError(t, err, "importTestingKey")

	defer test.MockVariableValue(&setting.Repository.Signing.SigningKey, rootKeyPair.PrimaryKey.KeyIdShortString())()
	defer test.MockVariableValue(&setting.Repository.Signing.SigningName, "kmup")()
	defer test.MockVariableValue(&setting.Repository.Signing.SigningEmail, "kmup@fake.local")()
	defer test.MockVariableValue(&setting.Repository.Signing.InitialCommit, []string{"never"})()
	defer test.MockVariableValue(&setting.Repository.Signing.CRUDActions, []string{"never"})()

	testGitSigning(t)
}

func TestSSHGit(t *testing.T) {
	tmpDir := t.TempDir() // use a temp dir to store the SSH keys
	err := os.Chmod(tmpDir, 0o700)
	assert.NoError(t, err)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err, "ed25519.GenerateKey")
	sshPubKey, err := ssh.NewPublicKey(pub)
	require.NoError(t, err, "ssh.NewPublicKey")

	err = os.WriteFile(tmpDir+"/id_ed25519.pub", ssh.MarshalAuthorizedKey(sshPubKey), 0o600)
	require.NoError(t, err, "os.WriteFile id_ed25519.pub")
	block, err := ssh.MarshalPrivateKey(priv, "")
	require.NoError(t, err, "ssh.MarshalPrivateKey")
	err = os.WriteFile(tmpDir+"/id_ed25519", pem.EncodeToMemory(block), 0o600)
	require.NoError(t, err, "os.WriteFile id_ed25519")

	defer test.MockVariableValue(&setting.Repository.Signing.SigningKey, tmpDir+"/id_ed25519.pub")()
	defer test.MockVariableValue(&setting.Repository.Signing.SigningName, "kmup")()
	defer test.MockVariableValue(&setting.Repository.Signing.SigningEmail, "kmup@fake.local")()
	defer test.MockVariableValue(&setting.Repository.Signing.SigningFormat, "ssh")()
	defer test.MockVariableValue(&setting.Repository.Signing.InitialCommit, []string{"never"})()
	defer test.MockVariableValue(&setting.Repository.Signing.CRUDActions, []string{"never"})()

	testGitSigning(t)
}

func testGitSigning(t *testing.T) {
	username := "user2"
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{Name: username})
	baseAPITestContext := NewAPITestContext(t, username, "repo1")

	onKmupRun(t, func(t *testing.T, u *url.URL) {
		u.Path = baseAPITestContext.GitPath()

		t.Run("Unsigned-Initial", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-unsigned", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateRepository", doAPICreateRepository(testCtx, false))
			t.Run("CheckMasterBranchUnsigned", doAPIGetBranch(testCtx, "master", func(t *testing.T, branch api.Branch) {
				assert.NotNil(t, branch.Commit)
				assert.NotNil(t, branch.Commit.Verification)
				assert.False(t, branch.Commit.Verification.Verified)
				assert.Empty(t, branch.Commit.Verification.Signature)
			}))
			t.Run("CreateCRUDFile-Never", crudActionCreateFile(
				t, testCtx, user, "master", "never", "unsigned-never.txt", func(t *testing.T, response api.FileResponse) {
					assert.False(t, response.Verification.Verified)
				}))
			t.Run("CreateCRUDFile-Never", crudActionCreateFile(
				t, testCtx, user, "never", "never2", "unsigned-never2.txt", func(t *testing.T, response api.FileResponse) {
					assert.False(t, response.Verification.Verified)
				}))
		})

		setting.Repository.Signing.CRUDActions = []string{"parentsigned"}
		t.Run("Unsigned-Initial-CRUD-ParentSigned", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-unsigned", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateCRUDFile-ParentSigned", crudActionCreateFile(
				t, testCtx, user, "master", "parentsigned", "signed-parent.txt", func(t *testing.T, response api.FileResponse) {
					assert.False(t, response.Verification.Verified)
				}))
			t.Run("CreateCRUDFile-ParentSigned", crudActionCreateFile(
				t, testCtx, user, "parentsigned", "parentsigned2", "signed-parent2.txt", func(t *testing.T, response api.FileResponse) {
					assert.False(t, response.Verification.Verified)
				}))
		})

		setting.Repository.Signing.CRUDActions = []string{"never"}
		t.Run("Unsigned-Initial-CRUD-Never", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-unsigned", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateCRUDFile-Never", crudActionCreateFile(
				t, testCtx, user, "parentsigned", "parentsigned-never", "unsigned-never2.txt", func(t *testing.T, response api.FileResponse) {
					assert.False(t, response.Verification.Verified)
				}))
		})

		setting.Repository.Signing.CRUDActions = []string{"always"}
		t.Run("Unsigned-Initial-CRUD-Always", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-unsigned", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateCRUDFile-Always", crudActionCreateFile(
				t, testCtx, user, "master", "always", "signed-always.txt", func(t *testing.T, response api.FileResponse) {
					require.NotNil(t, response.Verification, "no verification provided with response! %v", response)
					require.True(t, response.Verification.Verified)
					assert.Equal(t, "kmup@fake.local", response.Verification.Signer.Email)
				}))
			t.Run("CreateCRUDFile-ParentSigned-always", crudActionCreateFile(
				t, testCtx, user, "parentsigned", "parentsigned-always", "signed-parent2.txt", func(t *testing.T, response api.FileResponse) {
					require.NotNil(t, response.Verification, "no verification provided with response! %v", response)
					require.True(t, response.Verification.Verified)
					assert.Equal(t, "kmup@fake.local", response.Verification.Signer.Email)
				}))
		})

		setting.Repository.Signing.CRUDActions = []string{"parentsigned"}
		t.Run("Unsigned-Initial-CRUD-ParentSigned", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-unsigned", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateCRUDFile-Always-ParentSigned", crudActionCreateFile(
				t, testCtx, user, "always", "always-parentsigned", "signed-always-parentsigned.txt", func(t *testing.T, response api.FileResponse) {
					require.NotNil(t, response.Verification, "no verification provided with response! %v", response)
					require.True(t, response.Verification.Verified)
					assert.Equal(t, "kmup@fake.local", response.Verification.Signer.Email)
				}))
		})

		setting.Repository.Signing.InitialCommit = []string{"always"}
		t.Run("AlwaysSign-Initial", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-always", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateRepository", doAPICreateRepository(testCtx, false))
			t.Run("CheckMasterBranchSigned", doAPIGetBranch(testCtx, "master", func(t *testing.T, branch api.Branch) {
				require.NotNil(t, branch.Commit, "no commit provided with branch! %v", branch)
				require.NotNil(t, branch.Commit.Verification, "no verification provided with branch commit! %v", branch.Commit)
				require.True(t, branch.Commit.Verification.Verified)
				assert.Equal(t, "kmup@fake.local", branch.Commit.Verification.Signer.Email)
			}))
		})

		setting.Repository.Signing.CRUDActions = []string{"never"}
		t.Run("AlwaysSign-Initial-CRUD-Never", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-always-never", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateRepository", doAPICreateRepository(testCtx, false))
			t.Run("CreateCRUDFile-Never", crudActionCreateFile(
				t, testCtx, user, "master", "never", "unsigned-never.txt", func(t *testing.T, response api.FileResponse) {
					assert.False(t, response.Verification.Verified)
				}))
		})

		setting.Repository.Signing.CRUDActions = []string{"parentsigned"}
		t.Run("AlwaysSign-Initial-CRUD-ParentSigned-On-Always", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-always-parent", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateRepository", doAPICreateRepository(testCtx, false))
			t.Run("CreateCRUDFile-ParentSigned", crudActionCreateFile(
				t, testCtx, user, "master", "parentsigned", "signed-parent.txt", func(t *testing.T, response api.FileResponse) {
					require.True(t, response.Verification.Verified)
					assert.Equal(t, "kmup@fake.local", response.Verification.Signer.Email)
				}))
		})

		setting.Repository.Signing.CRUDActions = []string{"always"}
		t.Run("AlwaysSign-Initial-CRUD-Always", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-always-always", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreateRepository", doAPICreateRepository(testCtx, false))
			t.Run("CreateCRUDFile-Always", crudActionCreateFile(
				t, testCtx, user, "master", "always", "signed-always.txt", func(t *testing.T, response api.FileResponse) {
					require.True(t, response.Verification.Verified)
					assert.Equal(t, "kmup@fake.local", response.Verification.Signer.Email)
				}))
		})

		setting.Repository.Signing.Merges = []string{"commitssigned"}
		t.Run("UnsignedMerging", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-unsigned", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreatePullRequest", func(t *testing.T) {
				pr, err := doAPICreatePullRequest(testCtx, testCtx.Username, testCtx.Reponame, "master", "never2")(t)
				assert.NoError(t, err)
				t.Run("MergePR", doAPIMergePullRequest(testCtx, testCtx.Username, testCtx.Reponame, pr.Index))
			})
			t.Run("CheckMasterBranchUnsigned", doAPIGetBranch(testCtx, "master", func(t *testing.T, branch api.Branch) {
				assert.NotNil(t, branch.Commit)
				assert.NotNil(t, branch.Commit.Verification)
				assert.False(t, branch.Commit.Verification.Verified)
				assert.Empty(t, branch.Commit.Verification.Signature)
			}))
		})

		setting.Repository.Signing.Merges = []string{"basesigned"}
		t.Run("BaseSignedMerging", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-unsigned", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreatePullRequest", func(t *testing.T) {
				pr, err := doAPICreatePullRequest(testCtx, testCtx.Username, testCtx.Reponame, "master", "parentsigned2")(t)
				assert.NoError(t, err)
				t.Run("MergePR", doAPIMergePullRequest(testCtx, testCtx.Username, testCtx.Reponame, pr.Index))
			})
			t.Run("CheckMasterBranchUnsigned", doAPIGetBranch(testCtx, "master", func(t *testing.T, branch api.Branch) {
				assert.NotNil(t, branch.Commit)
				assert.NotNil(t, branch.Commit.Verification)
				assert.False(t, branch.Commit.Verification.Verified)
				assert.Empty(t, branch.Commit.Verification.Signature)
			}))
		})

		setting.Repository.Signing.Merges = []string{"commitssigned"}
		t.Run("CommitsSignedMerging", func(t *testing.T) {
			defer tests.PrintCurrentTest(t)()
			testCtx := NewAPITestContext(t, username, "initial-unsigned", auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopeWriteUser)
			t.Run("CreatePullRequest", func(t *testing.T) {
				pr, err := doAPICreatePullRequest(testCtx, testCtx.Username, testCtx.Reponame, "master", "always-parentsigned")(t)
				assert.NoError(t, err)
				t.Run("MergePR", doAPIMergePullRequest(testCtx, testCtx.Username, testCtx.Reponame, pr.Index))
			})
			t.Run("CheckMasterBranchUnsigned", doAPIGetBranch(testCtx, "master", func(t *testing.T, branch api.Branch) {
				assert.NotNil(t, branch.Commit)
				assert.NotNil(t, branch.Commit.Verification)
				assert.True(t, branch.Commit.Verification.Verified)
			}))
		})
	})
}

func crudActionCreateFile(_ *testing.T, ctx APITestContext, user *user_model.User, from, to, path string, callback ...func(*testing.T, api.FileResponse)) func(*testing.T) {
	return doAPICreateFile(ctx, path, &api.CreateFileOptions{
		FileOptions: api.FileOptions{
			BranchName:    from,
			NewBranchName: to,
			Message:       fmt.Sprintf("from:%s to:%s path:%s", from, to, path),
			Author: api.Identity{
				Name:  user.FullName,
				Email: user.Email,
			},
			Committer: api.Identity{
				Name:  user.FullName,
				Email: user.Email,
			},
		},
		ContentBase64: base64.StdEncoding.EncodeToString([]byte("This is new text for " + path)),
	}, callback...)
}

func importTestingKey() (*openpgp.Entity, error) {
	if _, _, err := process.GetManager().Exec("gpg --import tests/integration/private-testing.key", "gpg", "--import", "tests/integration/private-testing.key"); err != nil {
		return nil, err
	}
	keyringFile, err := os.Open("tests/integration/private-testing.key")
	if err != nil {
		return nil, err
	}
	defer keyringFile.Close()

	block, err := armor.Decode(keyringFile)
	if err != nil {
		return nil, err
	}

	keyring, err := openpgp.ReadKeyRing(block.Body)
	if err != nil {
		return nil, fmt.Errorf("Keyring access failed: '%w'", err)
	}

	// There should only be one entity in this file.
	return keyring[0], nil
}
