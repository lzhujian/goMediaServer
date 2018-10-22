package app

import (
	"github.com/lzhujian/goMediaServer/manager"
	"github.com/lzhujian/goMediaServer/server"
	"log"
	"net/http"
)

/*
	live/{name} http处理函数, 新建流channel或者从channel拉流
*/
func liveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("liveHandler, req=", r)

	c, err := manager.GetManager().GetChannel(r.RequestURI)
	if err != nil {
		log.Println("GetChannel failed, uri:", r.RequestURI)
		return
	}

	if r.Method == "PUT" || r.Method == "POST" {
		err = c.HandlePublish(w, r)
		if err != nil {
			manager.GetManager().DeleteChannel(r.RequestURI)
			log.Printf("handle publish failed, err=%v, delete channel from manager.", err)
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
