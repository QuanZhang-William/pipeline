/*
Copyright 2021 The Tekton Authors

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

package experimentpod

import (
	"context"
	"fmt"

	"github.com/tektoncd/pipeline/pkg/pod"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewTransformer returns a pod.Transformer that will pod affinity if needed
func NewTransformer(_ context.Context, annotations map[string]string) pod.Transformer {
	return func(p *corev1.Pod) (*corev1.Pod, error) {
		// if it is an anchor pod, don't append pod affinity
		if isFirstPod := annotations["first-pod"]; isFirstPod == "true" {
			return p, nil
		}

		if p.Spec.Affinity == nil {
			p.Spec.Affinity = &corev1.Affinity{}
		}

		anchorPod, ok := annotations["anchor-pod"]
		if !ok {
			return p, fmt.Errorf("missing anchor pod")
		}

		mergeAffinityWithAnchorPod(p.Spec.Affinity, anchorPod)
		return p, nil
	}
}

func mergeAffinityWithAnchorPod(affinity *corev1.Affinity, anchorPodName string) {
	podAffinityTerm := podAffinityTermUsingAnchorPod(anchorPodName)

	if affinity.PodAffinity == nil {
		affinity.PodAffinity = &corev1.PodAffinity{}
	}

	affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution =
		append(affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution, *podAffinityTerm)
}

func podAffinityTermUsingAnchorPod(anchorPodName string) *corev1.PodAffinityTerm {
	return &corev1.PodAffinityTerm{LabelSelector: &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"anchor-pod": anchorPodName,
		},
	},
		TopologyKey: "kubernetes.io/hostname",
	}
}
