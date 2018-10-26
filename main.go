package main

import (
	"flag"
	"github.com/lzhujian/goMediaServer/app"
	"github.com/lzhujian/goMediaServer/server"
	"log"
	"os"
)

var (
	version = "master"
	httpflvAddr = flag.String("httpflv-addr", ":8001", "HTTP-FLV server listen address")
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	log.Println("Start media server, version:", version)

	// 启动 http-flv server
	httpServer := server.NewHttpServer(*httpflvAddr)
	app.AddHandler(httpServer)
	httpServer.Start()
}
