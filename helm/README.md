# TutupLapak Helm Chart

This Helm chart deploys the TutupLapak marketplace application with all its dependencies including PostgreSQL, MinIO, Prometheus, and Grafana.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- kubectl configured to connect to your Kubernetes cluster

## Quick Start

### 1. Build and Push Docker Image

First, build and push your TutupLapak application Docker image:

```bash
# Build the Docker image
docker build -t tutuplapak:latest .

# Tag for your registry (replace with your registry)
docker tag tutuplapak:latest your-registry.com/tutuplapak:latest

# Push to registry
docker push your-registry.com/tutuplapak:latest
```

### 2. Deploy with Default Values

```bash
# Deploy to default namespace
helm install tutuplapak ./tutuplapak

# Deploy to specific namespace
helm install tutuplapak ./tutuplapak -n tutuplapak --create-namespace
```

### 3. Deploy with Custom Values

```bash
# Deploy with custom values file
helm install tutuplapak ./tutuplapak -f custom-values.yaml

# Deploy with inline values
helm install tutuplapak ./tutuplapak --set app.replicaCount=3 --set app.image.tag=v1.0.0
```

### 4. Using the Deployment Script

```bash
# Deploy with default settings
./deploy.sh

# Deploy with custom release name and namespace
./deploy.sh --release my-tutuplapak --namespace my-namespace

# Deploy with custom values file
./deploy.sh --values custom-values.yaml

# Dry run to see what would be deployed
./deploy.sh --dry-run

# Upgrade existing deployment
./deploy.sh --upgrade
```

## Configuration

### Values File Structure

The main configuration is in `values.yaml`. Key sections:

#### Application Configuration
```yaml
app:
  replicaCount: 3
  image:
    repository: tutuplapak
    tag: "latest"
  service:
    type: ClusterIP
    port: 8080
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 128Mi
```

#### Database Configuration
```yaml
database:
  postgresql:
    enabled: true
    auth:
      postgresPassword: "postgres_password"
      database: "tutuplapak"
      username: "tutuplapak"
    primary:
      persistence:
        enabled: true
        size: 10Gi
```

#### MinIO Configuration
```yaml
minio:
  enabled: true
  auth:
    rootUser: "minioadmin"
    rootPassword: "minioadmin"
  defaultBuckets: "tutuplapak-uploads"
  persistence:
    enabled: true
    size: 50Gi
```

#### Monitoring Configuration
```yaml
prometheus:
  enabled: true
  server:
    persistentVolume:
      enabled: true
      size: 10Gi

grafana:
  enabled: true
  adminPassword: "admin"
  persistence:
    enabled: true
    size: 5Gi
```

## Environment Variables

The application uses the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Auto-generated |
| `JWT_SECRET` | JWT signing secret | From secret |
| `JWT_EXPIRY` | JWT token expiry | 24h |
| `MINIO_ENDPOINT` | MinIO server endpoint | Auto-generated |
| `MINIO_ACCESS_KEY` | MinIO access key | minioadmin |
| `MINIO_SECRET_KEY` | MinIO secret key | From secret |
| `MINIO_BUCKET` | MinIO bucket name | tutuplapak-uploads |
| `PROMETHEUS_ENABLED` | Enable Prometheus metrics | true |
| `METRICS_PORT` | Metrics port | 9090 |

## Services and Ports

| Service | Port | Description |
|---------|------|-------------|
| TutupLapak API | 8080 | Main application API |
| TutupLapak Metrics | 9090 | Prometheus metrics endpoint |
| PostgreSQL | 5432 | Database |
| MinIO API | 9000 | Object storage API |
| MinIO Console | 9001 | Object storage web UI |
| Prometheus | 9090 | Monitoring |
| Grafana | 3000 | Dashboards |

## Ingress

The chart includes optional ingress configuration:

```yaml
app:
  ingress:
    enabled: true
    className: "nginx"
    hosts:
      - host: tutuplapak.local
        paths:
          - path: /api
            pathType: Prefix
          - path: /metrics
            pathType: Prefix
```

## Monitoring

### Prometheus Metrics

The application exposes Prometheus metrics on port 9090:

- `tutuplapak_http_requests_total` - HTTP request counter
- `tutuplapak_http_request_duration_seconds` - HTTP request duration histogram
- `tutuplapak_database_operations_total` - Database operation counter
- `tutuplapak_upload_requests_total` - File upload counter
- `tutuplapak_memory_usage_bytes` - Memory usage gauge

### Grafana Dashboards

Grafana is automatically configured with:
- Prometheus as data source
- Pre-built dashboards for TutupLapak monitoring
- Alerting rules for common issues

Access Grafana:
```bash
kubectl port-forward svc/tutuplapak-grafana 3000:3000 -n tutuplapak
# Open http://localhost:3000
# Username: admin, Password: admin
```

## Scaling

### Horizontal Pod Autoscaling

To enable HPA, add to your values:

```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
  targetMemoryUtilizationPercentage: 80
```

### Vertical Scaling

Adjust resource limits in values.yaml:

```yaml
app:
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 200m
      memory: 256Mi
```

## Security

### Secrets Management

Sensitive data is stored in Kubernetes secrets:

```yaml
secrets:
  create: true
  data:
    jwt-secret: "your-jwt-secret-key-change-in-production"
    db-password: "postgres_password"
    minio-secret: "minioadmin"
```

### Network Policies

Enable network policies for additional security:

```yaml
networkPolicy:
  enabled: true
  ingress: []
  egress: []
```

## Troubleshooting

### Common Issues

1. **Pod not starting**
   ```bash
   kubectl describe pod <pod-name> -n tutuplapak
   kubectl logs <pod-name> -n tutuplapak
   ```

2. **Database connection issues**
   ```bash
   kubectl logs deployment/tutuplapak -n tutuplapak | grep -i database
   ```

3. **MinIO connection issues**
   ```bash
   kubectl logs deployment/tutuplapak -n tutuplapak | grep -i minio
   ```

### Useful Commands

```bash
# Check all resources
kubectl get all -n tutuplapak

# View logs
kubectl logs -f deployment/tutuplapak -n tutuplapak

# Port forward for testing
kubectl port-forward svc/tutuplapak 8080:8080 -n tutuplapak

# Check ingress
kubectl get ingress -n tutuplapak

# Check persistent volumes
kubectl get pv,pvc -n tutuplapak

# Check secrets
kubectl get secrets -n tutuplapak

# Check configmaps
kubectl get configmaps -n tutuplapak
```

## Upgrading

### Upgrade Application

```bash
# Update image tag
helm upgrade tutuplapak ./tutuplapak --set app.image.tag=v1.1.0

# Or with values file
helm upgrade tutuplapak ./tutuplapak -f new-values.yaml
```

### Upgrade Dependencies

```bash
# Update Helm dependencies
helm dependency update ./tutuplapak

# Upgrade with updated dependencies
helm upgrade tutuplapak ./tutuplapak
```

## Uninstalling

```bash
# Uninstall release
helm uninstall tutuplapak -n tutuplapak

# Or using the script
./deploy.sh --uninstall
```

## Development

### Local Development

For local development, you can use the chart with local values:

```yaml
# local-values.yaml
app:
  image:
    repository: tutuplapak
    tag: "dev"
  replicaCount: 1

database:
  postgresql:
    enabled: true
    auth:
      postgresPassword: "dev_password"

minio:
  enabled: true
  auth:
    rootPassword: "dev_password"

prometheus:
  enabled: false

grafana:
  enabled: false
```

Deploy with:
```bash
helm install tutuplapak-dev ./tutuplapak -f local-values.yaml
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the chart
5. Submit a pull request

## License

This project is licensed under the MIT License.
