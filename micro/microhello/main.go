package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"log"
	"os"
	"os/signal"
	"time"
)

func AddHelloService(nc *nats.Conn) (micro.Service, error) {
	helloHandler := func(req *micro.Request) error {
		req.Respond([]byte(fmt.Sprintf("A hearty micro Hello to ya' [%s]", time.Now().String())))
		return nil
	}

	config := micro.Config{
		Name:        "MicroHelloService",
		Version:     "1.0.0",
		Description: "Say hello",
		Endpoint: micro.Endpoint{
			Subject: "hello",
			Handler: helloHandler,
		},

		// DoneHandler can be set to customize behavior on stopping a service.
		DoneHandler: func(srv micro.Service) {
			info := srv.Info()
			fmt.Printf("Stopped service %q with ID %q\n", info.Name, info.ID)
		},

		// ErrorHandler can be used to customize behavior on service execution error.
		ErrorHandler: func(srv micro.Service, err *micro.NATSError) {
			info := srv.Info()
			fmt.Printf("Service %q returned an error on subject %q: %s", info.Name, err.Subject, err.Description)
		},
	}

	fmt.Printf("Starting service %q...\n", config.Name)

	srv, err := micro.AddService(nc, config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Started service %q with ID %q\n", srv.Info().Name, srv.Info().ID)

	return srv, nil
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Printf("Starting NATS microservice hosting infrastructure... (CTRL-C to halt)\n")

	nc, err := nats.Connect("localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	mySvc, err := AddHelloService(nc)
	if err != nil {
		fmt.Printf("Could not add service: %s", err.Error())
		return
	}
	defer func() {
		mySvc.Stop()
		time.Sleep(5 * time.Millisecond)
	}()

	select {
	case <-signalChan:
		signal.Stop(signalChan)
		fmt.Printf("\nHalting NATS microservice hosting infrastructure...\n")
		return
	}
}
