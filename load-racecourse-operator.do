#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

CONTAINER_TOOL="${CONTAINER_TOOL:-docker}"
IMAGE="${IMAGE:-localhost/racecourse-operator:latest}"

"$CONTAINER_TOOL" build -t "$IMAGE" racecourse/operator >&2
kind load docker-image "$IMAGE"

"$CONTAINER_TOOL" image inspect "$IMAGE" -f '{{.Id}}'
