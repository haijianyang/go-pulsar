package gopulsar

import (
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/pkg/errors"
)

// PulsarConsumerOptions provide options to configure the PulsarConsumer.
type PulsarConsumerOptions struct {
	// ClientOptions is used to construct a Pulsar Client instance.
	ClientOptions
	// ConsumerOptions is used to configure and create instances of Consumer
	ConsumerOptions
	// Client supersedes the ClientOptions field, and allows you to inject
	// an already-made Pulsar Client for use in the consumer.
	Client Client
	// BufferSize determines the size of the channel uses to coordinate the
	// processing of the messages. This determines the maximum number of
	// in-flight messages.
	BufferSize int
	// Concurrency dictates how many goroutines to spawn to handle the messages.
	Concurrency int
	// MaxDeliveryCount determines the maximum delivery number of the message.
	MaxDeliveryCount int
}

// PulsarConsumer adds a convenient wrapper around dequeuing and managing concurrency.
type PulsarConsumer struct {
	// Errors is a channel that you can receive from to centrally handle any
	// errors that may occur either by your Handle or by internal
	// processing functions. Because this is an unbuffered channel, you must
	// have a listener on it. If you don't parts of the consumer could stop
	// functioning when errors occur due to the blocking nature of unbuffered
	// channels.
	Errors   chan error
	options  *PulsarConsumerOptions
	client   Client
	consumer Consumer
	handle   Handle
	wg       *sync.WaitGroup
	stopCh   chan struct{}
}

// NewPulsarConsumer creates a PulsarConsumer with custom PulsarConsumerOptions.
func NewPulsarConsumer(options *PulsarConsumerOptions) (*PulsarConsumer, error) {
	client := options.Client
	if client == nil {
		var err error
		client, err = pulsar.NewClient(options.ClientOptions)
		if err != nil {
			return nil, err
		}
	}

	options.MessageChannel = make(chan ConsumerMessage, options.BufferSize)
	consumer, err := client.Subscribe(options.ConsumerOptions)
	if err != nil {
		return nil, err
	}

	return &PulsarConsumer{
		options:  options,
		client:   client,
		consumer: consumer,
		stopCh:   make(chan struct{}, 1),
		wg:       &sync.WaitGroup{},
		Errors:   make(chan error),
	}, nil
}

// Run starts all of the worker goroutines and starts processing.
// All errors will be sent to the Errors channel.
// Run will block until Shutdown is called
// and all of the in-flight messages have been processed.
func (c *PulsarConsumer) Run(handle Handle) {
	c.handle = handle

	stop := newSignalHandler()
	go func() {
		<-stop
		c.Shutdown()
	}()

	c.wg.Add(c.options.Concurrency)

	for i := 0; i < c.options.Concurrency; i++ {
		go c.work()
	}

	c.wg.Wait()
}

// work is called in a separate goroutine. The number of work goroutines is
// determined by Concurreny. Once it gets a message from the message channel,
// it calls the corrensponding Handle depending on the subscription it
// came from. If no error is returned from the Handle, the message is
// acknowledged in Pulsar.
func (c *PulsarConsumer) work() {
	defer c.wg.Done()

	for {
		select {
		case message := <-c.options.MessageChannel:
			err := c.process(message)
			if err != nil {
				c.Errors <- errors.Wrapf(err, "error calling handle for %q topic and %q message", message.Topic(), message.ID())

				if c.options.MaxDeliveryCount <= 0 || c.options.MaxDeliveryCount > int(message.RedeliveryCount()+1) {
					continue
				}
			}

			c.consumer.Ack(message)
		case <-c.stopCh:
			return
		}
	}
}

func (c *PulsarConsumer) process(message Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = errors.Wrap(e, "handle panic")
				return
			}
			err = errors.Errorf("handle panic: %v", r)
		}
	}()

	err = c.handle(message)

	return
}

func (c *PulsarConsumer) ReconsumeLater(msg Message, delay time.Duration) {
	c.consumer.ReconsumeLater(msg, delay)
}

func (c *PulsarConsumer) Shutdown() {
	c.stopCh <- struct{}{}
}
