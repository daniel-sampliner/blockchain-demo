# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

FROM docker.io/oven/bun:1-alpine as base
WORKDIR /usr/src/app

FROM base as release
RUN apk add \
	catatonit \
	git

COPY . .

RUN chown -R bun:bun .
USER bun

RUN cd client \
	&& bun install --frozen-lockfile \
	&& bun run build

RUN cd server \
	&& bun install --frozen-lockfile

WORKDIR /usr/src/app/server
EXPOSE 3000/tcp
ENTRYPOINT [ "catatonit", "--", "bun", "server.js" ]
