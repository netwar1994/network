package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/netwar1994/network/pkg/card"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main(){
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

func execute() (err error) {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Println(cerr)
			if err == nil {
				err = cerr
			}
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	r := bufio.NewReader(conn)
	const delim = '\n'
	line, err := r.ReadString(delim)
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
		log.Printf("invalid request line: %s", line)
		return
	}

	path := parts[1]

	switch path {
	case "/":
		err = writeIndex(conn)
	case "/operations.csv":
		err = writeOperationsCSV(conn)
	case "/operations.json":
		err = writeOperationsJSON(conn)
	case "/operations.xml":
		err = writeOperationsXML(conn)
	default:
		err = write404(conn)
	}
	if err != nil {
		log.Println(err)
		return
	}
}

func writeResponse(w io.Writer, status int, headers []string, content []byte) error {
	const CRLF = "\r\n"
	var err error

	writer := bufio.NewWriter(w)
	_, err = writer.WriteString(fmt.Sprintf("HTTP/1.1 %d OK%s", status, CRLF))
	if err != nil {
		return err
	}

	for _, h := range headers {
		_, err = writer.WriteString(h + CRLF)
		if err != nil {
			return err
		}
	}

	_, err = writer.WriteString(CRLF)
	if err != nil {
		return err
	}
	_, err = writer.Write(content)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func writeIndex(w io.Writer) error {
	username := "Василий"
	balance := "1 000.50"

	page, err := ioutil.ReadFile("web/template/index.html")
	if err != nil {
		return err
	}
	page = bytes.ReplaceAll(page, []byte("{username}"), []byte(username))
	page = bytes.ReplaceAll(page, []byte("{balance}"), []byte(balance))

	return writeResponse(w, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

func writeOperationsCSV(w io.Writer) error {
	transactions := card.MakeTransactions(1)
	page, err := card.ExportCSV(transactions)
	if err != nil {
		log.Println(err)
	}

	return writeResponse(w, 200, []string{
		"Content-Type: text/csv",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

func writeOperationsJSON(w io.Writer) error {
	transactions := card.MakeTransactions(1)
	page, err := card.ExportJson(transactions)
	if err != nil {
		log.Println(err)
	}

	return writeResponse(w, 200, []string{
		"Content-Type: application/json",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

func writeOperationsXML(w io.Writer) error {
	transactions := card.MakeTransactions(1)
	page, err := card.ExportXML(transactions)
	if err != nil {
		log.Println(err)
	}

	return writeResponse(w, 200, []string{
		"Content-Type: application/xml",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

func write404(writer io.Writer) error {
	page, err := ioutil.ReadFile("web/template/error404.html")
	if err != nil {
		return err
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}