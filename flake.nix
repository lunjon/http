{
  description = "Flake for building and developing http cli";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        formatter = pkgs.nixfmt-rfc-style;
        packages.default = pkgs.buildGo123Module {
          name = "http";
          src = ./.;
          vendorHash = "sha256-h27uHmOQMECkGHFsDggGfm+hRohTVYIkvF7zWFdwlTM=";
          doCheck = false;
        };
      }
    );
}
