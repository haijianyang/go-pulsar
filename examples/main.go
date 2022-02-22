package main

import (
	"fmt"

	gopulsar "github.com/haijianyang/go-pulsar"
)

func main() {
	consumer, err := gopulsar.NewPulsarConsumer(&gopulsar.PulsarConsumerOptions{
		ClientOptions: gopulsar.ClientOptions{
			URL: "pulsar://localhost:6650",
		},
		ConsumerOptions: gopulsar.ConsumerOptions{
			Topic:            "my-topic",
			SubscriptionName: "my-topic-consumer",
			Type:             gopulsar.Shared,
		},
		BufferSize:  600,
		Concurrency: 60,
	})
	if err != nil {
		panic(err)
	}

	consumer.Run(func(message gopulsar.Message) error {
		fmt.Printf("handle message %+v\n", message)

		return nil
	})

	go func() {
		for err := range consumer.Errors {
			fmt.Printf("process error: %+v\n", err)
		}
	}()
}
