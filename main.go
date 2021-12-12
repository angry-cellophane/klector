package main

import (
	"io.klector/klector/commands"
	"os"
)

func main() {
	os.Exit(commands.Execute())
}
