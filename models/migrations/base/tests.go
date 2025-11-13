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

package base

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kumose/kmup/models/unittest"
	"github.com/kumose/kmup/modules/git"
	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/tempdir"
	"github.com/kumose/kmup/modules/test"
	"github.com/kumose/kmup/modules/testlogger"

	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

// FIXME: this file shouldn't be in a normal package, it should only be compiled for tests

// PrepareTestEnv prepares the test environment and reset the database. The skip parameter should usually be 0.
// Provide models to be sync'd with the database - in particular any models you expect fixtures to be loaded from.
//
// fixtures in `models/migrations/fixtures/<TestName>` will be loaded automatically
func PrepareTestEnv(t *testing.T, skip int, syncModels ...any) (*xorm.Engine, func()) {
	t.Helper()
	ourSkip := 2
	ourSkip += skip
	deferFn := testlogger.PrintCurrentTest(t, ourSkip)
	require.NoError(t, unittest.SyncDirs(filepath.Join(filepath.Dir(setting.AppPath), "tests/kmup-repositories-meta"), setting.RepoRootPath))

	if err := deleteDB(); err != nil {
		t.Fatalf("unable to reset database: %v", err)
		return nil, deferFn
	}

	x, err := newXORMEngine()
	require.NoError(t, err)
	if x != nil {
		oldDefer := deferFn
		deferFn = func() {
			oldDefer()
			if err := x.Close(); err != nil {
				t.Errorf("error during close: %v", err)
			}
			if err := deleteDB(); err != nil {
				t.Errorf("unable to reset database: %v", err)
			}
		}
	}
	if err != nil {
		return x, deferFn
	}

	if len(syncModels) > 0 {
		if err := x.Sync(syncModels...); err != nil {
			t.Errorf("error during sync: %v", err)
			return x, deferFn
		}
	}

	fixturesDir := filepath.Join(filepath.Dir(setting.AppPath), "models", "migrations", "fixtures", t.Name())

	if _, err := os.Stat(fixturesDir); err == nil {
		t.Logf("initializing fixtures from: %s", fixturesDir)
		if err := unittest.InitFixtures(
			unittest.FixturesOptions{
				Dir: fixturesDir,
			}, x); err != nil {
			t.Errorf("error whilst initializing fixtures from %s: %v", fixturesDir, err)
			return x, deferFn
		}
		if err := unittest.LoadFixtures(); err != nil {
			t.Errorf("error whilst loading fixtures from %s: %v", fixturesDir, err)
			return x, deferFn
		}
	} else if !os.IsNotExist(err) {
		t.Errorf("unexpected error whilst checking for existence of fixtures: %v", err)
	} else {
		t.Logf("no fixtures found in: %s", fixturesDir)
	}

	return x, deferFn
}

func LoadTableSchemasMap(t *testing.T, x *xorm.Engine) map[string]*schemas.Table {
	tables, err := x.DBMetas()
	require.NoError(t, err)
	tableMap := make(map[string]*schemas.Table)
	for _, table := range tables {
		tableMap[table.Name] = table
	}
	return tableMap
}

func MainTest(m *testing.M) {
	testlogger.Init()

	kmupRoot := test.SetupKmupRoot()
	kmupBinary := "kmup"
	if runtime.GOOS == "windows" {
		kmupBinary += ".exe"
	}
	setting.AppPath = filepath.Join(kmupRoot, kmupBinary)
	if _, err := os.Stat(setting.AppPath); err != nil {
		testlogger.Fatalf("Could not find kmup binary at %s\n", setting.AppPath)
	}

	kmupConf := os.Getenv("KMUP_CONF")
	if kmupConf == "" {
		kmupConf = filepath.Join(filepath.Dir(setting.AppPath), "tests/sqlite.ini")
		_, _ = fmt.Fprintf(os.Stderr, "Environment variable $KMUP_CONF not set - defaulting to %s\n", kmupConf)
	}

	if !filepath.IsAbs(kmupConf) {
		setting.CustomConf = filepath.Join(kmupRoot, kmupConf)
	} else {
		setting.CustomConf = kmupConf
	}

	tmpDataPath, cleanup, err := tempdir.OsTempDir("kmup-test").MkdirTempRandom("data")
	if err != nil {
		testlogger.Fatalf("Unable to create temporary data path %v\n", err)
	}
	defer cleanup()

	setting.CustomPath = filepath.Join(setting.AppWorkPath, "custom")
	setting.AppDataPath = tmpDataPath

	unittest.InitSettingsForTesting()
	if err = git.InitFull(); err != nil {
		testlogger.Fatalf("Unable to InitFull: %v\n", err)
	}
	setting.LoadDBSetting()
	setting.InitLoggersForTest()

	exitStatus := m.Run()

	if err := removeAllWithRetry(setting.RepoRootPath); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "os.RemoveAll: %v\n", err)
	}
	os.Exit(exitStatus)
}
