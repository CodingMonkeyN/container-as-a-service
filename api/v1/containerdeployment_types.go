/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ContainerDeploymentSpec struct {
	Image           string            `json:"image"`
	Namespace       string            `json:"namespace,omitempty"`
	CPU             string            `json:"cpu,omitempty"`
	Memory          string            `json:"memory,omitempty"`
	Port            int32             `json:"port"`
	Storage         *StorageSpec      `json:"storage,omitempty"`
	EnvironmentVars map[string]string `json:"env,omitempty"`
}

type StorageSpec struct {
	Size      string `json:"size"`
	MountPath string `json:"mount-path"`
}

type ContainerDeploymentStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Schema for ContainerDeployment
type ContainerDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContainerDeploymentSpec   `json:"spec,omitempty"`
	Status ContainerDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ContainerDeploymentList contains a list of ContainerDeployment
type ContainerDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContainerDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContainerDeployment{}, &ContainerDeploymentList{})
}
