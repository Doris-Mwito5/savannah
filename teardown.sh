#!/bin/bash

echo "🧹 Cleaning up Savannah POS deployment..."

kubectl delete -f k8s/ingress.yaml --ignore-not-found=true
kubectl delete -f k8s/deployment.yaml --ignore-not-found=true
kubectl delete -f k8s/migration-job.yaml --ignore-not-found=true
kubectl delete -f k8s/postgresql.yaml --ignore-not-found=true
kubectl delete -f k8s/configmap.yaml --ignore-not-found=true
kubectl delete -f k8s/secrets.yaml --ignore-not-found=true
kubectl delete -f k8s/namespace.yaml --ignore-not-found=true

echo "✅ Cleanup completed!"