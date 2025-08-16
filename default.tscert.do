#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo apply-k8s-sops

sleep 1
kubectl --namespace tsnsrv wait \
	--for condition=Ready \
	--selector app=tsnsrv \
	--selector instance="tsnsrv-${2%%.*}" \
	pod >&2

ret=$(curl \
	--silent \
	--output /dev/null \
	--write-out '%{http_code}' \
	--head \
	--request TRACE \
	--max-time 30 \
	--retry 10 \
	--retry-max-time 300 \
	-k \
	"https://$2")

[[ $ret == 405 ]]
