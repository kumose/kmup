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

package actions

import (
	actions_model "github.com/kumose/kmup/models/actions"
	"github.com/kumose/kmup/modules/container"
)

// GetAllRerunJobs get all jobs that need to be rerun when job should be rerun
func GetAllRerunJobs(job *actions_model.ActionRunJob, allJobs []*actions_model.ActionRunJob) []*actions_model.ActionRunJob {
	rerunJobs := []*actions_model.ActionRunJob{job}
	rerunJobsIDSet := make(container.Set[string])
	rerunJobsIDSet.Add(job.JobID)

	for {
		found := false
		for _, j := range allJobs {
			if rerunJobsIDSet.Contains(j.JobID) {
				continue
			}
			for _, need := range j.Needs {
				if rerunJobsIDSet.Contains(need) {
					found = true
					rerunJobs = append(rerunJobs, j)
					rerunJobsIDSet.Add(j.JobID)
					break
				}
			}
		}
		if !found {
			break
		}
	}

	return rerunJobs
}
