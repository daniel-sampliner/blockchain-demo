<!--
SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>

SPDX-License-Identifier: GLWTPL
-->

# blockchain-demo

## Prerequisites

Necessary tools are contained and managed with Nix in [default.nix](./default.nix).
Refer to [upstream documentation](https://nixos.org/download/) for instructions to install.

If not using Nix, manually install:

* [GNU Make](https://www.gnu.org/software/make/)
* [kind](https://kind.sigs.k8s.io/)
* [redo](https://redo.readthedocs.io/en/latest/)

## Instructions

1. Create kind kubernetes cluster with command:

    ```bash
    :; kind create cluster --config config.yaml --wait 5m
    ```

1. Load all Kubernetes manifests:

    ```bash
    :; redo -j$(nproc)
    ```

1. Wait a couple minutes for the cluster to stabilize.

1. Port-forward to the racecourse service:

    ```bash
    :; kubectl -n racecourse port-forward services/racecourse 3000:webapp
    ```

1. Open http://localhost:8080 in web browser.

1. Provide node connection URL: `http://signer.besu-qbft:8545`

1. 🏇🏇🏇
