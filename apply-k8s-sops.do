#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

kustomize build --enable-alpha-plugins --enable-exec k8s/tsnsrv \
	| kubectl apply -f - >&2
