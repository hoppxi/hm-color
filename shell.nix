{
  pkgs ? import <nixpkgs> { },
  system ? builtins.currentSystem,
}:

let
  mcuc = import (pkgs.fetchFromGitHub {
    owner = "hoppxi";
    repo = "mcu-cli";
    rev = "main";
    sha256 = "sha256-LOroCdgUCyZN/S1HON6PAUhQXCDx8llSZ+xhZjsoLPs=";
  }) { inherit pkgs; };
in
pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    swww
    mcuc
  ];
}
