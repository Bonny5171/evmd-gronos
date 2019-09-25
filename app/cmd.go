package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

// Declaring flags
var (
	Trace        bool
	Version      bool
	ShortVersion bool
)

func init() {
	flag.BoolVarP(&Trace, "trace", "t", true, "show trace log")
	flag.BoolVarP(&Version, "version", "v", false, "show version")
	flag.BoolVarP(&ShortVersion, "VERSION", "V", false, "show tiny version")
	flag.CommandLine.MarkHidden("VERSION")
}

// Arg returns the i'th command-line argument. Arg(0) is the first remaining argument
// after flags have been processed. Arg returns an empty string if the
// requested element does not exist.
func Arg(i int) string {
	return flag.Arg(i)
}

// Args returns the non-flag command-line arguments.
func Args() []string {
	return flag.Args()
}

// StartFlags initialize flags arguments to the app.
func StartFlags() {
	flag.Usage = showUsageFlags
	flag.Parse()

	if ShortVersion {
		fmt.Print(VERSION)
		os.Exit(0)
	}

	if Version {
		fmt.Printf("Version: %s\n", VERSION)
		os.Exit(0)
	}
}

func showUsageFlags() {
	fmt.Fprintf(os.Stdout, "Gronos - Everymind 2019©\n\n")
	fmt.Fprintf(os.Stdout, "Usage: %s [optional flags]\n\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "Optional Flags:\n\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stdout, "\n")
}
