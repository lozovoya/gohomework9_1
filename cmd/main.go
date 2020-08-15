package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

//type Listener interface {
//	Accept() (Conn, error)
//	Close() error
//	Addr() Addr
//}
//
//type Conn interface {
//	Read (b []byte)
//}

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
		handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(cerr)
			return
		}
	}()

	reader := bufio.NewReader(conn)
	const delim = '\n'
	for {
		line, err := reader.ReadString(delim)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return
			}
			log.Println("received: %s\n", line)
		}
		log.Println("received: %s\n", line)
	}
}
