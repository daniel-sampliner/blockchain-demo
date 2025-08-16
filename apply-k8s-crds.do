#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

make -C racecourse/operator generate manifests >&2
kubectl apply -k racecourse/operator/config/default >&2
