package main

import (
	"a-dns/cmd"
	"a-dns/internal/i18n"
	"os"
)

func main() {
	if err := i18n.Init("fr"); err != nil {
		os.Exit(1)
	}
	cmd.Execute()
}
