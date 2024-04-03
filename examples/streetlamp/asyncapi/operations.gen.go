// DO NOT EDIT. This file is autogenerated by aapi-codegen
package asyncapi

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
)

var (
	ErrUnknownMessageId          = errors.New("unknown message id")
	ErrExpectedStruct            = errors.New("expected a struct")
	ErrChannelMismatchesAddress  = errors.New("channel mismatches address")
	ErrParameterCouldNotSetValue = errors.New("could not set value on parameter")
)

// Receive messages of channels by implementing this interface
type AsyncApiInterface interface {
	// Handles operation ReceiveStreetlights with message streetlight on the Streetlights channel
	ReceiveStreetlightsWithStreetlightMsg(msg StreetlightRecvMsg) error

	// Handles operation ReceiveLightMeasurement with message lightMeasured on the LightingMeasured channel
	ReceiveLightMeasurementWithLightMeasuredMsg(msg LightMeasuredRecvMsg, param LightingMeasuredParam) error

	// Handles operation ReceiveTurnOn with message turnOnOff on the LightTurnOn channel
	ReceiveTurnOnWithTurnOnOffMsg(msg TurnOnOffRecvMsg, param LightTurnOnParam) error

	// Handles operation ReceiveDimLight with message dimLight on the LightsDim channel
	ReceiveDimLightWithDimLightMsg(msg DimLightRecvMsg, param LightsDimParam) error
}

// Router plugin with all operations
type AsyncApi struct {
	handlers   AsyncApiInterface
	subscriber message.Subscriber

	publisher message.Publisher

	router *message.Router
}

func NewAsyncApi(publisher message.Publisher, subscriber message.Subscriber, handlers AsyncApiInterface) *AsyncApi {
	return &AsyncApi{
		publisher:  publisher,
		subscriber: subscriber,
		handlers:   handlers,

		router: nil,
	}
}

func (a *AsyncApi) Plugin(r *message.Router) error {
	a.router = r

	// Handler for receive operation ReceiveStreetlights
	r.AddNoPublisherHandler(
		"smartylighting.streetlights.1.0.streetlightsHandler",
		"smartylighting.streetlights.1.0.streetlights",
		a.subscriber,
		a.wrapperReceiveStreetlights,
	)

	return nil
}

// Subscribes to channel LightingMeasured
func (a AsyncApi) SubscribeToReceiveLightMeasurement(param LightingMeasuredParam) (*message.Handler, error) {
	if a.router == nil {
		panic("plugin uninitialised. call router.AddPlugin(asyncApi.Plugin) first")
	}

	address := marshalChannelAddress("smartylighting.streetlights.1.0.event.{streetlightId}.lighting.measured", param)

	handler := a.router.AddNoPublisherHandler(
		fmt.Sprintf("%sHandler", address),
		address,
		a.subscriber,
		a.wrapperReceiveLightMeasurement,
	)

	if err := a.router.RunHandlers(context.Background()); err != nil {
		return nil, err
	}

	return handler, nil
}

// Subscribes to channel LightTurnOn
func (a AsyncApi) SubscribeToReceiveTurnOn(param LightTurnOnParam) (*message.Handler, error) {
	if a.router == nil {
		panic("plugin uninitialised. call router.AddPlugin(asyncApi.Plugin) first")
	}

	address := marshalChannelAddress("smartylighting.streetlights.1.0.action.{streetlightId}.turn.on", param)

	handler := a.router.AddNoPublisherHandler(
		fmt.Sprintf("%sHandler", address),
		address,
		a.subscriber,
		a.wrapperReceiveTurnOn,
	)

	if err := a.router.RunHandlers(context.Background()); err != nil {
		return nil, err
	}

	return handler, nil
}

// Subscribes to channel LightsDim
func (a AsyncApi) SubscribeToReceiveDimLight(param LightsDimParam) (*message.Handler, error) {
	if a.router == nil {
		panic("plugin uninitialised. call router.AddPlugin(asyncApi.Plugin) first")
	}

	address := marshalChannelAddress("smartylighting.streetlights.1.0.action.{streetlightId}.dim", param)

	handler := a.router.AddNoPublisherHandler(
		fmt.Sprintf("%sHandler", address),
		address,
		a.subscriber,
		a.wrapperReceiveDimLight,
	)

	if err := a.router.RunHandlers(context.Background()); err != nil {
		return nil, err
	}

	return handler, nil
}

func marshalChannelAddress(address string, param interface{}) string {
	v := reflect.ValueOf(param)
	for i := 0; i < v.NumField(); i++ {
		pName := v.Type().Field(i).Tag.Get("parameter")
		pValue := v.Field(i).String()

		address = strings.ReplaceAll(address, fmt.Sprintf("{%s}", pName), pValue)
	}

	return address
}

func unmarshalChannelAddress(channel string, address string, param interface{}) error {
	if t := reflect.TypeOf(param); t.Kind() != reflect.Pointer {
		return ErrExpectedPointerType
	}

	pStruct := reflect.ValueOf(param).Elem()
	addressRegex := fmt.Sprintf("^%s$", address)
	for i := 0; i < pStruct.NumField(); i++ {
		pName := pStruct.Type().Field(i).Tag.Get("parameter")
		pField := pStruct.Type().Field(i).Name

		addressRegex = strings.ReplaceAll(addressRegex, fmt.Sprintf("{%s}", pName), fmt.Sprintf("(?<%s>.*)", pField))
	}

	re, err := regexp.Compile(addressRegex)
	if err != nil {
		return err
	} else if !re.MatchString(channel) {
		return ErrChannelMismatchesAddress
	} else if pStruct.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}

	subMatch := re.FindStringSubmatch(channel)
	for i, pField := range re.SubexpNames() {
		if i == 0 {
			continue
		}

		f := pStruct.FieldByName(pField)
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
			f.SetString(subMatch[i])
		} else {
			return ErrParameterCouldNotSetValue
		}
	}

	return nil
}

// Parameter of channel LightingMeasured
type LightingMeasuredParam struct {
	StreetlightId string `parameter:"streetlightId"`
}

// Parameter of channel LightTurnOn
type LightTurnOnParam struct {
	StreetlightId string `parameter:"streetlightId"`
}

// Parameter of channel LightsDim
type LightsDimParam struct {
	StreetlightId string `parameter:"streetlightId"`
}

// Parameter of channel LightTurnOff
type LightTurnOffParam struct {
	StreetlightId string `parameter:"streetlightId"`
}

// Message interface for operation PlaceStreetlight
type PlaceStreetlightSendMsg interface {
	visitPlaceStreetlightSendMsg(uuid string) (*Message, error)
}

func (m StreetlightMsgPayload) visitPlaceStreetlightSendMsg(uuid string) (*Message, error) {
	msg, err := NewMessage(uuid, m)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Publish message of operation PlaceStreetlight into channel Streetlights
func (a AsyncApi) PublishPlaceStreetlight(uuid string, message PlaceStreetlightSendMsg) error {
	msg, err := message.visitPlaceStreetlightSendMsg(uuid)
	if err != nil {
		return err
	}
	return a.publisher.Publish("smartylighting.streetlights.1.0.streetlights", msg.Raw())
}

// Message interface for operation LightMeasurement
type LightMeasurementSendMsg interface {
	visitLightMeasurementSendMsg(uuid string) (*Message, error)
}

func (m LightMeasuredMsgPayload) visitLightMeasurementSendMsg(uuid string) (*Message, error) {
	msg, err := NewMessage(uuid, m)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Publish message of operation LightMeasurement into channel LightingMeasured
func (a AsyncApi) PublishLightMeasurement(uuid string, message LightMeasurementSendMsg, param LightingMeasuredParam) error {
	msg, err := message.visitLightMeasurementSendMsg(uuid)
	if err != nil {
		return err
	}
	address := marshalChannelAddress("smartylighting.streetlights.1.0.event.{streetlightId}.lighting.measured", param)
	return a.publisher.Publish(address, msg.Raw())
}

// Message interface for operation TurnOn
type TurnOnSendMsg interface {
	visitTurnOnSendMsg(uuid string) (*Message, error)
}

func (m TurnOnOffMsgPayload) visitTurnOnSendMsg(uuid string) (*Message, error) {
	msg, err := NewMessage(uuid, m)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Publish message of operation TurnOn into channel LightTurnOn
func (a AsyncApi) PublishTurnOn(uuid string, message TurnOnSendMsg, param LightTurnOnParam) error {
	msg, err := message.visitTurnOnSendMsg(uuid)
	if err != nil {
		return err
	}
	address := marshalChannelAddress("smartylighting.streetlights.1.0.action.{streetlightId}.turn.on", param)
	return a.publisher.Publish(address, msg.Raw())
}

// Message interface for operation TurnOff
type TurnOffSendMsg interface {
	visitTurnOffSendMsg(uuid string) (*Message, error)
}

func (m TurnOnOffMsgPayload) visitTurnOffSendMsg(uuid string) (*Message, error) {
	msg, err := NewMessage(uuid, m)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Publish message of operation TurnOff into channel LightTurnOff
func (a AsyncApi) PublishTurnOff(uuid string, message TurnOffSendMsg, param LightTurnOffParam) error {
	msg, err := message.visitTurnOffSendMsg(uuid)
	if err != nil {
		return err
	}
	address := marshalChannelAddress("smartylighting.streetlights.1.0.action.{streetlightId}.turn.off", param)
	return a.publisher.Publish(address, msg.Raw())
}

// Message interface for operation DimLight
type DimLightSendMsg interface {
	visitDimLightSendMsg(uuid string) (*Message, error)
}

func (m DimLightMsgPayload) visitDimLightSendMsg(uuid string) (*Message, error) {
	msg, err := NewMessage(uuid, m)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Publish message of operation DimLight into channel LightsDim
func (a AsyncApi) PublishDimLight(uuid string, message DimLightSendMsg, param LightsDimParam) error {
	msg, err := message.visitDimLightSendMsg(uuid)
	if err != nil {
		return err
	}
	address := marshalChannelAddress("smartylighting.streetlights.1.0.action.{streetlightId}.dim", param)
	return a.publisher.Publish(address, msg.Raw())
}

// Wraps the handler for the receive operation ReceiveStreetlights
func (a AsyncApi) wrapperReceiveStreetlights(m *message.Message) error {
	msg := Message{}
	if err := msg.Load(m); err != nil {
		return err
	}

	switch msg.MsgId() {
	case "Streetlight":
		var recvMsg = StreetlightRecvMsg{
			Message: msg,
		}

		return a.handlers.ReceiveStreetlightsWithStreetlightMsg(recvMsg)
	default:
		return ErrUnknownMessageId
	}
}

// Wraps the handler for the receive operation ReceiveLightMeasurement
func (a AsyncApi) wrapperReceiveLightMeasurement(m *message.Message) error {
	msg := Message{}
	if err := msg.Load(m); err != nil {
		return err
	}

	// extract parameters from channel
	channel := message.SubscribeTopicFromCtx(m.Context())
	var param LightingMeasuredParam
	if err := unmarshalChannelAddress(channel, "smartylighting.streetlights.1.0.event.{streetlightId}.lighting.measured", &param); err != nil {
		return err
	}

	switch msg.MsgId() {
	case "LightMeasured":
		var recvMsg = LightMeasuredRecvMsg{
			Message: msg,
		}

		return a.handlers.ReceiveLightMeasurementWithLightMeasuredMsg(recvMsg, param)
	default:
		return ErrUnknownMessageId
	}
}

// Wraps the handler for the receive operation ReceiveTurnOn
func (a AsyncApi) wrapperReceiveTurnOn(m *message.Message) error {
	msg := Message{}
	if err := msg.Load(m); err != nil {
		return err
	}

	// extract parameters from channel
	channel := message.SubscribeTopicFromCtx(m.Context())
	var param LightTurnOnParam
	if err := unmarshalChannelAddress(channel, "smartylighting.streetlights.1.0.action.{streetlightId}.turn.on", &param); err != nil {
		return err
	}

	switch msg.MsgId() {
	case "TurnOnOff":
		var recvMsg = TurnOnOffRecvMsg{
			Message: msg,
		}

		return a.handlers.ReceiveTurnOnWithTurnOnOffMsg(recvMsg, param)
	default:
		return ErrUnknownMessageId
	}
}

// Wraps the handler for the receive operation ReceiveDimLight
func (a AsyncApi) wrapperReceiveDimLight(m *message.Message) error {
	msg := Message{}
	if err := msg.Load(m); err != nil {
		return err
	}

	// extract parameters from channel
	channel := message.SubscribeTopicFromCtx(m.Context())
	var param LightsDimParam
	if err := unmarshalChannelAddress(channel, "smartylighting.streetlights.1.0.action.{streetlightId}.dim", &param); err != nil {
		return err
	}

	switch msg.MsgId() {
	case "DimLight":
		var recvMsg = DimLightRecvMsg{
			Message: msg,
		}

		return a.handlers.ReceiveDimLightWithDimLightMsg(recvMsg, param)
	default:
		return ErrUnknownMessageId
	}
}
