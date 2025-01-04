package main

import (
	"os"

	"github.com/KentaroAOKI/mssqlbeat/cmd"

	_ "github.com/KentaroAOKI/mssqlbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
