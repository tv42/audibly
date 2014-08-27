# `audibly` -- Audibly report command status

Run it like `audibly COMMAND [ARGS..]`, and `audibly` will report the
success or error of `COMMAND` with two different sounds.

It currently uses [eSpeak](http://espeak.sourceforge.net/) as the
speech synthesizer, and requires it to be installed. This may change
later.

``` console
# this will say "success: echo"
$ audibly echo hello, world

# this will say "success: greeting"
$ audibly -name=greeting -- echo hello, world

# this will say "hurrah"
$ audibly -success=hurrah -- echo hello, world

# this will say "oh noes" and exit with status 1
$ audibly -failure="oh noes" false

# this will say "status 42" and exit with status 42
$ audibly -failure='status {{.Status}}' sh -c 'exit 42'
```

You can also use templates to generate the message spoken.
The following fields are defined:

- `Name`: value of the `-name=` flag, or the basename of the command
- `Status`: exit status of the command
- `Cmd`: the [os/exec Cmd value](http://golang.org/pkg/os/exec/#Cmd)

The template mechanism is documented at
http://golang.org/pkg/text/template/
