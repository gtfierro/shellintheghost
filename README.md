# Shell in the Ghost

This is a package to implement remote shells over BOSSWAVE.

The `conn` subpackage implements Go's `net.Conn` interface over BOSSWAVE, so that can be used for other projects.


## Install

```bash
$ go get github.com/gtfierro/shellintheghost
$ go install github.com/gtfierro/shellintheghost
```

## Usage

Requires a server and a client.

### Server

You need a base URI of the service and a set of terminals names. Each terminal corresponds to a single instance
of the shell, so if you need to support multiple concurrent users, then it is a good idea to create separate terminals.

```
shellintheghost server -u scratch.ns/terminals -t gabe
```

This will create a terminal at `scratch.ns/terminals/s.shell/_/i.term/slot/gabe` with corresponding output at 
`scratch.ns/terminals/s.shell/_/i.term/signal/gabe`. Additional terminals can be created with more instances of the `-t`  flag.

Servers can expose whatever shell they want. The default is `/bin/bash`, but using the `-s` flag you can specify 
`/usr/bin/python` or `which supervisorctl` or whatever you want.

### Client

```
shellintheghost client -u scratch.ns/terminals -t gabe
```

Will connect to the server -- the client handles expanding the URI out w/ the service/interface names.

#### Gotchas

You should press `enter` at least once when starting so that you can see your prompt, else the server might publish
it and it won't show up on your screen. The shell is still accepting input though.

The current version does not account for the size of your client's terminal window. Most things work, but this makes
using `vim` and `byobu` and similar tools difficult. You can use `stty` to inform the terminal of the size you want

```bash
# for a perfectly square terminal window
# type these in after logging in
$ stty rows 100
$ stty cols 100
$ # verify the size
$ stty size
100 100
```

Exit with `Ctl-D` or `exit` or whatever you usually use.
