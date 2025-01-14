#!/bin/bash

# create namespace and container
kubectl apply -f namespace.yaml --kubeconfig ../../../../kubeconfig
kubectl apply -f pod_kata.yaml --kubeconfig ../../../../kubeconfig
kubectl apply -f pod.yaml --kubeconfig ../../../../kubeconfig

# wait for containers to start
kubectl wait --timeout=10m --for=condition=Ready -f pod.yaml --kubeconfig ../../../../kubeconfig
kubectl wait --timeout=10m --for=condition=Ready -f pod_kata.yaml --kubeconfig ../../../../kubeconfig

# check normal pod
PS_AUX="$(kubectl exec sec-test -n sec-test --kubeconfig ../../../../kubeconfig -- ps aux | wc -l)"
IFCONFIG="$(kubectl exec sec-test -n sec-test --kubeconfig ../../../../kubeconfig -- ifconfig | wc -l)"

# check kata pod
PS_AUX_KATA="$(kubectl exec sec-test-kata -n sec-test --kubeconfig ../../../../kubeconfig -- ps aux | wc -l)"
IFCONFIG_KATA="$(kubectl exec sec-test-kata -n sec-test --kubeconfig ../../../../kubeconfig -- ifconfig | wc -l)"

if (( PS_AUX_KATA >= PS_AUX ));
then
    echo "Non kata pod has less or equal running processes. This shpuld not happen."
    echo "Kata pod: $PS_AUX_KATA"
    echo "Normal pod: $PS_AUX"
fi

if (( IFCONFIG_KATA >= IFCONFIG ));
then
    echo "Non kata pod has less or equal network configurations. This should not happen."
    echo "Kata pod: $IFCONFIG_KATA"
    echo "Normal pod: $IFCONFIG"
fi

# delete namespace with containers
kubectl delete -f namespace.yaml --kubeconfig ../../../../kubeconfig