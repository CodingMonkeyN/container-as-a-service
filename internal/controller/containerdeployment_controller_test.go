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
	"context"
	corev1 "k8s.io/api/core/v1"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1 "github.com/CodingMonkeyN/container-as-a-service/api/v1"
)

var _ = Describe("ContainerDeployment Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"
		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "test",
		}
		testNamespace := types.NamespacedName{
			Name:      "test",
			Namespace: "test-2",
		}
		containerdeployment := &appsv1.ContainerDeployment{}

		BeforeEach(func() {
			By("creating the test namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: typeNamespacedName.Namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("creating the custom resource for the Kind ContainerDeployment")
			err := k8sClient.Get(ctx, typeNamespacedName, containerdeployment)

			if err != nil && errors.IsNotFound(err) {
				resource := &appsv1.ContainerDeployment{
					ObjectMeta: metav1.ObjectMeta{
						Name: resourceName,
					},
					Spec: appsv1.ContainerDeploymentSpec{
						Image:     "nginx:latest",
						Namespace: testNamespace.Namespace,
						Port:      80,
						Memory:    "64Mi",
						CPU:       "250m",
						EnvironmentVars: map[string]string{
							"test-key": "test-value",
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ContainerDeploymentReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		var (
			interval = 100 * time.Millisecond
			timeout  = 5 * time.Second
		)
		It("should have successfully created the namespace", func() {
			Eventually(func() bool {
				ns := &corev1.Namespace{}
				err := k8sClient.Get(ctx, typeNamespacedName, ns)
				return err == nil
			}, timeout, interval).Should(BeTrue(), "Expected namespace to be created")
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &appsv1.ContainerDeployment{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())
			if err == nil {
				By("Cleanup the specific resource instance ContainerDeployment")
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}

			By("Cleanup the namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: typeNamespacedName.Namespace,
				},
			}
			err = k8sClient.Delete(ctx, ns)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
