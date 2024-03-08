package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/onyx-and-iris/q3rcon-proxy/pkg/udpproxy"
)

func start(proxy string) {
	port, target := func() (string, string) {
		x := strings.Split(proxy, ":")
		return x[0], x[1]
	}()

	c, err := udpproxy.New(fmt.Sprintf("%s:%s", host, port), fmt.Sprintf("127.0.0.1:%s", target))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("q3rcon-proxy initialized: [proxy] (%s:%s) [target] (127.0.0.1:%s)", host, port, target)

	log.Fatal(c.ListenAndServe())
}

var (
	proxies, host string
)

func init() {
	proxies = os.Getenv("Q3RCON_PROXY")
	if proxies == "" {
		log.Fatal("env Q3RCON_PROXY required")
	}

	host = os.Getenv("Q3RCON_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
}

func main() {
	for _, proxy := range strings.Split(proxies, ";") {
		go start(proxy)
	}

	<-make(chan int)
}
