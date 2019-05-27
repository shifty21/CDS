with import <nixpkgs> {};
stdenv.mkDerivation rec
{
  name = "himeno";
  src = ./. ;

  buildInputs = [ openssl ] ;

  buildPhase = ''
  make clean
  make
  '';

  installPhase = ''
  mkdir -p $out/bin
  cp himeno $out/bin
  '';
}
