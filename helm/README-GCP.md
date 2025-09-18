# Google Cloud Platform (GCP) Integration for Kubernetes

This document explains how to set up Google Cloud Platform integration for the TutupLapak application running on Kubernetes.

## Overview

The application is configured to use Google Cloud services through a service account key that is securely mounted as a Kubernetes secret. The service account key is base64 encoded and stored in the Helm values file.

## Setup Instructions

### 1. Prepare the Service Account Key

1. Ensure your `ps-key.json` file is in the project root directory
2. The file should contain your Google Cloud service account credentials

### 2. Encode and Configure the Key

Run the setup script to automatically encode the service account key and update the Helm values:

```bash
./scripts/setup-gcp-key.sh
```

This script will:
- Encode the `ps-key.json` file to base64
- Update `helm/tutuplapak/values.yaml` with the encoded key
- Provide instructions for deployment

### 3. Manual Setup (Alternative)

If you prefer to do it manually:

```bash
# Encode the service account key
cat ps-key.json | base64 -w 0

# Copy the output and update helm/tutuplapak/values.yaml
# Replace the empty value for gcp-service-account-key with the encoded string
```

### 4. Deploy to Kubernetes

```bash
# For development
helm upgrade --install tutuplapak ./helm/tutuplapak

# For production
helm upgrade --install tutuplapak ./helm/tutuplapak -f ./helm/tutuplapak/values-production.yaml
```

## Configuration Details

### Kubernetes Secret

The service account key is stored as a Kubernetes secret with the following structure:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tutuplapak-secrets
type: Opaque
data:
  gcp-service-account-key: <base64-encoded-json-key>
```

### Volume Mount

The service account key is mounted in the pod at `/etc/gcp/ps-key.json`:

```yaml
volumeMounts:
  - name: gcp-service-account
    mountPath: /etc/gcp
    readOnly: true
```

### Environment Variable

The application uses the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to locate the service account key:

```yaml
env:
  - name: GOOGLE_APPLICATION_CREDENTIALS
    value: "/etc/gcp/ps-key.json"
```

## Security Considerations

1. **Never commit the `ps-key.json` file** - It's already in `.gitignore`
2. **The encoded key in values.yaml should be committed** - This is the intended way to distribute the key
3. **Rotate the service account key regularly** - Update the values.yaml file when rotating keys
4. **Use least privilege principle** - Only grant necessary permissions to the service account

## Using Google Cloud Services

Once deployed, your application can use Google Cloud services by:

1. **Using the Google Cloud client libraries** with Application Default Credentials
2. **Setting the environment variable** `GOOGLE_APPLICATION_CREDENTIALS` to `/etc/gcp/ps-key.json`
3. **The client libraries will automatically discover and use the credentials**

### Example Usage in Go

```go
import (
    "context"
    "cloud.google.com/go/storage"
    "google.golang.org/api/option"
)

// Initialize Google Cloud Storage client
ctx := context.Background()
client, err := storage.NewClient(ctx, option.WithCredentialsFile("/etc/gcp/ps-key.json"))
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## Troubleshooting

### Common Issues

1. **Permission denied errors**: Ensure the service account has the necessary IAM roles
2. **Key not found**: Verify the volume mount and secret configuration
3. **Invalid credentials**: Check that the base64 encoding is correct

### Verification

To verify the setup is working:

```bash
# Check if the secret exists
kubectl get secret tutuplapak-secrets

# Check if the volume is mounted correctly
kubectl exec -it <pod-name> -- ls -la /etc/gcp/

# Verify the credentials file
kubectl exec -it <pod-name> -- cat /etc/gcp/ps-key.json
```

## Production Considerations

For production deployments:

1. **Use Workload Identity** instead of service account keys when possible
2. **Implement proper secret rotation** procedures
3. **Monitor access logs** for the service account
4. **Use separate service accounts** for different environments
5. **Consider using Google Secret Manager** for more advanced secret management

## Files Modified

- `helm/tutuplapak/values.yaml` - Added GCP service account key configuration
- `helm/tutuplapak/values-production.yaml` - Added GCP configuration for production
- `helm/tutuplapak/templates/deployment.yaml` - Added volume mount and environment variable
- `scripts/setup-gcp-key.sh` - Created setup script for encoding the key
- `.gitignore` - Added `ps-key.json` to prevent accidental commits
