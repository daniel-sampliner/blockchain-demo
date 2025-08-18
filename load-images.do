#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo-always

deps=()

if [[ -v WITH_TS ]]; then
	deps+=(nuke-old-ts.build)
fi

deps+=(
	loadbalancer.build
	racecourse-alt.build
	racecourse-operator.build
	racecourse.build
)

redo-ifchange "${deps[@]}"

for dep in "${deps[@]}"; do
	read -r img _ <"$dep"
	kind load docker-image "$img"
	podman exec -it kind-control-plane \
		crictl inspecti --output go-template --template '{{index .status.repoTags 0}} {{.status.id}}' "$img" \
		| tee >(redo-stamp)
done
