# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

FROM docker.io/library/alpine:3 as base

RUN apk add bash curl jq

COPY nuke-old-ts /nuke-old-ts
ENTRYPOINT ["/nuke-old-ts"]
