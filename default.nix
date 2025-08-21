(import
  (fetchTarball {
    url = "https://github.com/edolstra/flake-compat/archive/master.tar.gz";
    sha256 = "0b0cfjqzbgc1cif2n0kv8i8xqbq7p1psf49wirpcsgg9g8c8iiyr";
  })
  {
    src = ./.;
  }
).defaultNix
