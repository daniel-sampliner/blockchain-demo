#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

set -e

kubectl() {
	command kubectl --namespace besu-qbft "$@"
}

id=0
eth_api() {
	((++id))
	jq -n \
		--arg method "$1" \
		--arg id "$id" \
		'
			.jsonrpc = "2.0"
			| .method = $method
			| .id = ($id | tonumber)
			| .params = input
		' \
		| curl -sSf \
			--retry 10 \
			--retry-max-time 300 \
			--max-time 10 \
			blockchain.liger-beaver.ts.net --json @-
}

redo-ifchange apply-k8s

sleep 1
kubectl wait \
	--for condition=Ready \
	--selector app=besu \
	--timeout=5m \
	pod >&2

kubectl delete pod --selector app=signer >&2

kubectl wait \
	--for condition=Ready \
	--selector app=signer \
	--timeout=30s \
	pod >&2


# For some reason racecourse always times out deploying the first
# contract on greenfield deployment. We can avoid this by simply sending
# a dummy transaction first.

eth_api 'eth_blockNumber' <<<'[]' >&2

fakeContract=0x6060604052341561000f57600080fd5b60eb8061001d6000396000f300606060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063c6888fa1146044575b600080fd5b3415604e57600080fd5b606260048080359060200190919050506078565b6040518082815260200191505060405180910390f35b60007f24abdb5865df5079dcc5ac590ff6f01d5c16edbc5fab4e195d9febd1114503da600783026040518082815260200191505060405180910390a16007820290509190505600a165627a7a7230582040383f19d9f65246752244189b02f56e8d0980ed44e7a56c0b200458caad20bb0029
eth_api 'eth_accounts' <<<'[]' \
	| jq -n \
		--arg data "$fakeContract" \
		'[(.from = input.result[0] | .data = $data | .gas = "0x20000")]' \
	| eth_api 'eth_sendTransaction' >&2
