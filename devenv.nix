{ pkgs, ... }:

{
  packages = [
    pkgs.jo
  ];

  languages.go.enable = true;
}
