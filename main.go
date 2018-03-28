package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hpcloud/tail"
	stan "github.com/nats-io/go-nats-streaming"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func connectionCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %s", err)
	}
}

func main() {
	clusterName := getEnv("CLUSTER_NAME", "test-cluster")
	natsServer := getEnv("NATS_SERVER", "nats://localhost:4222")
	natsClientName := getEnv("NATS_CLIENT_NAME", "tailnats")
	tailFile := getEnv("TAIL_FILE", "/var/log/test.log")
	if len(os.Args) > 1 {
		tailFile = os.Args[1]
	}
	natsSubject := strings.TrimSuffix(filepath.Base(tailFile), filepath.Ext(tailFile))

	conn, err := stan.Connect(
		clusterName,
		natsClientName,
		stan.NatsURL(natsServer),
	)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer connectionCloser(conn)

	err = conn.Publish("natslog.subscribe", []byte(natsSubject))
	if err != nil {
		log.Fatalf("error subscribing: %v", err)
	}

	t, err := tail.TailFile(tailFile, tail.Config{Follow: true})
	for line := range t.Lines {
		perr := conn.Publish(natsSubject, []byte(fmt.Sprintf("%s\n", line.Text)))
		if perr != nil {
			log.Printf("Publish error %v", perr)
		}
	}
}
