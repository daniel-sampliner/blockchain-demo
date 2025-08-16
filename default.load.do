#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

CONTAINER_TOOL="${CONTAINER_TOOL:-docker}"
IMAGE="localhost/${2:?}:latest"

declare -A image_dir=(
	[racecourse]=racecourse
	[racecourse-operator]=racecourse/operator
	[loadbalancer]=LoadBalancer
)

"$CONTAINER_TOOL" build -t "$IMAGE" "${image_dir[$2]}" >&2
kind load docker-image "$IMAGE"

"$CONTAINER_TOOL" image inspect "$IMAGE" -f '{{.Id}}'
