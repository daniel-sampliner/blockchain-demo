#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

redo \
	apply-k8s \
	apply-k8s-sops \
	load-images \
	blockchain.liger-beaver.ts.net.tscert \
	racecourse.liger-beaver.ts.net.tscert \
	;
