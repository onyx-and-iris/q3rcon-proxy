package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

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
	proxies, host, debug string
)

func getenvInt(key string) (int, error) {
	s := os.Getenv(key)
	if s == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func init() {
	proxies = os.Getenv("Q3RCON_PROXY")
	if proxies == "" {
		log.Fatal("env Q3RCON_PROXY required")
	}

	host = os.Getenv("Q3RCON_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	debug, err := getenvInt("Q3RCON_DEBUG")
	if err != nil {
		log.Fatal(err)
	}

	if debug == 1 {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

}

func main() {
	for _, proxy := range strings.Split(proxies, ";") {
		go start(proxy)
	}

	<-make(chan int)
}
