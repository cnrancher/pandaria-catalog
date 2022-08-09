#!/bin/bash
set -e

echo "Branch name: $DRONE_SOURCE_BRANCH"
echo "Git http url: $DRONE_GIT_HTTP_URL"

git clone -b $DRONE_SOURCE_BRANCH $DRONE_GIT_HTTP_URL source && cd source

# Pull charts-build-scripts.
make pull-scripts
# Validation.
make validate
