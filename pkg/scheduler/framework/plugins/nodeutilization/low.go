/*
Copyright 201 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package leastresourceallocation

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/priorities"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/migration"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

// LowNodeUtilization is a priority that favors nodes with fewer requested resources.
type LowNodeUtilization struct {
	handle framework.FrameworkHandle
}

var _ = framework.ScorePlugin(&LowNodeUtilization{})

// LowNodeUtilizationName is the name of the plugin used in the plugin registry and configurations.
const LowNodeUtilizationName = "LowNodeUtilization"

// Name returns name of the plugin. It is used in logs, etc.
func (br *LowNodeUtilization) Name() string {
	return LowNodeUtilizationName
}

func (br *LowNodeUtilization) Score(state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, exist := br.handle.NodeInfoSnapshot().NodeInfoMap[nodeName]
	if !exist {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("node %q does not exist in NodeInfoSnapshot", nodeName))
	}
	meta := migration.PriorityMetadata(state)

	s, err := priorities.LeastRequestedPriorityMap(pod, meta, nodeInfo)
	return s.Score, migration.ErrorToFrameworkStatus(err)
}

// ScoreExtensions of the Score plugin.
func (pl *LowNodeUtilization) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

// NewLowNodeUtilization initializes a new plugin and returns it.
func NewLowNodeUtilization(_ *runtime.Unknown, h framework.FrameworkHandle) (framework.Plugin, error) {
	return &LowNodeUtilization{handle: h}, nil
}
