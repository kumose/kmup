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

package gitrepo

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/kumose/kmup/modules/git/gitcmd"
)

// GetDiffShortStatByCmdArgs counts number of changed files, number of additions and deletions
// TODO: it can be merged with another "GetDiffShortStat" in the future
func GetDiffShortStatByCmdArgs(ctx context.Context, repo Repository, trustedArgs gitcmd.TrustedCmdArgs, dynamicArgs ...string) (numFiles, totalAdditions, totalDeletions int, err error) {
	// Now if we call:
	// $ git diff --shortstat 1ebb35b98889ff77299f24d82da426b434b0cca0...788b8b1440462d477f45b0088875
	// we get:
	// " 9902 files changed, 2034198 insertions(+), 298800 deletions(-)\n"
	cmd := gitcmd.NewCommand("diff", "--shortstat").AddArguments(trustedArgs...).AddDynamicArguments(dynamicArgs...)
	stdout, err := RunCmdString(ctx, repo, cmd)
	if err != nil {
		return 0, 0, 0, err
	}

	return parseDiffStat(stdout)
}

var shortStatFormat = regexp.MustCompile(
	`\s*(\d+) files? changed(?:, (\d+) insertions?\(\+\))?(?:, (\d+) deletions?\(-\))?`)

func parseDiffStat(stdout string) (numFiles, totalAdditions, totalDeletions int, err error) {
	if len(stdout) == 0 || stdout == "\n" {
		return 0, 0, 0, nil
	}
	groups := shortStatFormat.FindStringSubmatch(stdout)
	if len(groups) != 4 {
		return 0, 0, 0, fmt.Errorf("unable to parse shortstat: %s groups: %s", stdout, groups)
	}

	numFiles, err = strconv.Atoi(groups[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unable to parse shortstat: %s. Error parsing NumFiles %w", stdout, err)
	}

	if len(groups[2]) != 0 {
		totalAdditions, err = strconv.Atoi(groups[2])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("unable to parse shortstat: %s. Error parsing NumAdditions %w", stdout, err)
		}
	}

	if len(groups[3]) != 0 {
		totalDeletions, err = strconv.Atoi(groups[3])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("unable to parse shortstat: %s. Error parsing NumDeletions %w", stdout, err)
		}
	}
	return numFiles, totalAdditions, totalDeletions, err
}
