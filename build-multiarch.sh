#!/bin/bash

# Multi-architecture build script for TutupLapak API
# This script builds Docker images for both AMD64 and ARM64 architectures

set -e

echo "🏗️  Building TutupLapak API for multiple architectures..."

# Build for AMD64 (Intel/AMD)
echo "📦 Building for AMD64 (linux/amd64)..."
docker buildx build --platform linux/amd64 -t tutuplapak-api:amd64 -f Dockerfile.multiarch .

# Build for ARM64 (Apple Silicon, ARM servers)
echo "📦 Building for ARM64 (linux/arm64)..."
docker buildx build --platform linux/arm64 -t tutuplapak-api:arm64 -f Dockerfile.multiarch .

# Create a multi-architecture manifest
echo "🔗 Creating multi-architecture manifest..."
docker buildx build --platform linux/amd64,linux/arm64 -t tutuplapak-api:latest -f Dockerfile.multiarch .

echo "✅ Multi-architecture build complete!"
echo ""
echo "Available images:"
echo "  - tutuplapak-api:amd64  (Intel/AMD)"
echo "  - tutuplapak-api:arm64  (Apple Silicon/ARM)"
echo "  - tutuplapak-api:latest (Multi-architecture)"
echo ""
echo "To run on your current architecture:"
echo "  docker run -p 8080:8080 tutuplapak-api:latest"
