/*
Copyright 2025.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IssueSpec defines the desired state of Issue
type IssueSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Severity      string      `json:"severity"`
	IssueType     string      `json:"issueType"`
	Scope         IssueScope  `json:"scope,omitempty"`
	RelatedIssues []string    `json:"relatedIssues,omitempty"`
	Links         []IssueLink `json:"links,omitempty"`
}

type IssueScope struct {
	ResourceType      string `json:"resourceType"`
	ResourceName      string `json:"resourceName"`
	ResourceNamespace string `json:"resourceNamespace,omitempty"`
}

// IssueLink provides links to related resources
type IssueLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// IssueCondition defines the observed state of Issue
type IssueCondition struct {
	// INSERT ADDITIONAL CONDITION FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State             string `json:"state,omitempty"`
	ResolvedTimestamp string `json:"resolvedTimestamp,omitempty"`
	DetectedTimestamp string `json:"detectedTimestamp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=isu,shortName=iss

// Issue is the Schema for the issues API
type Issue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec      IssueSpec      `json:"spec,omitempty"`
	Condition IssueCondition `json:"condition,omitempty"`
}

// +kubebuilder:object:root=true

// IssueList contains a list of Issue
type IssueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Issue `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Issue{}, &IssueList{})
}
