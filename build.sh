#!/bin/sh

set -e # Blow up the second shit goes south

echo "Copying .env to src/config/.env..."
cp .env src/config/.env

echo "Building zeroctl..."
go build -o zeroctl src/main.go

echo "Moving binary to /usr/local/bin/..."
sudo mv zeroctl /usr/local/bin/

echo "Setting executable permissions..."
sudo chmod +x /usr/local/bin/zeroctl

echo "Stopping current zeroctl process..."
pkill -f "zeroctl daemon"

echo "Starting zeroctl..."
/usr/local/bin/zeroctl daemon &

echo "zeroctl restarted successfully."
