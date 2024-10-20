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

package controller

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "github.com/CodingMonkeyN/container-as-a-service/api/v1" // Dein CRD-Paket
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("ContainerDeployment Controller", Ordered, func() {
	var (
		containerDeployment *v1.ContainerDeployment
	)

	integrationTestNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "integration-test",
		},
	}

	deploymentNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}

	BeforeAll(func() {
		containerDeployment = &v1.ContainerDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-container",
				Namespace: integrationTestNs.Name,
			},
			Spec: v1.ContainerDeploymentSpec{
				Image: "nginx:latest",
				Port:  80,
				Memory: resource.Quantity{
					Format: "64Mi",
				},
				CPU: resource.Quantity{
					Format: "250m",
				},
				Replicas: ptr.To(int32(2)),
				EnvironmentVars: map[string]string{
					"TEST_KEY": "TEST_VALUE",
				},
			},
		}

		Expect(k8sClient.Create(ctx, integrationTestNs)).To(Succeed())

		Expect(k8sClient.Create(ctx, containerDeployment)).To(Succeed())
	})

	AfterAll(func() {
		_ = k8sClient.Delete(ctx, integrationTestNs)

		deleteResource(&appsv1.Deployment{}, containerDeployment)
		deleteResource(&corev1.Service{}, containerDeployment)
		deleteResource(&networkingv1.Ingress{}, containerDeployment)
		_ = k8sClient.Delete(ctx, deploymentNs)
	})

	Context("When ContainerDeployment is created", func() {
		It("should create the deployment", func() {
			Eventually(func() error {
				deployment := &appsv1.Deployment{}
				return k8sClient.Get(ctx, client.ObjectKey{
					Name:      containerDeployment.Name,
					Namespace: containerDeployment.Namespace,
				}, deployment)
			}, 10*time.Second, 100*time.Millisecond).Should(Succeed())
		})

		It("should have the correct replicas", func() {
			Eventually(func() *int32 {
				deployment := &appsv1.Deployment{}
				err := k8sClient.Get(ctx, client.ObjectKey{
					Name:      containerDeployment.Name,
					Namespace: containerDeployment.Namespace,
				}, deployment)
				if err != nil {
					return ptr.To(int32(0))
				}
				return ptr.To(deployment.Status.ReadyReplicas)
			}, 10*time.Second, 100*time.Millisecond).Should(Equal(containerDeployment.Spec.Replicas), "Expected number of ready replicas to match the spec")
		})

		It("should ensure the Pod has the correct runtimeClass", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{
					Name:      containerDeployment.Name,
					Namespace: containerDeployment.Namespace,
				}, deployment)
			}, 10*time.Second, 100*time.Millisecond).Should(Succeed())

			podList := &corev1.PodList{}
			Eventually(func() error {
				return k8sClient.List(ctx, podList, client.InNamespace(containerDeployment.Namespace), client.MatchingLabels(deployment.Spec.Selector.MatchLabels))
			}, 10*time.Second, 100*time.Millisecond).Should(Succeed())

			Expect(len(podList.Items)).To(BeNumerically(">", 0), "Expected at least one Pod to be created")
			// TODO: ENABLE ONCE KATAS IS INSTALLED
			/*
				for _, pod := range podList.Items {
					Expect(pod.Spec.RuntimeClassName).NotTo(BeNil(), "Expected runtimeClassName to be set")
					Expect(*pod.Spec.RuntimeClassName).To(Equal("kata-qemu"), "Expected runtimeClassName to match 'kata-qemu'")
				}*/
		})

		It("should ensure the Pod has the correct environment variables", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{
					Name:      containerDeployment.Name,
					Namespace: containerDeployment.Namespace,
				}, deployment)
			}, 10*time.Second, 100*time.Millisecond).Should(Succeed())

			podList := &corev1.PodList{}
			Eventually(func() error {
				return k8sClient.List(ctx, podList, client.InNamespace(containerDeployment.Namespace), client.MatchingLabels(deployment.Spec.Selector.MatchLabels))
			}, 10*time.Second, 100*time.Millisecond).Should(Succeed())

			Expect(len(podList.Items)).To(BeNumerically(">", 0), "Expected at least one Pod to be created")

			defaultEnvs := map[string]string{
				"KUBERNETES_SERVICE_PORT_HTTPS": "",
				"KUBERNETES_SERVICE_PORT":       "",
				"KUBERNETES_PORT_443_TCP":       "",
				"KUBERNETES_PORT_443_TCP_PROTO": "",
				"KUBERNETES_PORT_443_TCP_ADDR":  "",
				"KUBERNETES_SERVICE_HOST":       "",
				"KUBERNETES_PORT":               "",
				"KUBERNETES_PORT_443_TCP_PORT":  "",
				"HOSTNAME":                      "",
			}
			expectedEnv := containerDeployment.Spec.EnvironmentVars

			for _, pod := range podList.Items {
				for _, container := range pod.Spec.Containers {
					actualEnv := make(map[string]string)
					for _, env := range container.Env {
						actualEnv[env.Name] = env.Value
					}

					for name, value := range expectedEnv {
						Expect(actualEnv).To(HaveKeyWithValue(name, value), fmt.Sprintf("Expected environment variable %s to have value %s", name, value))
					}

					for name, value := range defaultEnvs {
						Expect(actualEnv).To(HaveKeyWithValue(name, value), fmt.Sprintf("Expected default environment variable %s to have value %s", name, value))
					}
				}
			}
		})

		It("should create the service", func() {
			Eventually(func() error {
				service := &corev1.Service{}
				return k8sClient.Get(ctx, client.ObjectKey{
					Name:      containerDeployment.Name,
					Namespace: containerDeployment.Namespace,
				}, service)
			}, 10*time.Second, 100*time.Millisecond).Should(Succeed())
		})

		It("should create the ingress", func() {
			Eventually(func() error {
				ingress := &networkingv1.Ingress{}
				return k8sClient.Get(ctx, client.ObjectKey{
					Name:      containerDeployment.Name,
					Namespace: containerDeployment.Namespace,
				}, ingress)
			}, 10*time.Second, 100*time.Millisecond).Should(Succeed())
		})
	})
})

func deleteResource(resource client.Object, containerDeployment *v1.ContainerDeployment) {
	resource.SetName(containerDeployment.Name)
	resource.SetNamespace(containerDeployment.Namespace)
	_ = k8sClient.Delete(ctx, resource)
}
