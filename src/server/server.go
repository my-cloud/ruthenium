package server

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	port uint16
}

func NewServer(port uint16) *Server {
	return &Server{port}
}

func (server *Server) Port() uint16 {
	return server.port
}

func HelloWorld(writer http.ResponseWriter, _ *http.Request) {
	io.WriteString(writer, "hello world")
}

func (server *Server) Run() {
	http.HandleFunc("/", HelloWorld)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(server.port)), nil))
}
