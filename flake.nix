{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nuconv = {
      url = "github:LQR471814/nuconv";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs =
    {
      self,
      nixpkgs,
      nuconv,
    }:
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

        vendorHash = "sha256-kIwF73rxRF+BJREwXHlf5wKzV9+s1Su4x4iA3s1SOuM=";
        subPackages = [ "." ];
      };

      apps.${system}.default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/nu_plugin_caldav";
      };

      devShells.${system}.default = pkgs.mkShell {
        name = "devenv";
        nativeBuildInputs = [
          nuconv.packages.${system}.default
        ];

        shellHook = ''
          echo "Devshell activated."
        '';
      };
    };
}
