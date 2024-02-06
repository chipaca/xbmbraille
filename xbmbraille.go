package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func usage() {
	cmd := flag.CommandLine
	fmt.Fprintf(cmd.Output(), "Usage: %s [options] {-|filename}...\n", os.Args[0])
	cmd.PrintDefaults()
}

func main() {
	var neg bool
	var printName bool
	var clear bool
	var delay time.Duration

	flag.Usage = usage
	flag.BoolVar(&neg, "n", false, "negate (invert) image")
	flag.BoolVar(&clear, "c", false, "clear the terminal before printing each image")
	flag.BoolVar(&printName, "p", false, "print the filename before printing each image")
	flag.DurationVar(&delay, "d", 0, "wait this much time after printing each image")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}
	for _, arg := range args {
		var r io.Reader
		var f *os.File
		if arg == "-" {
			r = os.Stdin
		} else {
			var err error
			f, err = os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}
			r = f
		}
		xbm, err := fromReader(r, neg)
		if err != nil {
			log.Fatalf("when working on %q: %v", arg, err)
		}
		if f != nil {
			f.Close()
		}
		if clear {
			os.Stdout.WriteString("\033[H\033[2J\033[3J")
		}
		if printName {
			os.Stdout.WriteString(arg + "\n")
		}
		os.Stdout.WriteString(xbm.braille())
		if delay > 0 {
			time.Sleep(delay)
		}
	}
}
