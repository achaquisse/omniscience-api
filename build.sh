#!/bin/bash

set -e

echo "Building Omniscience API..."

cd src

echo "Downloading dependencies..."
go mod download

echo "Building binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../bin/omniscience-api .

echo "✅ Build completed successfully!"
echo "Binary location: bin/omniscience-api"
