/*


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

package v1alpha1

import (
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HelmReleaseSpec defines the desired state of HelmRelease
type HelmReleaseSpec struct {
	// +optional
	Atomic bool `json:"atomic"`

	Chart string `json:"chart,required"`

	// +optional
	CleanupOnFail bool `json:"cleanupOnFail"`

	// +optional
	DisableHooks bool `json:"disableHooks"`

	// +optional
	DisableOpenAPIValidation bool `json:"disableOpenAPIValidation"`

	// +optional
	ForceUpgrade bool `json:"forceUpgrade"`

	// +optional
	IncludeCRDs bool `json:"includeCRDS"`

	Version string `json:"version,required"`

	// +optional
	Wait bool `json:"wait"`

	// +optional
	WaitTimeout int `json:"waitTimeout"`

	// +optional
	DryRun bool `json:"dryRun"`

	// +optional
	ValuesOverride Values `json:"values"`
}

type Values struct {
	V map[string]interface{} `json:"-"`
}

// DeepCopyInto is an deepcopy function, copying the receiver, writing
// into out.
func (v *Values) DeepCopyInto(out *Values) {
	b, err := json.Marshal(v.V)
	if err != nil {
		panic(err)
	}
	var c map[string]interface{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}
	out.V = c
	return
}

// MarshalJSON marshals the HelmValues data to a JSON blob.
func (v Values) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.V)
}

// UnmarshalJSON sets the HelmValues to a copy of data.
func (v *Values) UnmarshalJSON(data []byte) error {
	var out map[string]interface{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}
	v.V = out
	return nil
}

// HelmReleaseStatus defines the observed state of HelmRelease
type HelmReleaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// HelmRelease is the Schema for the helmreleases API
// +kubebuilder:printcolumn:name="release-name",type=string,JSONPath=.metadata.name
// +kubebuilder:printcolumn:name="release-namespace",type=string,JSONPath=.metadata.namespace
// +kubebuilder:printcolumn:name="chart",type=string,JSONPath=.spec.chart
// +kubebuilder:printcolumn:name="chart-version",type=string,JSONPath=.spec.version
// +kubebuilder:printcolumn:name="age",type=date,JSONPath=.metadata.creationTimestamp
type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmReleaseSpec   `json:"spec,omitempty"`
	Status HelmReleaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HelmReleaseList contains a list of HelmRelease
type HelmReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmRelease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelmRelease{}, &HelmReleaseList{})
}
