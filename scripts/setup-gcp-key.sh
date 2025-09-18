#!/bin/bash

# Script to encode Google Cloud Service Account key for Kubernetes deployment
# Usage: ./scripts/setup-gcp-key.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if ps-key.json exists
if [ ! -f "ps-key.json" ]; then
    print_error "ps-key.json file not found in the current directory"
    print_warning "Please ensure ps-key.json is in the project root directory"
    exit 1
fi

# Encode the service account key
print_status "Encoding Google Cloud Service Account key..."
ENCODED_KEY=$(cat ps-key.json | base64 -w 0)

if [ -z "$ENCODED_KEY" ]; then
    print_error "Failed to encode the service account key"
    exit 1
fi

print_status "Service account key encoded successfully"

# Update the values.yaml file
print_status "Updating helm/tutuplapak/values.yaml with encoded key..."

# Create a temporary file with the updated values
TEMP_FILE=$(mktemp)
sed "s/gcp-service-account-key: \"\"/gcp-service-account-key: \"$ENCODED_KEY\"/" helm/tutuplapak/values.yaml > "$TEMP_FILE"

# Replace the original file
mv "$TEMP_FILE" helm/tutuplapak/values.yaml

print_status "Updated values.yaml successfully"
print_warning "Remember to commit the updated values.yaml file to your repository"
print_warning "The ps-key.json file should remain in .gitignore and not be committed"

print_status "Setup complete! You can now deploy with:"
echo "  helm upgrade --install tutuplapak ./helm/tutuplapak"
