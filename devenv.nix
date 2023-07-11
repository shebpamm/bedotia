{ pkgs, ... }:

{
  languages.python.enable = true;
  languages.python.poetry.enable = true;
  languages.python.poetry.activate.enable = true;
}
