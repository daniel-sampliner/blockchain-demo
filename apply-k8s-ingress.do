#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e
cd k8s/ingress-nginx/
kustomize build --enable-alpha-plugins --enable-exec . \
	| kubectl apply -f - >&2
