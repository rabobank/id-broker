package main

import (
	"fmt"

	"github.com/rabobank/id-broker/cf"
	"github.com/rabobank/id-broker/cfg"
	"github.com/rabobank/id-broker/managers"
	"github.com/rabobank/id-broker/server"
	"github.com/rabobank/id-broker/uaa"
)

func main() {
	fmt.Printf("id-broker starting, Version:%s, CommitHash:%s, buildTime:%s\n", cfg.Version, cfg.CommitHash, cfg.BuildTime)

	cfg.Initialize()
	cf.Initialize()
	uaa.Initialize()
	managers.Initialize()

	server.StartServer()
}
