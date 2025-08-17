# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

ARG VARIANT=standard

FROM docker.io/oven/bun:1-alpine as base
WORKDIR /usr/src/app


FROM base AS install
RUN apk add git

COPY client/bun.lock client/package.json /temp/client
RUN cd /temp/client && bun install --frozen-lockfile

COPY server/bun.lock server/package.json /temp/server
RUN cd /temp/server && bun install --frozen-lockfile


FROM docker.io/nixpkgs/nix:latest as alt_assets

WORKDIR /temp
COPY alt-assets /temp/alt-assets
RUN NIX_PATH=nixpkgs=channel:nixos-25.05 nix print-dev-env \
		--extra-experimental-features nix-command \
		-f alt-assets/default.nix > nix.env

RUN . ./nix.env \
	&& cd alt-assets \
	&& redo -j$(nproc)


FROM base as prerelease
ARG VARIANT

COPY --from=alt_assets /temp/alt-assets/public /temp/alt-assets/public
COPY --from=alt_assets /temp/alt-assets/src /temp/alt-assets/src
COPY --from=install /temp/client/node_modules client/node_modules
COPY --from=install /temp/server/node_modules server/node_modules
COPY . .

RUN cd client \
	&& if [ ${VARIANT} = alt ]; then \
		mv -v -- /temp/alt-assets/public/* public \
		&& mv -v -- /temp/alt-assets/src/img/* src/img \
		&& mv -v -- ../contracts/alt/* ../contracts \
	; fi \
	&& bun run build


FROM base as release
RUN apk add catatonit 

COPY --from=prerelease /usr/src/app /usr/src/app

RUN chown -R bun:bun .
USER bun

WORKDIR /usr/src/app/server
EXPOSE 3000/tcp
ENTRYPOINT [ "catatonit", "--", "bun", "server.js" ]
