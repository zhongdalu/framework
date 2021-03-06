package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/panshiqu/framework/define"
	"github.com/panshiqu/framework/network"
	"github.com/panshiqu/framework/proxy"
	"github.com/panshiqu/framework/utils"
)

func handleSignal(server *network.Server, client *network.Client) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	s := <-c
	log.Println("Got signal:", s)

	server.Stop()
	client.Stop()
}

func main() {
	config := &define.ConfigProxy{}
	if err := utils.ReadJSON("./config/proxy.json", config); err != nil {
		log.Println("ReadJSON ConfigProxy", err)
		return
	}

	server := network.NewServer(config.ListenIP)
	client := network.NewClient(config.DialIP)
	processor := proxy.NewProcessor(server, client, config)

	server.Register(processor)
	client.Register(processor)

	go handleSignal(server, client)
	go client.Start()

	if err := server.Start(); err != nil {
		log.Println("Start", err)
	}
}
