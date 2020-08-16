package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
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
	line, err := reader.ReadString(delim)
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		log.Printf("received: %s\n", line)
		return
	}
	log.Printf("received: %s\n", line)

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		log.Printf("Invlaid request line: #{line}")
		return
	}

	time.Sleep(time.Second * 2)
	path := parts[1]

	switch path {
	case "/":
		err = writeIndex(conn)
	case "/operations.csv":
		err = writeOperations(conn)
	default:
		err = write404(conn)
	}
	if err != nil {
		log.Println(err)
		return
	}
}

func writeIndex(writer io.Writer) error {
	username := "Васян"
	balance := "1000"

	page, err := ioutil.ReadFile("web/template/index.html")
	if err != nil {
		log.Println(err)
		return err
	}
	page = bytes.ReplaceAll(page, []byte("{username}"), []byte(username))
	page = bytes.ReplaceAll(page, []byte("{balance}"), []byte(balance))

	return writeResponse(writer, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

func writeOperations(writer io.Writer) error {
	page, err := ioutil.ReadFile("web/template/export.csv")
	if err != nil {
		log.Println(err)
		return err
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: text/csv",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

func write404(writer io.Writer) error {
	page, err := ioutil.ReadFile("web/template/404.html")
	if err != nil {
		return err
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

func writeResponse(
	writer io.Writer,
	status int,
	headers []string,
	content []byte,
) error {
	const CRLF = "\r\n"
	var err error

	w := bufio.NewWriter(writer)
	_, err = w.WriteString(fmt.Sprintf("HTTP/1.1 %d OK%s", status, CRLF))
	if err != nil {
		return err
	}

	for _, h := range headers {
		_, err = w.WriteString(h + CRLF)
		if err != nil {
			return err
		}
	}

	_, err = w.WriteString(CRLF)
	if err != nil {
		return err
	}
	_, err = w.Write(content)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}
