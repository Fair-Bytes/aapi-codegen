package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/Fair-Bytes/aapi-codegen/examples/streetlamp/asyncapi"
)

var (
	defaultLogger = watermill.NewStdLogger(true, false)
)

func main() {
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, defaultLogger)
	asyncApi := asyncapi.NewAsyncApi(pubsub, pubsub, handlers{})

	router, err := message.NewRouter(message.RouterConfig{}, defaultLogger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddPlugin(asyncApi.Plugin)

	router.Run(context.Background())
}

type handlers struct{}

// Handles operation ReceiveLightMeasurement with message LightMeasured on a channel with parameters
func (handlers) ReceiveLightMeasurement_LightMeasured(msg asyncapi.Message, payload asyncapi.LightingMeasured_LightMeasured, parameters asyncapi.LightingMeasuredParameters) error {
	msg.Ack()
	return nil
}

// Handles operation ReceiveTurnOn with message TurnOnOff on a channel with parameters
func (handlers) ReceiveTurnOn_TurnOnOff(msg asyncapi.Message, payload asyncapi.LightTurnOn_TurnOnOff, parameters asyncapi.LightTurnOnParameters) error {
	msg.Ack()
	return nil
}

// Handles operation ReceiveDimLight with message DimLight on a channel with parameters
func (handlers) ReceiveDimLight_DimLight(msg asyncapi.Message, payload asyncapi.LightsDim_DimLight, parameters asyncapi.LightsDimParameters) error {
	msg.Ack()
	return nil
}
