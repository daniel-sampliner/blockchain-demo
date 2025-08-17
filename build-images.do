#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

deps=(
	loadbalancer.build
	racecourse-alt.build
	racecourse-operator.build
	racecourse.build
)

if [[ -v WITH_TS ]]; then
	deps+=(nuke-old-ts.build)
fi

redo-ifchange "${deps[@]}"
