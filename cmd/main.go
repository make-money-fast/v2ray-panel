package main

import (
	"flag"
	"github.com/make-money-fast/v2ray/helpers"
	"github.com/make-money-fast/v2ray/server"
)

func init() {
	flag.BoolVar(&helpers.IsServer, "server", false, "is server")
	flag.BoolVar(&helpers.IsClient, "client", false, "is client")

	flag.StringVar(&helpers.HttpPort, "http-port", ":8091", "HTTP服务端口")
}

func main() {
	flag.Parse()

	helpers.IsServer = true

	server.StartServer()
}
