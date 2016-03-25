package cmd

import (
	"fmt"
	"os"
)


// fatal prints the message and then exits. If V(2) or greater, glog.Fatal
// is invoked for extended information.
func fatal(msg string) {
	fmt.Fprint(os.Stderr, msg)
	os.Exit(1)
}