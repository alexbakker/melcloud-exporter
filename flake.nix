{
  description = "Nix flake for melcloud-exporter";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, flake-utils, nixpkgs }: 
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in
      rec {
        defaultPackage = with pkgs;
          buildGoModule rec {
            name = "melcloud-exporter";
            src = self;

            vendorHash = "sha256-8GC9/rtfqiXvR1goEM+FR5n9wK3YNs2F/3TB7NXvkXw=";
          };
        devShell = with pkgs; mkShell {
          buildInputs = [
            go
          ];
        };
      }
    );
}
