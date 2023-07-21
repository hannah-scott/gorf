package main

import (
	"bufio"
	"os"
	"strings"
	"io"
)


func main () {
	reader := bufio.NewReader(os.Stdin)
	s := Machine{
		top: -1,
		err: "",
	}
	s.InitializeDict()

	for {
		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			os.Exit(1)
		}
		if len(strings.Fields(input)) > 0 {
			s.rstack = nil
			s.addr = 0
			s.Execute(strings.ToUpper(input))
			if s.err != "" {
				print(s.err + "\n")
				s.err = ""
				os.Exit(1)
			} else {
				print("ok\n")
			}
		}
		if err == io.EOF { os.Exit(0) }
	}
}