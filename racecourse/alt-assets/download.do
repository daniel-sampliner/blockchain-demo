#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

redo-always

tempdir=$(mktemp -d "$3.XXXXXX")
readonly tempdir
trap 'rm -rf -- "$tempdir"' EXIT

(
	set -e
	cd "$tempdir"
	wcurl \
		https://umamusume.com/favicon.ico \
		https://umamusume.com/_app/immutable/assets/kv.BopeP3kA.avif \
		https://umamusume.com/_app/immutable/assets/01_specialweek.DtnlpFlH.webp \
		https://en-portal.g.kuroco-img.app/v=1744884538/files/user/character/specialweek/specialweek_list.png \
		https://en-portal.g.kuroco-img.app/v=1744944608/files/user/character/goldship/goldship_list.png \
		https://en-portal.g.kuroco-img.app/v=1744944706/files/user/character/mejiromcqueen/mejiromcqueen_list.png \
		https://en-portal.g.kuroco-img.app/v=1749211656/files/user/character/riceshower/riceshower_01_list.png \
		https://en-portal.g.kuroco-img.app/v=1744944639/files/user/character/elcondorpasa/elcondorpasa_list.png \
		;
	sha256sum ./* | tee >(redo-stamp)
)

rm -rf -- "$2.dir"
mv -T -- "$tempdir" "$2.dir"
