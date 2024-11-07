# on

Run command on file event. Like `watch` but for file changes. [↗️Related blog article](https://evilcookie.de/on-run-commands-on-file-event.html).

## Getting started

usage:
```bash
usage: on [--create] [--write] [--rename] [--remove] [--chmod] <file> <cmd...>
```

Listen to all events for a file and print date for each event:
```bash
$ on ./main.go date
Tue Nov  5 12:27:53 PM CET 2024
Tue Nov  5 12:27:53 PM CET 2024
```

Run git diff for any write event:
```bash
$ on --write ./main.go git diff
diff --git a/main.go b/main.go
index e400b66..baa2c16 100644
--- a/main.go
+++ b/main.go
@@ -17,8 +17,7 @@ var (
    remove = flag.Bool("remove", false, "React on remove")
    chmod  = flag.Bool("chmod", false, "React on chmod")

-   verbose = flag.Bool("verbose", false, "Print debug information")
-
+   verbose     = flag.Bool("verbose", false, "Print debug information")
    listenToAll = false
 )
```

Debug output with all events:
```bash
$ on --verbose ./main.go true
ops: []fsnotify.Op{}
received REMOVE        "./main.go"
received CREATE        "./main.go"
received WRITE         "./main.go"
received CHMOD         "./main.go"
received RENAME        "./main.go"
received CREATE        "./main.go"
received WRITE         "./main.go"
received WRITE         "./main.go"
received CHMOD         "./main.go"
```

## Limitations

- given file needs to have a path, at least `./`
- ~~events are not debounced, typically writes appear in masses~~ added with https://github.com/RaphaelPour/on/pull/1
