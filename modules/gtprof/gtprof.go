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

package gtprof

// This is a Kmup-specific profiling package,
// the name is chosen to distinguish it from the standard pprof tool and "GNU gprof"

// LabelGracefulLifecycle is a label marking manager lifecycle phase
// Making it compliant with prometheus key regex https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels
// would enable someone interested to be able to continuously gather profiles into pyroscope.
// Other labels for pprof should also follow this rule.
const LabelGracefulLifecycle = "graceful_lifecycle"

// LabelPid is a label set on goroutines that have a process attached
const LabelPid = "pid"

// LabelPpid is a label set on goroutines that have a process attached
const LabelPpid = "ppid"

// LabelProcessType is a label set on goroutines that have a process attached
const LabelProcessType = "process_type"

// LabelProcessDescription is a label set on goroutines that have a process attached
const LabelProcessDescription = "process_description"
