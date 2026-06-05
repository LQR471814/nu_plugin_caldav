{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      packages.${system}.default = pkgs.buildGoModule {
        name = "nu_plugin_caldav";
        system = builtins.currentSystem;
        meta.mainProgram = "nu_plugin_caldav";

        src = ./.;

        vendorHash = "sha256-6AhxJPXHRpCRun4VaHGmT85ZBLR0fWHL0IKGVH+zaow=";
        subPackages = [ "." ];
      };

      apps.${system}.default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/nu_plugin_caldav";
      };
    };
}
