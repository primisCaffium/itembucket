#!/bin/bash
set -e

VERSION="1.0.0"
COMMIT="$(git rev-parse --short HEAD)"
BUILDDATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

go build -ldflags "-X main.version=${VERSION} -X main.commitHash=${COMMIT} -X main.buildDate=${BUILDDATE}" -o ib
mkdir -p ~/.itembucket/
mv ib ~/.itembucket/

echo "Itembucket installed successfully!"
echo ""
echo ">>> You need to add ~/.itembucket/ to your path before continuing <<<"
echo "Like so:"
echo "    export PATH=\"\$PATH\":\"~/.itembucket"
