package main

import (
	"fmt"
	"os"

	"github.com/julienbreux/rabdis/internal/rabdis/command"
)

func main() {
	cmd := command.NewCmdRoot(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
