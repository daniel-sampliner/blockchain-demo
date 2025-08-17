#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo-ifchange apply-k8s

sleep 1
kubectl --namespace besu-qbft wait \
	--for condition=Ready \
	--selector app=besu \
	--timeout=2m \
	pod >&2

kubectl --namespace besu-qbft delete \
	pod -l 'app=signer' >&2
