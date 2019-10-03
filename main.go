package main

import (
	"os/signal"

	"log"
	"os"
	"telegrambot/chat"
	"telegrambot/config"
	"telegrambot/order"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v\n", err)
	}

	orderService := order.NewInMemoryService(conf)

	chat, err := chat.NewService(
		conf,
		orderService,
	)
	if err != nil {
		log.Fatalf("creating chat service: %s", err)
	}

	err = chat.Run()
	if err != nil {
		log.Fatalf("running chat service: %s", err)
	}

	// Shutdown on SIGINT (CTRL-C).
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	chat.Stop()

	log.Fatalf("shutting down...")
}
