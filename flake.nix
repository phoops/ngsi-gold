{
  description = "Go client for NGSI-LD protocol";

  inputs.nixpkgs.url = "nixpkgs/release-22.05";

  outputs = { self, nixpkgs }: with nixpkgs.legacyPackages.x86_64-linux; {
    devShell.x86_64-linux = mkShell {
      hardeningDisable = [ "fortify" ];
      buildInputs = [
        go_1_18
        gotools # provides godoc
        go-task
      ];
    };
  };
}
