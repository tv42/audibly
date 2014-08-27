  $ printf '#!/bin/sh\necho "SPEAK $@ END"' >espeak
  $ chmod a+x espeak
  $ PATH="$PWD:$PATH"
  $ audibly echo hello, world
  hello, world
  SPEAK -- success: echo END
  $ audibly -name=greeting -- echo hello, world
  hello, world
  SPEAK -- success: greeting END
  $ audibly -success=hurrah -- echo hello, world
  hello, world
  SPEAK -- hurrah END
  $ audibly -failure="oh noes" false
  SPEAK -- oh noes END
  [1]
  $ audibly -failure='status {{.Status}}' sh -c 'exit 42'
  SPEAK -- status 42 END
  [42]
