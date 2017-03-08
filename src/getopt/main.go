package main

import (
	"fmt"
	"os"
	"strconv"

	"code.google.com/p/getopt"
)

func main() {

	var s *getopt.Set = getopt.New() // 新建指针

	var node string
	var help bool = false
	var version bool = false

	s.StringVarLong(&node, "node", 'N', "The node tp startup", "NODE")

	s.BoolVarLong(&help, "help", 'h', "Print this help")
	s.BoolVarLong(&version, "version", 'v', "Print version information")

	s.Parse(os.Args) // 参数为命令行参数

	if help {

		s.PrintUsage(os.Stderr)

		return
	}

	if version {

		os.Stderr.WriteString("lorisd 1.5\n")

		os.Stderr.WriteString("Build: $Revision$, built on $Date$\n")

		return
	}

	if node != "" {

		if _, e := strconv.ParseUint(node, 16, 32); e != nil {

			os.Stderr.WriteString(fmt.Sprintf("loris: unrecognized service 0x%v,%v", node, e))

			os.Exit(1)

		} else {

		}
	}
}
