# go-pulsar
Apache Pulsar Consumer &amp; Worker for Go.

## Install

```console
go get github.com/haijianyang/go-pulsar
```

## Quick Start

```go
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
```

## Document

### Credentials

```go
// consumer will use the provided client.
client := pulsar.NewClient(pulsar.ClientOptions{URL: "pulsar://localhost:6650"})
consumer, err := gopulsar.NewPulsarConsumer(&gopulsar.PulsarConsumerOptions{
  Client: client,
})
```

```go
// consumer will create Pulsar client.
consumer, err := gopulsar.NewPulsarConsumer(&gopulsar.PulsarConsumerOptions{
	ClientOptions: gopulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	},
})
```

### Concurrent

```go
consumer, err := gopulsar.NewPulsarConsumer(&gopulsar.PulsarConsumerOptions{
	ConsumerOptions: gopulsar.ConsumerOptions{},
	BufferSize:  600, // message queue buffer size is 600.
	Concurrency: 60, // consumer will create 60 goroutines to handle the messages.
})
```