package main

import (
	"fmt"
	"github.com/reiver/go-telnet"
	"net"
)

func main() {
	RunTelnetServer()
}

var (
	Port = 5555
	//Listener net.Listener
)

type Server telnet.Server

func handler(port int) telnet.Handler {
	return telnet.EchoHandler
}

func RunTelnetServer() error {
	server := &Server{Addr: fmt.Sprintf(":%v", Port), Handler: handler(Port)}
	err := server.MyListenAndServe()

	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
	return err
}

func StopTelnetServer(listener net.Listener) {
	fmt.Println("[Telnet] Shutting down ...")
	listener.Close()
}

func (server *Server) MyListenAndServe() error {

	addr := server.Addr
	if "" == addr {
		addr = ":telnet"
	}

	listener, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}

	//Listener = listener
	fmt.Println("[Telnet] Listening at port: ", Port)
	return telnet.Serve(listener, handler(Port))
}
