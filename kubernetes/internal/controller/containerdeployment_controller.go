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
	containerv1 "github.com/CodingMonkeyN/container-as-a-service/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"log"
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
// +kubebuilder:rbac:groups=apps.com.coding-monkey,resources=containerdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.com.coding-monkey,resources=containerdeployments/status,verbs=get
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=create;delete;list;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ContainerDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var containerDeployment containerv1.ContainerDeployment
	if err := r.Get(ctx, req.NamespacedName, &containerDeployment); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(containerDeployment.Name) > 15 {
		return ctrl.Result{}, errors.NewBadRequest("Name must be less than 16 characters")
	}

	deploymentError := createDeployment(r, containerDeployment, ctx)
	if deploymentError != nil {
		log.Println("Error creating deployment: ", deploymentError.Error())
		return ctrl.Result{}, deploymentError
	}

	backendPortName, serviceError := createService(r, containerDeployment, ctx)
	if serviceError != nil {
		log.Println("Error creating service: ", serviceError.Error())
		return ctrl.Result{}, serviceError
	}

	ingressError := createIngress(backendPortName, r, containerDeployment, ctx)
	if ingressError != nil {
		log.Println("Error creating ingress: ", ingressError.Error())
		return ctrl.Result{}, ingressError
	}

	return ctrl.Result{}, nil
}

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

	envVars = overrideDefaultEnvInjections(envVars)

	return envVars
}

func overrideDefaultEnvInjections(envVars []corev1.EnvVar) []corev1.EnvVar {
	envToOverride := map[string]string{
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
	for name, value := range envToOverride {
		envVars = append(envVars, corev1.EnvVar{Name: name, Value: value})
	}

	return envVars
}

func createDeployment(r *ContainerDeploymentReconciler,
	containerDeployment containerv1.ContainerDeployment,
	ctx context.Context,
) error {
	replicas := int32(1)
	if containerDeployment.Spec.Replicas != nil {
		replicas = *containerDeployment.Spec.Replicas
	}
	var mounts []corev1.VolumeMount
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      containerDeployment.Name,
			Namespace: containerDeployment.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(replicas),
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
					RuntimeClassName:   ptr.To("kata-qemu"),
					EnableServiceLinks: ptr.To(false),
					Containers: []corev1.Container{
						{
							Name:  containerDeployment.Name,
							Image: containerDeployment.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: containerDeployment.Spec.Port,
								},
							},
							ImagePullPolicy: corev1.PullAlways,
							VolumeMounts:    mounts,
							Env:             convertEnvMap(containerDeployment.Spec.EnvironmentVars),
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    containerDeployment.Spec.CPU,
									corev1.ResourceMemory: containerDeployment.Spec.Memory,
								},
							},
						},
					},
				},
			},
		},
	}

	var err error
	if err = r.Create(ctx, deploy); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	if err != nil && errors.IsAlreadyExists(err) {
		if updateError := r.Update(ctx, deploy); updateError != nil {
			return updateError
		}
	}

	return nil
}

func createService(r *ContainerDeploymentReconciler,
	containerDeployment containerv1.ContainerDeployment,
	ctx context.Context) (string, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      containerDeployment.Name,
			Namespace: containerDeployment.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: map[string]string{"app": containerDeployment.Name},
			Ports: []corev1.ServicePort{
				{
					Name:       containerDeployment.Name,
					Port:       80,
					TargetPort: intstr.FromInt32(containerDeployment.Spec.Port),
				},
			},
		},
	}

	var err error
	if err = r.Create(ctx, service); err != nil && !errors.IsAlreadyExists(err) {
		return "", err
	}

	if err != nil && errors.IsAlreadyExists(err) {
		if updateError := r.Update(ctx, service); updateError != nil {
			return "", updateError
		}
	}
	return containerDeployment.Name, nil
}

func createIngress(backendPortName string, r *ContainerDeploymentReconciler,
	containerDeployment containerv1.ContainerDeployment,
	ctx context.Context) error {
	pathType := networkingv1.PathTypePrefix
	ingressClassName := "traefik"
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      containerDeployment.Name,
			Namespace: containerDeployment.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class":    "traefik",
				"cert-manager.io/cluster-issuer": "lets-encrypt",
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClassName,
			TLS: []networkingv1.IngressTLS{
				{
					Hosts:      []string{containerDeployment.Name + ".codingmonkey.cloud"},
					SecretName: "codingmonkey-cloud-tls",
				},
			},
			Rules: []networkingv1.IngressRule{
				{
					Host: containerDeployment.Name + ".codingmonkey.cloud",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: containerDeployment.Name,
											Port: networkingv1.ServiceBackendPort{
												Name: backendPortName,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	var err error
	if err = r.Create(ctx, ingress); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	if err != nil && errors.IsAlreadyExists(err) {
		if updateError := r.Update(ctx, ingress); updateError != nil {
			return updateError
		}
	}
	return nil
}
