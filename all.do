#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

deps=(
	apply-k8s
	load-images
	wait-besu
)

if [[ -v WITH_TS ]]; then
	deps+=(
		apply-k8s-sops
		blockchain.liger-beaver.ts.net.tscert
		racecourse-alt.liger-beaver.ts.net.tscert
		racecourse.liger-beaver.ts.net.tscert
	)
fi

redo-ifchange "${deps[@]}"
