#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo \
	loadbalancer.load \
	nuke-old-ts.load \
	racecourse-operator.load \
	racecourse.load \
	;
