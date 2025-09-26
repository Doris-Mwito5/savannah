#!/bin/bash

set -e

echo "🚀 Starting Savannah POS Kubernetes Deployment..."

# Build the Docker image
echo "📦 Building Docker image..."
docker build -t savannah-pos:latest .

# Load image into Minikube
echo "⬆️ Loading image into Minikube..."
minikube image load savannah-pos:latest

# Create namespace
echo "📁 Creating namespace..."
kubectl apply -f k8s/namespace.yaml

# Apply configurations
echo "⚙️ Applying configurations..."
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml

# Deploy PostgreSQL
echo "🐘 Deploying PostgreSQL..."
kubectl apply -f k8s/postgresql.yaml

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=db-postgresql -n savannah-pos --timeout=180s

# Run database migrations
echo "🔄 Running database migrations..."
kubectl apply -f k8s/migration-job.yaml
kubectl wait --for=condition=complete job/db-migration -n savannah-pos --timeout=120s

# Deploy application
echo "🚀 Deploying application..."
kubectl apply -f k8s/deployment.yaml

# Deploy ingress
echo "🌐 Deploying ingress..."
kubectl apply -f k8s/ingress.yaml

# Wait for application to be ready
echo "⏳ Waiting for application to be ready..."
kubectl wait --for=condition=ready pod -l app=savannah-pos -n savannah-pos --timeout=120s

echo "✅ Deployment completed successfully!"
echo ""
echo "📊 Deployment status:"
kubectl get all -n savannah-pos

echo ""
echo "🌐 To access your application:"
echo "1. Add this to your /etc/hosts file:"
echo "   $(minikube ip) savannah-pos.local"
echo ""
echo "2. Access the application at: http://savannah-pos.local"
echo ""
echo "🔍 Check logs with: kubectl logs -f deployment/savannah-pos -n savannah-pos"