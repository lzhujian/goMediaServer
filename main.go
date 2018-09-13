package main

import (
	"github.com/lzhujian/goMediaServer/app"
	"github.com/lzhujian/goMediaServer/server"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	log.Println("Start media server")

	httpServer := server.NewHttpServer(8888)
	app.AddHandler(httpServer)
	httpServer.Start()
}
