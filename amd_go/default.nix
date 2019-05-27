with import <nixpkgs> {};
stdenv.mkDerivation rec
{
  name = "average_min_distance";
  src = ./. ;

  buildInputs = [ openssl ] ;

  buildPhase = ''
  make clean
  make
  '';

  installPhase = ''
  mkdir -p $out/bin
  cp amd $out/bin
  '';
}
