#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

deps=(
	apply-k8s
	load-images
)

if [[ -v WITH_TS ]]; then
	deps+=(
		apply-k8s-sops
		blockchain.liger-beaver.ts.net.tscert
		racecourse.liger-beaver.ts.net.tscert
	)
fi

redo "${deps[@]}"
