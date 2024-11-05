package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var (
	create = flag.Bool("create", false, "React on create")
	write  = flag.Bool("write", false, "React on write")
	rename = flag.Bool("rename", false, "React on rename")
	remove = flag.Bool("remove", false, "React on remove")
	chmod  = flag.Bool("chmod", false, "React on chmod")

	verbose = flag.Bool("verbose", false, "Print debug information")

	listenToAll = false
)

func ArrayHas[K comparable](arr []K, s K) bool {
	Debugf("arr=%#v s=%v\n", arr, s)
	for _, element := range arr {
		if element == s {
			return true
		}
	}

	return false
}

func Debugf(format string, args ...any) {
	if !*verbose {
		return
	}

	fmt.Printf(format, args...)
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("usage: on [--create] [--write] [--rename] [--remove] [--chmod] <file> <cmd...>")
		return
	}

	Debugf("os.Args: %#v\n", os.Args)
	Debugf("flag.Args: %#v\n", flag.Args())

	ops := make([]fsnotify.Op, 0)
	if *create {
		ops = append(ops, fsnotify.Create)
	}
	if *write {
		ops = append(ops, fsnotify.Write)
	}
	if *rename {
		ops = append(ops, fsnotify.Rename)
	}
	if *remove {
		ops = append(ops, fsnotify.Remove)
	}
	if *chmod {
		ops = append(ops, fsnotify.Chmod)
	}

	Debugf("ops: %#v\n", ops)

	// if no operation is given, listen to all
	if len(ops) == 0 {
		listenToAll = true
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("error initializing watcher: ", err)
		return
	}

	err = watcher.Add(filepath.Dir(flag.Args()[0]))
	if err != nil {
		fmt.Println("error setting up watcher: ", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Name != flag.Args()[0] {
				continue
			}

			Debugf("received %v\n", event)
			if !listenToAll && !ArrayHas(ops, event.Op) {
				continue
			}

			cmd := exec.Command(flag.Args()[1], flag.Args()[2:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			_ = cmd.Run()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("error:", err)
		}
	}
}
