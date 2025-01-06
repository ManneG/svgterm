{
  description = "A CLI tool for creating svg terminal animations";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";
  inputs.systems.url = "github:nix-systems/default";

  outputs = { self, systems, nixpkgs }:
    let
      version = "0.0.1-dev";

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      eachSystem = nixpkgs.lib.genAttrs (import systems);

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = eachSystem (system: import nixpkgs { inherit system; });

    in {

      formatter = eachSystem (system:
        let pkgs = nixpkgsFor.${system};
        in pkgs.writeShellApplication {
          name = "format";
          runtimeInputs = [ pkgs.nixfmt-classic pkgs.go ];
          text = ''
            find . -name '*.nix' -exec nixfmt {} +
            gofmt -w .
          '';
        });

      packages = eachSystem (system:
        let pkgs = nixpkgsFor.${system};
        in {
          svgterm = pkgs.buildGoModule {
            pname = "svgterm";
            inherit version;

            src = ./.;

            # vendorHash = pkgs.lib.fakeHash;

            vendorHash = null;
          };
          default = self.packages.${system}.svgterm;
        });

      devShells = eachSystem (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];
          };
        });
    };
}
