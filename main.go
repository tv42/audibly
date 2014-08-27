package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"
)

var (
	name    = flag.String("name", "", "how to identify this run (default: basename of command)")
	success = flag.String("success", "success: {{.Name}}", "say this message on success")
	failure = flag.String("failure", "failure: {{.Name}}", "say this message on failure")
)

var templates = template.New("speech")

type Info struct {
	name string
	Cmd  *exec.Cmd
}

func (i *Info) Name() string {
	if *name != "" {
		return *name
	}
	return filepath.Base(i.Cmd.Args[0])
}

func (i *Info) Status() int {
	if i.Cmd.ProcessState.Success() {
		return 0
	}
	if status, ok := i.Cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		if status.Exited() {
			return status.ExitStatus()
		}
		if status.Signaled() {
			return 128 + int(status.Signal())
		}
	}
	// uhh.. something. 1 is too common, and 2 is typical usage error.
	return 3
}

func run(args []string) (*Info, error) {
	i := &Info{
		Cmd: exec.Command(args[0], args[1:]...),
	}
	i.Cmd.Stdin = os.Stdin
	i.Cmd.Stdout = os.Stdout
	i.Cmd.Stderr = os.Stderr

	if err := i.Cmd.Run(); err != nil {
		switch err.(type) {
		case *exec.ExitError:
			// we don't consider runs that actually start the child as failures
		default:
			return i, err
		}
	}
	return i, nil
}

func report(i *Info) error {
	msg, err := format(i)
	if err != nil {
		return err
	}
	if err := speak(msg); err != nil {
		return err
	}
	return nil
}

func format(i *Info) (string, error) {
	tmpl := templates.Lookup("success")
	if !i.Cmd.ProcessState.Success() {
		tmpl = templates.Lookup("failure")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, i); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func speak(msg string) error {
	cmd := exec.Command("espeak", "--", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", prog)
	fmt.Fprintf(os.Stderr, "  %s [OPTS] [--] COMMAND [ARGS..]\n", prog)
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// exit 2 here as we consider these usage errors
	if _, err := templates.New("success").Parse(*success); err != nil {
		log.Print(err)
		os.Exit(2)
	}
	if _, err := templates.New("failure").Parse(*failure); err != nil {
		log.Print(err)
		os.Exit(2)
	}

	i, err := run(flag.Args())
	if err != nil {
		log.Fatalf("starting command: %v", err)
	}
	if err := report(i); err != nil {
		log.Fatalf("espeak: %v", err)
	}
	os.Exit(i.Status())
}
