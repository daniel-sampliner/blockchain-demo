<!--
SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>

SPDX-License-Identifier: GLWTPL
-->

# k8s

## besu-qbft

Runs the private blockchain.

### besu

Deploys multiple `StatefulSet`s of [`besu`](https://besu.hyperledger.org/) pods that run the blockchain network.

- [`besu/base`](./besu-qbft/besu/base): common YAML manifests and `besu` configuration for all nodes.
- [`besu/overlays`](./besu-qbft/besu/overlays): node-specific configuration.

#### Scale-up

1. Duplicate the current highest numbered node's overlay directory, incrementing its number.

    For example:

    ```bash
    :; cp -r besu/overlays/node-3 besu/overlays/node-4
    ```

1. Edit the node's `kustomization.yaml` to increment the node number.
1. Edit the node's `secrets.yaml` with a new public/private keypair.
1. Edit [`besu/overlays/kustomization.yaml`](./besu-qbft/besu/overlays/kustomization.yaml) to include the new overlay directory.
1. Apply the top-level `besu` `kustomization.yaml`.

    For example:

    ```bash
    :; kubectl apply -k ./besu
    ```

### firefly

Deploys a `Deployment` of [`firefly-signer`](https://github.com/hyperledger/firefly-signer).
This is responsible for signing ethereum transactions.

Additional ethereum keys/passwords can be added into the keystore by:

1. Place private key.json and password files into [firefly/keys](./besu-qbft/firefly/keys).
1. Update [`firefly/kustomization.yaml`](./besu-qbft/firefly/kustomization.yaml) with new files.

### loadbalancer

Deploys a `Deployment` of [`LoadBalancer`](https://github.com/OnGridSystems/LoadBalancer/).
This is responsible for statefully load balancing json-rpc requests between the `besu` nodes.

## racecourse

Deploys:

1. `racecourse-controller` `Deployment`: Kubernetes Operator responsible for instantiating the Racecourse application.
1. Multiple `Racecourse` CRDs: each of these is a separate instance of the Racecourse application.

    Each instance can be configured with:

    ```
    FIELDS:
       image	<string>
         Image is the container image to use for the deployment

       ingressHost	<string>
         IngressHost is the host to use in the Ingress resource

       replicas	<integer>
         Replicas is the number of replicas in the deployment
    ```
