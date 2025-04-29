package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Severity",type=string,JSONPath=`.spec.severity`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.issueType`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.state`

// Issue is the Schema for the issues API
type Issue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IssueSpec   `json:"spec,omitempty"`
	Status IssueStatus `json:"status,omitempty"`
}

// IssueSpec defines the desired state of Issue
type IssueSpec struct {
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Severity      string      `json:"severity"`
	IssueType     string      `json:"issueType"`
	Scope         IssueScope  `json:"scope,omitempty"`
	RelatedIssues []string    `json:"relatedIssues,omitempty"`
	Links         []IssueLink `json:"links,omitempty"`
}

// IssueScope defines what resource this issue is related to
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

// IssueStatus defines the observed state of Issue
type IssueStatus struct {
	State             string `json:"state,omitempty"`
	ResolvedTimestamp string `json:"resolvedTimestamp,omitempty"`
	DetectedTimestamp string `json:"detectedTimestamp,omitempty"`
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
