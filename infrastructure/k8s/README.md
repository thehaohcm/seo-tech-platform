# Kubernetes Manifests for SEO Tech Platform

This directory contains Kubernetes manifests for deploying the SEO Tech Platform to AWS EKS.

## Structure

```
k8s/
├── namespace.yaml              - Namespace definition
├── configmap.yaml              - Configuration
├── secrets.yaml                - Secrets (don't commit real values)
├── postgres.yaml               - PostgreSQL deployment
├── redis.yaml                  - Redis deployment
├── crawler-service.yaml        - Crawler service
├── analyzer-service.yaml       - Analyzer service
├── api-gateway.yaml            - API Gateway
├── web-dashboard.yaml          - Web Dashboard
└── ingress.yaml                - Ingress controller
```

## Deployment

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check status
kubectl get pods -n seo-platform

# View logs
kubectl logs -f <pod-name> -n seo-platform
```

## Prerequisites

- AWS EKS cluster running
- kubectl configured
- Docker images pushed to ECR
