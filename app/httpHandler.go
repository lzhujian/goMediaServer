package app

import (
	"github.com/lzhujian/goMediaServer/manager"
	"github.com/lzhujian/goMediaServer/server"
	"log"
	"net/http"
)

func liveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("liveHandler, req=", r, "method=", r.Method)

	c, err := manager.GetManager().GetChannel(r.RequestURI)
	if err != nil {
		log.Println("GetChannel failed, uri:", r.RequestURI)
		return
	}

	if r.Method == "PUT" || r.Method == "POST" {
		err = c.HandlePublish(w, r)
		if err != nil {
			log.Println("handle publish failed, err:", err)
			return
		}
	} else if r.Method == "GET" {
		err = c.HandlePlay(w, r)
		if err != nil {
			log.Println("handle play failed, err:", err)
			return
		}
	}
}

func AddHandler(server *server.HttpServer) {
	server.HandlerMap["/live/{name}"] = liveHandler
}
