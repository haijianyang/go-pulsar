package gopulsar

import "github.com/apache/pulsar-client-go/pulsar"

type Handle func(pulsar.Message) error

/* Common */

type (
	Message               = pulsar.Message
	Schema                = pulsar.Schema
	MessageDecryptionInfo = pulsar.MessageDecryptionInfo
)

/* Client */

type (
	ClientOptions = pulsar.ClientOptions
	Client        = pulsar.Client
)

/* Consumer */

var (
	Exclusive = pulsar.Exclusive
	Shared    = pulsar.Shared
	Failover  = pulsar.Failover
	KeyShared = pulsar.KeyShared

	SubscriptionPositionLatest   = pulsar.SubscriptionPositionLatest
	SubscriptionPositionEarliest = pulsar.SubscriptionPositionEarliest
)

type (
	ConsumerMessage      = pulsar.ConsumerMessage
	DLQPolicy            = pulsar.DLQPolicy
	KeySharedPolicy      = pulsar.KeySharedPolicy
	ConsumerInterceptors = pulsar.ConsumerInterceptors
	ConsumerOptions      = pulsar.ConsumerOptions
	Consumer             = pulsar.Consumer
)
