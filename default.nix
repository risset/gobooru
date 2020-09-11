with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "gobooru";

  buildInputs = with pkgs; [
    go
    gopkgs
    delve
    gopls
    go-tools
  ];
}
