package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Fair-Bytes/aapi-codegen/examples/streetlamp/asyncapi"
)

var (
	defaultLogger = watermill.NewStdLogger(true, false)
)

func main() {
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, defaultLogger)

	conn, err := sql.Open("pgx", "postgresql://localhost:5432/aapicodegen_db?user=aapicodegen_user&password=aapicodegen_pass")
	if err != nil {
		panic(err)
	}

	h := &handlers{}
	asyncApi, err := asyncapi.NewAsyncApi(pubsub, pubsub, h, asyncapi.PostboxConfig{
		PostgresDB: conn,
	})
	if err != nil {
		panic(err)
	}
	h.api = asyncApi

	router, err := message.NewRouter(message.RouterConfig{}, defaultLogger)
	if err != nil {
		panic(err)
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddPlugin(asyncApi.Plugin)

	streetlightId := uuid.New()
	go placeStreetlightAfterDelay(asyncApi, streetlightId)
	go sendTestMessages(asyncApi, streetlightId, conn)

	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}

func sendTestMessages(a *asyncapi.AsyncApi, streetlightId uuid.UUID, db *sql.DB) {
	lumens := 0

	for {
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		pub, err := a.NewTx(tx)
		if err != nil {
			panic(err)
		}

		lumens += 10
		now := time.Now()
		err = pub.PublishLightMeasurement(uuid.NewString(), asyncapi.LightMeasuredMsgPayload{
			LightMeasuredPayload: asyncapi.LightMeasuredPayload{
				Lumens: &lumens,
				SentAt: &now,
			},
		}, asyncapi.LightingMeasuredParam{
			StreetlightId: streetlightId.String(),
		})
		if err != nil {
			panic(err)
		}

		tx.Commit()

		time.Sleep(5 * time.Second)
	}
}

func placeStreetlightAfterDelay(a *asyncapi.AsyncApi, streetlightId uuid.UUID) {
	time.Sleep(7 * time.Second)

	err := a.PublishPlaceStreetlight(uuid.NewString(), asyncapi.StreetlightMsgPayload{
		StreetlightPayload: asyncapi.StreetlightPayload{
			Id: streetlightId,
			Address: asyncapi.Address{
				Street: "Musterstraße",
				City:   "Musterhausen",
			},
		},
	})
	if err != nil {
		panic(err)
	}
}

type handlers struct {
	api *asyncapi.AsyncApi
}

// Handles operation ReceiveStreetlights with message streetlight on the Streetlights channel
func (h handlers) ReceiveStreetlightsWithStreetlightMsg(msg asyncapi.StreetlightRecvMsg) error {
	payload, err := msg.UnmarshalPayload()
	if err != nil {
		return err
	}
	log.Printf("[NEW] Streetlight %s was placed in %s, %s\n", payload.Id.String(), payload.Address.Street, payload.Address.City)

	_, err = h.api.SubscribeToReceiveLightMeasurement(asyncapi.LightingMeasuredParam{
		StreetlightId: payload.Id.String(),
	})
	if err != nil {
		log.Println(err)
		msg.Nack()
	} else {
		msg.Ack()
	}

	return nil
}

// Handles operation ReceiveLightMeasurement with message lightMeasured on the LightingMeasured channel
func (handlers) ReceiveLightMeasurementWithLightMeasuredMsg(msg asyncapi.LightMeasuredRecvMsg, param asyncapi.LightingMeasuredParam) error {
	payload, err := msg.UnmarshalPayload()
	if err != nil {
		return err
	}
	log.Printf("[Measurement] Streetlight %s: %d at %s\n", param.StreetlightId, *payload.Lumens, payload.SentAt.Format(time.RFC3339))
	msg.Ack()
	return nil
}

// Handles operation ReceiveTurnOn with message turnOnOff on the LightTurnOn channel
func (handlers) ReceiveTurnOnWithTurnOnOffMsg(msg asyncapi.TurnOnOffRecvMsg, param asyncapi.LightTurnOnParam) error {
	payload, err := msg.UnmarshalPayload()
	if err != nil {
		return err
	}
	log.Printf("[ON] Streetlight %s: %s\n", param.StreetlightId, *payload.Command)
	msg.Ack()
	return nil
}

// Handles operation ReceiveDimLight with message dimLight on the LightsDim channel
func (handlers) ReceiveDimLightWithDimLightMsg(msg asyncapi.DimLightRecvMsg, param asyncapi.LightsDimParam) error {
	payload, err := msg.UnmarshalPayload()
	if err != nil {
		return err
	}
	log.Printf("[DIM] Streetlight %s:\n", param.StreetlightId)
	for _, dimLightPoint := range payload.DimLightPayload {
		log.Panicf("\t%d at %s\n", dimLightPoint.Percentage, dimLightPoint.SentAt.Format(time.RFC3339))
	}
	msg.Ack()
	return nil
}
