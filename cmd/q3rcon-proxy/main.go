package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/onyx-and-iris/q3rcon-proxy/pkg/udpproxy"
)

func main() {
	logLevel, err := getEnvInt("Q3RCON_LOGLEVEL")
	if err != nil {
		log.Fatalf("unable to parse Q3RCON_LEVEL: %s", err.Error())
	}
	if slices.Contains(log.AllLevels, log.Level(logLevel)) {
		log.SetLevel(log.Level(logLevel))
	}

	proxies := os.Getenv("Q3RCON_PROXY")
	if proxies == "" {
		log.Fatal("env Q3RCON_PROXY required")
	}

	host := os.Getenv("Q3RCON_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	staleTimeout, err := getEnvInt("Q3RCON_STALE_SESSION_TIMEOUT")
	if err != nil {
		log.Fatalf("unable to parse Q3RCON_STALE_SESSION_TIMEOUT: %s", err.Error())
	}

	for _, proxy := range strings.Split(proxies, ";") {
		go start(host, proxy, staleTimeout)
	}

	<-make(chan int)
}

func start(host, proxy string, staleTimeout int) {
	port, target := func() (string, string) {
		x := strings.Split(proxy, ":")
		return x[0], x[1]
	}()

	c, err := udpproxy.New(
		fmt.Sprintf("%s:%s", host, port),
		fmt.Sprintf("127.0.0.1:%s", target),
		udpproxy.WithStaleTimeout(time.Duration(staleTimeout)*time.Minute))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("q3rcon-proxy initialized: [proxy] (%s:%s) [target] (127.0.0.1:%s)", host, port, target)

	log.Fatal(c.ListenAndServe())
}
