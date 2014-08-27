  $ audibly ./does-not-exist
  audibly: starting command: fork/exec ./does-not-exist: no such file or directory
  [1]

  $ printf '#!/bin/sh\necho ESPEAK BORK\nexit 42' >espeak
  $ chmod a+x espeak
  $ PATH="$PWD:$PATH"
  $ audibly echo hello, world
  hello, world
  ESPEAK BORK
  audibly: espeak: exit status 42
  [1]
