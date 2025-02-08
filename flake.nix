{
  description = "A very basic flake";
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
          # modRoot = "./tfn-go";
          vendorHash = "sha256-AVu+R0pvH1H6Mp4OyzM03P1LeOIRPW+HHU2sj88KqGo=";
          doCheck = false;
        };
      }
    );
}
