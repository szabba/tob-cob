{
  description = "A Nix-flake-based Go development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = { allowUnfree = true; };
        };
      in {
        devShells = {
          default = pkgs.mkShell {
            # Packages included in the environment
            buildInputs = with pkgs; [
              # Nix
              nixfmt

              # Git
              gitMinimal

              # C
              gcc
              pkg-config

              xorg.libX11
              xorg.libXcursor
              xorg.libXrandr
              xorg.libXinerama
              xorg.libXi
              mesa
              glew

              # Go
              go
              gotools
              golangci-lint
              gopls
              go-outline
              gopkgs

              # Editors
              vscode

              # CI
              earthly
            ];

            # Run when the shell is started up
            shellHook = ''
              ${pkgs.go}/bin/go version
            '';
          };
        };
      });
}
