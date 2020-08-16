package main

import (
	"github.com/lozovoya/gohomework9_1/pkg/webserver"
	"log"
	"net"
	"os"
)

func main() {

	if err := execute(); err != nil {
		os.Exit(1)
	}

}

func execute() (err error) {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Println(err)
			return
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go webserver.Handle(conn)
	}
}
