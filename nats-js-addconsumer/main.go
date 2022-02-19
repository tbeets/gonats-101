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
	"encoding/json"
	"flag"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

func usage() {
	log.Printf("Usage: nats-js-addconsumer [-s server] [-creds file] [-nkey file] [-tlscert file] [-tlskey file] [-tlscacert file] <streamname> <consumername> <subfilter>\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var nkeyFile = flag.String("nkey", "", "NKey Seed File")
	var tlsClientCert = flag.String("tlscert", "", "TLS client certificate file")
	var tlsClientKey = flag.String("tlskey", "", "Private key file for client certificate")
	var tlsCACert = flag.String("tlscacert", "", "CA certificate to verify peer against")
	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) != 3 {
		showUsageAndExit(1)
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS JetStream Sample Publisher")}

	if *userCreds != "" && *nkeyFile != "" {
		log.Fatal("specify -seed or -creds")
	}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Use TLS client authentication
	if *tlsClientCert != "" && *tlsClientKey != "" {
		opts = append(opts, nats.ClientCert(*tlsClientCert, *tlsClientKey))
	}

	// Use specific CA certificate
	if *tlsCACert != "" {
		opts = append(opts, nats.RootCAs(*tlsCACert))
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

	// Create JetStream Context from NATS connection
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	str, con, subFilter := args[0], args[1], args[2]

	// Create a JS Consumer
	//
	// Durable consumers have Durable name... omitting name means Ephemeral
	// Pull consumers have dynamic DeliverSubject (DDS) so OMIT static DeliverSubject (SDS) in configuration
	//
	// Consumers inherit the persistence store of their associated stream (i.e. whether delivery state recorded in
	// memory or filestore).
	//
	conInfo, err := js.AddConsumer(str, &nats.ConsumerConfig{
		Durable:       con,
		FilterSubject: subFilter,
		AckPolicy:     nats.AckExplicitPolicy,
	})

	if err != nil {
		log.Fatal(err)
	}

	/* JS Consumer configuration options
	{
	  "durable_name": "thebarone",
	  "deliver_policy": "all",
	  "ack_policy": "explicit",
	  "ack_wait": 30000000000,
	  "max_deliver": -1,
	  "filter_subject": "foo8.bar",
	  "replay_policy": "instant",
	  "max_waiting": 512,
	  "max_ack_pending": 20000
	}
	*/

	// Pretty print our happy result to show all defaults etc.
	if conInfo != nil {
		conCfg := conInfo.Config
		jsonCfg, err := json.MarshalIndent(&conCfg, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s", string(jsonCfg))

	}
}
