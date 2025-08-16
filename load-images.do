#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo \
	loadbalancer.load \
	racecourse.load \
	racecourse-operator.load \
	;
