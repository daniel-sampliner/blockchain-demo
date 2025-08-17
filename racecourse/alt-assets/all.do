#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo-ifchange \
	public/favicon.png \
	src/img/horse-1.png \
	src/img/horse-2.png \
	src/img/horse-3.png \
	src/img/horse-4.png \
	src/img/horse-5.png \
	src/img/loading.gif \
	src/img/splash.png \
	;
