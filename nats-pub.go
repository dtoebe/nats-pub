package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats"
)

func useage() {
	fmt.Println("Nats-io Message publisher:\n\tCMD: nats-pub <options: [-s server-urls][-t timeout-int]> <subject> <message>")
}

func main() {
	urls := flag.String("s", nats.DefaultURL, fmt.Sprintf("Nats server URLs seperated by commas (Default: %s)", nats.DefaultURL))
	timeOut := flag.Int("t", 0, "If timeout is set (>0) in milliseconds then we send an RPC request with a response")

	flag.Usage = useage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		useage()
		return
	}

	nc, err := nats.Connect(*urls)
	if err != nil {
		fmt.Println("Unable to connecrt to nats server:", err.Error())
		os.Exit(1)
	}
	defer nc.Close()

	sub, msg := args[0], []byte(args[1])

	switch *timeOut {
	case 0:
		nc.Publish(sub, msg)
	default:
		m, err := nc.Request(sub, msg, time.Millisecond*time.Duration(*timeOut))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("\nRESPONSE:\n\nSubject: %s\nMessage: %s", m.Subject, string(m.Data))
	}

}
