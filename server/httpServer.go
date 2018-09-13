package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"time"
)

/*
	HttpServer struct
	Port: http服务器监听端口
	HandlerMap: handler name 和 handler 处理函数 map
*/
type HttpServer struct {
	Port       uint16
	HandlerMap map[string]func(w http.ResponseWriter, r *http.Request)
}

// New http server
func NewHttpServer(port uint16) *HttpServer {
	return &HttpServer{
		Port:       port,
		HandlerMap: make(map[string]func(w http.ResponseWriter, r *http.Request)),
	}
}

func ConnState(c net.Conn, cs http.ConnState) {
	idleTime := time.Second * 10

	switch cs {
	case http.StateIdle, http.StateNew:
		c.SetReadDeadline(time.Now().Add(idleTime))
	case http.StateActive:
		c.SetReadDeadline(time.Time{})
	}
}

// 设置http server router，启动http server
func (s *HttpServer) Start() {
	muxHandler := mux.NewRouter()

	for path, handler := range s.HandlerMap {
		muxHandler.HandleFunc(path, handler)
	}

	addr := fmt.Sprintf(":%d", s.Port)
	server := http.Server{
		Addr:        addr,
		Handler:     muxHandler,
		ReadTimeout: 0,
		ConnState:   ConnState,
	}

	//go func() {
	err := server.ListenAndServe()

	if err != nil {
		log.Println("server listen and serve failed", err)
	}
	//}()
}
