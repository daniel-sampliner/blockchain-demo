#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

declare -A target=(
	[src/img/loading.gif]=01_specialweek.DtnlpFlH.webp
)
[[ -v target[$1] ]]

redo-ifchange download

mkdir -p "${1%/*}"
src="download.dir/${target[$1]}"
magick "$src" "gif:$3"
