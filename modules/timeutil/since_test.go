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

package timeutil

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kumose/kmup/modules/setting"
	"github.com/kumose/kmup/modules/translation"

	"github.com/stretchr/testify/assert"
)

var BaseDate time.Time

// time durations
const (
	DayDur   = 24 * time.Hour
	WeekDur  = 7 * DayDur
	MonthDur = 30 * DayDur
	YearDur  = 12 * MonthDur
)

func TestMain(m *testing.M) {
	setting.StaticRootPath = "../../"
	setting.Names = []string{"english"}
	setting.Langs = []string{"en-US"}
	// setup
	translation.InitLocales(context.Background())
	BaseDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	// run the tests
	retVal := m.Run()

	os.Exit(retVal)
}

func TestTimeSincePro(t *testing.T) {
	assert.Equal(t, "now", timeSincePro(BaseDate, BaseDate, translation.NewLocale("en-US")))

	// test that a difference of `diff` yields the expected string
	test := func(expected string, diff time.Duration) {
		actual := timeSincePro(BaseDate, BaseDate.Add(diff), translation.NewLocale("en-US"))
		assert.Equal(t, expected, actual)
		assert.Equal(t, "future", timeSincePro(BaseDate.Add(diff), BaseDate, translation.NewLocale("en-US")))
	}
	test("1 second", time.Second)
	test("2 seconds", 2*time.Second)
	test("1 minute", time.Minute)
	test("1 minute, 1 second", time.Minute+time.Second)
	test("1 minute, 59 seconds", time.Minute+59*time.Second)
	test("2 minutes", 2*time.Minute)
	test("1 hour", time.Hour)
	test("1 hour, 1 second", time.Hour+time.Second)
	test("1 hour, 59 minutes, 59 seconds", time.Hour+59*time.Minute+59*time.Second)
	test("2 hours", 2*time.Hour)
	test("1 day", DayDur)
	test("1 day, 23 hours, 59 minutes, 59 seconds",
		DayDur+23*time.Hour+59*time.Minute+59*time.Second)
	test("2 days", 2*DayDur)
	test("1 week", WeekDur)
	test("2 weeks", 2*WeekDur)
	test("1 month", MonthDur)
	test("3 months", 3*MonthDur)
	test("1 year", YearDur)
	test("2 years, 3 months, 1 week, 2 days, 4 hours, 12 minutes, 17 seconds",
		2*YearDur+3*MonthDur+WeekDur+2*DayDur+4*time.Hour+
			12*time.Minute+17*time.Second)
}

func TestMinutesToFriendly(t *testing.T) {
	// test that a number of minutes yields the expected string
	test := func(expected string, minutes int) {
		actual := MinutesToFriendly(minutes, translation.NewLocale("en-US"))
		assert.Equal(t, expected, actual)
	}
	test("1 minute", 1)
	test("2 minutes", 2)
	test("1 hour", 60)
	test("1 hour, 1 minute", 61)
	test("1 hour, 2 minutes", 62)
	test("2 hours", 120)
}
