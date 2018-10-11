// sProxy project main.go
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/myml/go-socks5"
	"golang.org/x/net/proxy"
)

func main() {
	uri, err := url.Parse("http://t4.flys.cf:4433")
	if err != nil {
		log.Panic(err)
	}
	client(uri)
}

func client(uri *url.URL) {
	conf := socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			resp, err := http.Get(uri.String())
			if err != nil {
				return nil, err
			}
			port := resp.Header.Get("port")
			if port == "" {
				return nil, errors.New("can not get port")
			}
			dialer, err := proxy.SOCKS5("tcp", uri.Hostname()+":"+port, nil, nil)
			if err != nil {
				return nil, err
			}
			log.Println(network, addr, port)
			return dialer.Dial(network, addr)
		}}
	h, err := socks5.New(&conf)
	if err != nil {
		log.Panic(err)
	}
	h.ListenAndServe("tcp", ":7070")
}
