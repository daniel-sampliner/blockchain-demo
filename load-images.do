#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

deps=(
	loadbalancer.load
	racecourse-alt.load
	racecourse-operator.load
	racecourse.load
)

if [[ -v WITH_TS ]]; then
	deps+=(
		nuke-old-ts.load
	)
fi
redo "${deps[@]}"
