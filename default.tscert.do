#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo apply-k8s-sops
kubectl --namespace tsnsrv wait \
	--for condition=Ready \
	--selector app=tsnsrv \
	pod >&2

ret=$(curl \
	--silent \
	--output /dev/null \
	--write-out '%{http_code}' \
	--head \
	--request TRACE \
	--max-time 10 \
	--retry 10 \
	--retry-max-time 60 \
	"https://$2")

[[ $ret == 405 ]]
