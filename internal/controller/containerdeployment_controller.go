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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	containerv1 "github.com/CodingMonkeyN/container-as-a-service/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ContainerDeploymentReconciler reconciles a ContainerDeployment object
type ContainerDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.com.coding-monkey,resources=containerdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.com.coding-monkey,resources=containerdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.com.coding-monkey,resources=containerdeployments/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ContainerDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	// Fetch the ContainerDeployment instance
	var containerDeployment containerv1.ContainerDeployment
	if err := r.Get(ctx, req.NamespacedName, &containerDeployment); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if containerDeployment.Spec.Namespace != "" {
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: containerDeployment.Spec.Namespace,
			},
		}

		log.Info("Creating namespace")
		err := r.Client.Create(ctx, namespace)
		if err != nil && !errors.IsAlreadyExists(err) {
			return ctrl.Result{}, err
		}
		log.Info("Namespace created")
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      containerDeployment.Name,
			Namespace: containerDeployment.Spec.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": containerDeployment.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": containerDeployment.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  containerDeployment.Name,
							Image: containerDeployment.Spec.Image,

							Ports: []corev1.ContainerPort{
								{
									ContainerPort: containerDeployment.Spec.Port,
								},
							},
							Env: convertEnvMap(containerDeployment.Spec.EnvironmentVars),
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(containerDeployment.Spec.CPU),
									corev1.ResourceMemory: resource.MustParse(containerDeployment.Spec.Memory),
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.Create(ctx, deploy); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContainerDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&containerv1.ContainerDeployment{}).
		Complete(r)
}

func convertEnvMap(envMap map[string]string) []corev1.EnvVar {
	var envVars []corev1.EnvVar
	for name, value := range envMap {
		envVars = append(envVars, corev1.EnvVar{Name: name, Value: value})
	}

	return envVars
}
