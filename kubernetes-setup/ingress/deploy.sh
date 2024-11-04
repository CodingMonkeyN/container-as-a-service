#!/bin/sh

kubectl apply -f namespace.yaml
kubectl apply -f role.yaml
kubectl apply -f account.yaml
kubectl apply -f role-binding.yaml
kubectl apply -f deployment.yaml



