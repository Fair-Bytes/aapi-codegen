// Code generated by https://github.com/Fair-Bytes/aapi-codegen version 0.0.2 DO NOT EDIT.
package asyncapi

import (
	"errors"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
)

var (
	ErrConfigMissingSubscriberOrDB = errors.New("postbox config expected either subscriber or postgresDB")
)

var (
	defaultPostboxTopic             = "asyncapi_postbox"
	defaultConsumerGroup            = "asyncapi_group"
	defaultLogger                   = watermill.NewStdLogger(true, false)
	defaultPostgresSubscriberConfig = sql.SubscriberConfig{
		SchemaAdapter:    sql.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   sql.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
		ConsumerGroup:    defaultConsumerGroup,
	}
	defaultPostgresPublisherConfig = sql.PublisherConfig{
		SchemaAdapter: sql.DefaultPostgreSQLSchema{},
	}
)

type PostboxConfig struct {
	Logger          watermill.LoggerAdapter
	PostgresDB      sql.Beginner
	Subscriber      *sql.Subscriber
	PublisherConfig *sql.PublisherConfig
	PostboxTopic    string
}

func (c *PostboxConfig) setDefaultsAndValidate() error {
	if c.Logger == nil {
		c.Logger = defaultLogger
	}

	if c.Subscriber == nil && c.PostgresDB != nil {
		if defaultSubscriber, err := sql.NewSubscriber(
			c.PostgresDB,
			defaultPostgresSubscriberConfig,
			c.Logger,
		); err != nil {
			return err
		} else {
			c.Subscriber = defaultSubscriber
		}
	} else if c.Subscriber == nil && c.PostgresDB == nil {
		return ErrConfigMissingSubscriberOrDB
	}

	if c.PublisherConfig == nil {
		c.PublisherConfig = &defaultPostgresPublisherConfig
	}

	if c.PostboxTopic == "" {
		c.PostboxTopic = defaultPostboxTopic
	}

	return nil
}

type postboxPublisher struct {
	publisher     message.Publisher
	postboxConfig PostboxConfig
}

func (p postboxPublisher) NewTx(tx sql.ContextExecutor) (*postboxPublisher, error) {
	pub, err := sql.NewPublisher(tx, *p.postboxConfig.PublisherConfig, p.postboxConfig.Logger)
	if err != nil {
		return nil, err
	}

	publisher := forwarder.NewPublisher(pub, forwarder.PublisherConfig{
		ForwarderTopic: p.postboxConfig.PostboxTopic,
	})

	return &postboxPublisher{
		publisher: publisher,
	}, err
}
