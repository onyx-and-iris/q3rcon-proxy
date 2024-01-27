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

	c, err := udpproxy.New(fmt.Sprintf("0.0.0.0:%s", port), fmt.Sprintf("127.0.0.1:%s", target))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("q3rcon-proxy initialized: [proxy] (0.0.0.0:%s) [target] (127.0.0.1:%s)", port, target)

	log.Fatal(c.ListenAndServe())
}

func main() {
	proxies := os.Getenv("Q3RCON_PROXY")
	if proxies == "" {
		log.Fatal("env Q3RCON_PROXY required")
	}

	for _, proxy := range strings.Split(proxies, ";") {
		go start(proxy)
	}

	<-make(chan int)
}
