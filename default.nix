# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

let
  pkgs = import <nixpkgs> { };
in
{
  shell =
    let
      ksops = pkgs.runCommand "kops" { } ''
        mkdir -p $out/bin
        ln -sf ${pkgs.kustomize-sops}/lib/viaduct.ai/v1/ksops/ksops $out/bin/ksops
      '';
    in
    pkgs.mkShell {
      packages = builtins.attrValues {
        inherit (pkgs)
          besu
          bun
          gnumake
          go
          go-ethereum
          kind
          kubebuilder
          redo-apenwarr
          reuse
          ;

        inherit ksops;
      };
    };
}
