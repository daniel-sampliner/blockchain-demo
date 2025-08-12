# SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
#
# SPDX-License-Identifier: GLWTPL

let
  pkgs = import <nixpkgs> { };
in
{
  shell = pkgs.mkShell {
    packages = builtins.attrValues {
      inherit (pkgs)
        kind
        reuse
        ;
    };
  };
}
