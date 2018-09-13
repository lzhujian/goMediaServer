package app

import (
	"fmt"
	"github.com/lzhujian/goMediaServer/server"
	"log"
	"net/http"
)

func liveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("liveHandler, req=", r)
	fmt.Fprintf(w, "hello there, %s!", r.URL.Path[1:])
}

func AddHandler(server *server.HttpServer) {
	server.HandlerMap["/live/{name}"] = liveHandler
}
