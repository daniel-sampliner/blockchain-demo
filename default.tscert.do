#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo apply-k8s-ingress

ret=$(curl \
	--silent \
	--output /dev/null \
	--write-out '%{http_code}' \
	--head \
	--request TRACE \
	"https://$2")

[[ $ret == 405 ]]
