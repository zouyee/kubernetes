/*
Copyright 2019 The Kubernetes Authors.
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

package nodeutilization

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/priorities"
	"k8s.io/kubernetes/pkg/scheduler/framework/plugins/migration"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

// BalancedNodeUtilization is a plugin that calculates the difference between the cpu and memory fraction
// of capacity, and prioritizes the host based on how close the two metrics are to each other.
type BalancedNodeUtilization struct {
	handle framework.FrameworkHandle
}

var _ = framework.ScorePlugin(&BalancedNodeUtilization{})

// Name is the name of the plugin used in the plugin registry and configurations.
const BalancedNodeUtilizationName = "BalancedNodeUtilization"

// Name returns name of the plugin. It is used in logs, etc.
func (br *BalancedNodeUtilization) Name() string {
	return BalancedNodeUtilizationName
}

func (br *BalancedNodeUtilization) Score(state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, exist := br.handle.NodeInfoSnapshot().NodeInfoMap[nodeName]
	if !exist {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("node %q does not exist in NodeInfoSnapshot", nodeName))
	}
	meta := migration.PriorityMetadata(state)

	s, err := priorities.BalancedResourceAllocationMap(pod, meta, nodeInfo)
	return s.Score, migration.ErrorToFrameworkStatus(err)
}

// ScoreExtensions of the Score plugin.
func (pl *BalancedNodeUtilization) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

// NewBalancedNodeUtilization initializes a new plugin and returns it.
func NewBalancedNodeUtilization(_ *runtime.Unknown, h framework.FrameworkHandle) (framework.Plugin, error) {
	return &BalancedNodeUtilization{handle: h}, nil
}
