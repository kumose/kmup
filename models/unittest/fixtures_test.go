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

package unittest_test

import (
	"path/filepath"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	user_model "github.com/kumose/kmup/models/user"
	"github.com/kumose/kmup/modules/test"

	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
)

var NewFixturesLoaderVendor = func(e *xorm.Engine, opts unittest.FixturesOptions) (unittest.FixturesLoader, error) {
	return nil, nil
}

/*
// the old code is kept here in case we are still interested in benchmarking the two implementations
func init() {
	NewFixturesLoaderVendor = func(e *xorm.Engine, opts unittest.FixturesOptions) (unittest.FixturesLoader, error) {
		return NewFixturesLoaderVendorGoTestfixtures(e, opts)
	}
}

func NewFixturesLoaderVendorGoTestfixtures(e *xorm.Engine, opts unittest.FixturesOptions) (*testfixtures.Loader, error) {
	files, err := unittest.FixturesFileFullPaths(opts.Dir, opts.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to get fixtures files: %w", err)
	}
	var dialect string
	switch e.Dialect().URI().DBType {
	case schemas.POSTGRES:
		dialect = "postgres"
	case schemas.MYSQL:
		dialect = "mysql"
	case schemas.MSSQL:
		dialect = "mssql"
	case schemas.SQLITE:
		dialect = "sqlite3"
	default:
		return nil, fmt.Errorf("unsupported RDBMS for integration tests: %q", e.Dialect().URI().DBType)
	}
	loaderOptions := []func(loader *testfixtures.Loader) error{
		testfixtures.Database(e.DB().DB),
		testfixtures.Dialect(dialect),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Files(files...),
	}
	if e.Dialect().URI().DBType == schemas.POSTGRES {
		loaderOptions = append(loaderOptions, testfixtures.SkipResetSequences())
	}
	return testfixtures.New(loaderOptions...)
}
*/

func prepareTestFixturesLoaders(t testing.TB) unittest.FixturesOptions {
	_ = user_model.User{}
	opts := unittest.FixturesOptions{Dir: filepath.Join(test.SetupKmupRoot(), "models", "fixtures"), Files: []string{
		"user.yml",
	}}
	require.NoError(t, unittest.CreateTestEngine(opts))
	return opts
}

func TestFixturesLoader(t *testing.T) {
	opts := prepareTestFixturesLoaders(t)
	loaderInternal, err := unittest.NewFixturesLoader(unittest.GetXORMEngine(), opts)
	require.NoError(t, err)
	loaderVendor, err := NewFixturesLoaderVendor(unittest.GetXORMEngine(), opts)
	require.NoError(t, err)
	t.Run("Internal", func(t *testing.T) {
		require.NoError(t, loaderInternal.Load())
		require.NoError(t, loaderInternal.Load())
	})
	t.Run("Vendor", func(t *testing.T) {
		if loaderVendor == nil {
			t.Skip()
		}
		require.NoError(t, loaderVendor.Load())
		require.NoError(t, loaderVendor.Load())
	})
}

func BenchmarkFixturesLoader(b *testing.B) {
	opts := prepareTestFixturesLoaders(b)
	require.NoError(b, unittest.CreateTestEngine(opts))
	loaderInternal, err := unittest.NewFixturesLoader(unittest.GetXORMEngine(), opts)
	require.NoError(b, err)
	loaderVendor, err := NewFixturesLoaderVendor(unittest.GetXORMEngine(), opts)
	require.NoError(b, err)

	// BenchmarkFixturesLoader/Vendor
	// BenchmarkFixturesLoader/Vendor-12         	    1696	    719416 ns/op
	// BenchmarkFixturesLoader/Internal
	// BenchmarkFixturesLoader/Internal-12       	    1746	    670457 ns/op
	b.Run("Internal", func(b *testing.B) {
		for b.Loop() {
			require.NoError(b, loaderInternal.Load())
		}
	})
	b.Run("Vendor", func(b *testing.B) {
		if loaderVendor == nil {
			b.Skip()
		}
		for b.Loop() {
			require.NoError(b, loaderVendor.Load())
		}
	})
}
