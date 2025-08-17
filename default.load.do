#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

tag=latest
image="localhost/${2:?}"
extra_args=()
if [[ ${2:?} == racecourse-alt ]]; then
	tag=alt
	image="localhost/racecourse"
	extra_args+=(--build-arg VARIANT=alt)
fi

image="$image:$tag"

declare -A image_dir=(
	[loadbalancer]=LoadBalancer
	[nuke-old-ts]=nuke-old-ts
	[racecourse-alt]=racecourse
	[racecourse-operator]=racecourse-operator
	[racecourse]=racecourse
)

docker build -t "$image" "${extra_args[@]}" "${image_dir[$2]}" >&2
kind load docker-image "$image"

docker image inspect "$image" -f '{{.Id}}'
