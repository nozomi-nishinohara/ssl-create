package main

import (
	"os"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createSSL() {
	if !exists("ssl") {
		os.Mkdir("ssl", 0700)
	}
	if !(exists("ssl/server.crt") && exists("ssl/server.key")) {
		tpl, pkey := ca()
		serverSSL(tpl, pkey)
	}
}

type GenCommand struct{}

func (c *GenCommand) Help() string {
	return ""
}

func (c *GenCommand) Synopsis() string {
	return ""
}

func (c *GenCommand) Run(args []string) int {
	createSSL()
	return 0
}
