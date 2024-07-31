package main

import (
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/api"
)

func main() {

	//rabbitmq.RabbitMQServer()
	server := api.NewServer(".")
	server.Start(4000)
}
