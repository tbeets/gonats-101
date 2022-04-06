// Copyright 2012-2021 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var countCh = make(chan struct{}, 128)

func usage() {
	log.Printf("Usage: nats-req-multi [-s server] [-creds file] [-nkey file] [-d {reply duration}] [-m {max replies}] <subject> <msg>\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func doReqWait(nc *nats.Conn, subj string, body []byte, dur int, max int) error {
	start := time.Now()

	msg := nats.Msg{
		Subject: subj,
		Reply:   nc.NewRespInbox(),
		Header:  nil,
		Data:    body,
		Sub:     nil,
	}

	s, err := nc.Subscribe(msg.Reply, func(m *nats.Msg) {
		countCh <- struct{}{}

		log.Printf("Received on %q rtt %v", m.Subject, time.Since(start))

		if len(m.Header) > 0 {
			for h, vals := range m.Header {
				for _, val := range vals {
					log.Printf("%s: %s", h, val)
				}
			}

			fmt.Println()
		}

		fmt.Println(string(m.Data))
		if !strings.HasSuffix(string(m.Data), "\n") {
			fmt.Println()
		}
	})
	if err != nil {
		return err
	}
	defer s.Unsubscribe()

	err = nc.PublishMsg(&msg)
	if err != nil {
		return err
	}

	received := 0
Loop:
	for {
		select {
		case <-countCh:
			received++
			if received >= max {
				break Loop
			}
		case <-time.After(time.Duration(dur) * time.Second):
			break Loop
		}
	}

	// we don't want any responses after we break.
	s.Unsubscribe()

	return nil
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var nkeyFile = flag.String("nkey", "", "NKey Seed File")
	var showHelp = flag.Bool("h", false, "Show help message")
	var duration = flag.Int("d", 2, "Reply interest duration (seconds)")
	var max = flag.Int("m", 1, "Maximum number of replies")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) < 2 {
		showUsageAndExit(1)
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Requestor")}

	if *userCreds != "" && *nkeyFile != "" {
		log.Fatal("specify -seed or -creds")
	}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Use Nkey authentication.
	if *nkeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(*nkeyFile)
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, opt)
	}

	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	subj, payload := args[0], []byte(args[1])

	log.Printf("Published [%s] : '%s'", subj, payload)

	// msg, err := nc.Request(subj, payload, 2*time.Second)
	doReqWait(nc, subj, payload, *duration, *max)
	if err != nil {
		if nc.LastError() != nil {
			log.Fatalf("%v for request", nc.LastError())
		}
		log.Fatalf("%v for request", err)
	}

}
