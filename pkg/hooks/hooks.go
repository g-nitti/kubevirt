/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright The KubeVirt Authors.
 *
 */

package hooks

import (
	"encoding/json"

	k8sv1 "k8s.io/api/core/v1"

	v1 "kubevirt.io/api/core/v1"
)

const HookSidecarListAnnotationName = "hooks.kubevirt.io/hookSidecars"
const HookSocketsSharedDirectory = "/var/run/kubevirt-hooks"

const ContainerNameEnvVar = "CONTAINER_NAME"

type HookSidecarList []HookSidecar

type ConfigMap struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	HookPath string `json:"hookPath"`
}

type PVC struct {
	Name              string `json:"name"`
	VolumePath        string `json:"volumePath"`
	SharedComputePath string `json:"sharedComputePath"`
}

type HookSidecar struct {
	Image           string                           `json:"image,omitempty"`
	ImagePullPolicy k8sv1.PullPolicy                 `json:"imagePullPolicy"`
	Command         []string                         `json:"command,omitempty"`
	Args            []string                         `json:"args,omitempty"`
	ConfigMap       *ConfigMap                       `json:"configMap,omitempty"`
	PVC             *PVC                             `json:"pvc,omitempty"`
	DownwardAPI     v1.NetworkBindingDownwardAPIType `json:"-"`
	Name			string                           `json:"name"`
}

func UnmarshalHookSidecarList(vmiObject *v1.VirtualMachineInstance) (HookSidecarList, error) {
	hookSidecarList := make(HookSidecarList, 0)

	if rawRequestedHookSidecarList, requestedHookSidecarListDefined := vmiObject.GetAnnotations()[HookSidecarListAnnotationName]; requestedHookSidecarListDefined {
		if err := json.Unmarshal([]byte(rawRequestedHookSidecarList), &hookSidecarList); err != nil {
			return nil, err
		}
	}

	return hookSidecarList, nil
}
