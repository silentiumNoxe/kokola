package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var args, err = parseArgs()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(args.Target)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	var line int64
	var buff = make([]byte, 0x10)
	for {
		n, err := file.ReadAt(buff, line)

		if n > 0 {
			print(strings.ToUpper(fmt.Sprintf("[%08x]", line)))
			for i, x := range buff {
				print(strings.ToUpper(fmt.Sprintf(" %02x", x)))
				buff[i] = 0
			}

			println()

			line += 0x10
		}

		if line > 0xF0 && err == nil {
			var r = bufio.NewReaderSize(os.Stdin, 1)
			print("Press ENTER to continue...")
			b, _ := r.ReadByte()

			// Clear previous line
			// \033[1A - one line up
			// \033[K - delete the line
			print("\033[1A\033[K")

			if b == 0x03 {
				return
			}
		}

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}

func parseArgs() (Args, error) {
	var help = flag.Bool("h", false, "Show help")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	var args = flag.Args()
	if len(args) == 0 {
		return Args{}, errors.New("no target file specified")
	}

	return Args{Target: args[0]}, nil
}

type Args struct {
	Target string
}
