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

            vendorSha256 = "sha256-5iXY9UBuVVEbxuJSB6uuwGfvaDlpX/ux3gNJKOszIbw=";
          };
        devShell = with pkgs; mkShell {
          buildInputs = [
            go
          ];
        };
      }
    );
}
