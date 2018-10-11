// sProxy project main.go
package main

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/myml/go-socks5"
)

func main() {
	server()
}

type Proxy struct{}

func (p *Proxy) New(args *struct{}, reply *int) error {

	l, err := net.Listen("tcp", ":")
	if err != nil {
		return err
	}
	addr := l.Addr().String()
	port := addr[strings.LastIndex(addr, ":")+1:]
	*reply, err = strconv.Atoi(port)
	if err != nil {
		return err
	}

	go func() {
		defer l.Close()
		t := time.AfterFunc(time.Second*3, func() {
			l.Close()
		})
		c, err := l.Accept()
		if err != nil {
			return
		}
		t.Stop()
		log.Println("accept", c.LocalAddr())
		s, err := socks5.New(&socks5.Config{})
		err = s.ServeConn(c)
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}
func proxy(resp http.ResponseWriter, req *http.Request) {
	l, err := net.Listen("tcp", ":")
	if err != nil {
		log.Println(err)
		return
	}
	addr := l.Addr().String()
	port := addr[strings.LastIndex(addr, ":")+1:]
	resp.Header().Set("port", port)

	go func() {
		defer l.Close()
		t := time.AfterFunc(time.Second*3, func() {
			l.Close()
		})
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		t.Stop()
		log.Println("accept", c.LocalAddr())
		defer time.AfterFunc(time.Hour, func() {
			c.Close()
		}).Stop()
		s, err := socks5.New(&socks5.Config{})
		err = s.ServeConn(c)
		if err != nil {
			log.Println(err)
		}
	}()
}
func server() {
	http.HandleFunc("/", proxy)
	http.ListenAndServe(":4433", nil)
}
