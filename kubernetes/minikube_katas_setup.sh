minikube start --vm-driver kvm2 --memory 6144 --network-plugin=cni --enable-default-cni --container-runtime=cri-o --bootstrapper=kubeadm

git clone https://github.com/kata-containers/kata-containers.git
cd kata-containers/tools/packaging/kata-deploy
kubectl apply -f kata-rbac/base/kata-rbac.yaml
kubectl apply -f kata-deploy/base/kata-deploy.yaml