package main

import (
	"context"
	"log"
	"net"

	"github.com/things-go/go-socks5"
	"golang.org/x/net/proxy"
)

func main() {
	torDialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		log.Fatalf("Failed to start Tor proxy dialer: %v", err)
	}
	server := socks5.NewServer(
		socks5.WithDial(func(ctx context.Context, network, addr string) (net.Conn, error) {
			if network == "tcp" {
				log.Printf("[TCP] Routing %s through Tor", addr)
				return torDialer.Dial(network, addr)
			}
			log.Printf("[UDP] Routing %s normally", addr)
			return net.Dial(network, addr)
		}),
	)
	log.Println("Middleware SOCKS5 proxy running on 127.0.0.1:1080...")
	if err := server.ListenAndServe("tcp", "127.0.0.1:1080"); err != nil {
		log.Fatal(err)
	}
}
