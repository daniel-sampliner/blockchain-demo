#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

declare -A target=(
	[public/favicon.png]=favicon.ico
	[src/img/horse-1.png]=specialweek_list.png
	[src/img/horse-2.png]=goldship_list.png
	[src/img/horse-3.png]=mejiromcqueen_list.png
	[src/img/horse-4.png]=riceshower_01_list.png
	[src/img/horse-5.png]=elcondorpasa_list.png
	[src/img/splash.png]=kv.BopeP3kA.avif
)
[[ -v target[$1] ]]

redo-ifchange download

mkdir -p "${1%/*}"
src="download.dir/${target[$1]}"
magick "$src" "png:$3"
